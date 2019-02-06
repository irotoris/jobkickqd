package jobkickqd

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/sirupsen/logrus"
)

type JobQueueExecutor interface {
	Run()
}

type PubSubJobQueue struct {
	projectID        string
	subscriptionName string
	pubsubClient     pubsub.Client
	subscription     pubsub.Subscription
}

func NewPubSubJobQueueExecutor(ctx context.Context, projectID, subscriptionName string) (*PubSubJobQueue, error) {
	qd := new(PubSubJobQueue)
	qd.projectID = projectID
	qd.subscriptionName = subscriptionName
	pubsubClient, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	qd.pubsubClient = *pubsubClient
	sub := pubsubClient.Subscription(subscriptionName)
	qd.subscription = *sub

	return qd, nil
}

func (jq *PubSubJobQueue) Run(ctx, cctx context.Context, ld PubSubMessageDriver) error {
	logrus.Infof("Start job queue polling and command executor. project:%s, job queue subscription:%s, command log topic:%s", jq.projectID, jq.subscriptionName, ld.topicName)
	var mu sync.Mutex
	err := jq.subscription.Receive(cctx, func(ctx context.Context, m *pubsub.Message) {
		logrus.Infof("Received a job message:%s :%s :%s", m.ID, string(m.Data), m.Attributes)
		if m.Attributes["app"] != "jobkickqd" {
			return
		}
		// TODO: check and stop duplicate execution
		m.Ack()
		mu.Lock()
		defer mu.Unlock()
		var jm DefaultJobMessage
		if err := json.Unmarshal(m.Data, &jm); err != nil {
			logrus.Errorf("json.Unmarshal() failed.: %s", err)
			return
		}
		timeoutDuration := time.Duration(jm.Timeout*100) * time.Second

		// TODO: to async
		jobExecutionID := jm.JobID + m.ID
		attributes := map[string]string{"app": "jobkickqd", "job_execution_id": jobExecutionID}
		j := NewJob(jm.JobID, jobExecutionID, jm.Command, jm.Environment, timeoutDuration)
		if err := j.Execute(ctx); err != nil {
			logrus.Errorf("Failed to create new job object.: %s", err)
			return
		}
		if _, err := ld.Write(ctx, j.ExecutionLog, attributes); err != nil {
			logrus.Errorf("Failed to write log to a log driver.: %s", err)
			return
		}

	})
	if err != context.Canceled {
		return err
	}
	return nil
}
