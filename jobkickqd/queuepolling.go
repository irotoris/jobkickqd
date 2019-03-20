package jobkickqd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/sirupsen/logrus"
)

type JobQueueExecutor interface {
	Run()
}

type PubSubJobQueue struct {
	projectID    string
	pubsubClient pubsub.Client
	topic        pubsub.Topic
	subscription pubsub.Subscription
	daemonApp    string
}

func NewPubSubJobQueueExecutor(ctx context.Context, projectID, topicName, subscriptionName, daemonApp string) (*PubSubJobQueue, error) {
	qd := new(PubSubJobQueue)
	qd.projectID = projectID
	pubsubClient, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	qd.pubsubClient = *pubsubClient
	topic := pubsubClient.Topic(topicName)
	sub, err := pubsubClient.CreateSubscription(ctx, subscriptionName, pubsub.SubscriptionConfig{
		Topic:       topic,
		AckDeadline: 10 * time.Second,
	})
	if err != nil {
		logrus.Warnf("%s", err)
		sub = pubsubClient.Subscription(subscriptionName)
	}
	qd.topic = *topic
	qd.subscription = *sub
	qd.daemonApp = daemonApp

	return qd, nil
}

func (jq *PubSubJobQueue) Run(ctx, cctx context.Context, ld PubSubMessageDriver) error {
	logrus.Infof("Start job queue polling and command executor. project:%s, job queue topic:%s, subscription:%s, command log topic:%s", jq.projectID, jq.topic.ID(), jq.subscription.ID(), ld.topicName)
	err := jq.subscription.Receive(cctx, func(ctx context.Context, m *pubsub.Message) {
		logrus.Infof("Received a job message:%s :%s :%s :%s", m.PublishTime, m.ID, string(m.Data), m.Attributes)

		// TODO: check and stop duplicate execution
		m.Ack()

		if m.Attributes["app"] != jq.daemonApp {
			logrus.Debugf("This message app is not match. %s != %s", jq.daemonApp, m.Attributes["app"], )
			return
		}
		now := time.Now()
		if now.Sub(m.PublishTime) > 5 * time.Minute {
			logrus.Warnf("This message was published more than 5 minutes ago. Skip job...")
			return
		}

		var jm DefaultJobMessage
		if err := json.Unmarshal(m.Data, &jm); err != nil {
			logrus.Errorf("json.Unmarshal() failed.: %s", err)
			return
		}
		timeoutDuration := time.Duration(jm.Timeout*1) * time.Second

		// TODO: to async
		jobExecutionID := jm.JobID + m.ID

		j := NewJob(jm.JobID, jobExecutionID, jm.Command, jm.Environment, timeoutDuration)
		if err := j.Execute(ctx); err != nil {
			logrus.Errorf("Failed to create new job object.: %s", err)
			return
		}
		attributes := map[string]string{"app": jq.daemonApp, "job_execution_id": jobExecutionID, "job_exit_code": fmt.Sprintf("%d", j.JobExitCode)}

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
