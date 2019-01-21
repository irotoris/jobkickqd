package jobkickqd

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type JobQueueExecutor interface {
	Run()
}

type PubSubJobQueue struct {
	projectID string
	subscriptionName string
	pubsubClient pubsub.Client
	subscription pubsub.Subscription
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

	return qd ,nil
}

func (jq *PubSubJobQueue)Run(ctx, cctx context.Context, ld PubSubMessageDriver) error {
	logrus.Infof("Start job queue polling and command executor. project:%s, job queue subscription:%s, command log topic:%s", jq.projectID, jq.subscriptionName, ld.topicName)
	var mu sync.Mutex
	err := jq.subscription.Receive(cctx, func(ctx context.Context, m *pubsub.Message) {
		logrus.Infof("Received a job message: %s: %s", string(m.Data), m.Attributes)
		if m.Attributes["app"] != "jobkickqd" {
			return
		}
		// TODO: check and stop duplicate execution
		m.Ack()
		mu.Lock()
		defer mu.Unlock()
		var jm DefaultJobMessage
		attributes := make(map[string]string)
		if err := json.Unmarshal(m.Data, &jm); err != nil {
			logrus.Errorf("json.Unmarshal() failed.: %s", err)
			return
		}
		timeoutDuration := time.Duration(jm.Timeout*100) * time.Second

		// TODO: to async
		j := NewJob(jm.JobID, jm.Command, jm.Environment, timeoutDuration)
		if err := j.Execute(ctx); err != nil {
			logrus.Errorf("Failed to create new job object.: %s", err)
			return
		}
		if err := ld.Write(ctx, j.ExecutionLog, attributes); err != nil {
			logrus.Errorf("Failed to write log to a log driver.: %s", err)
			return
		}

	})
	if err != context.Canceled {
		return err
	}
	return nil
}
