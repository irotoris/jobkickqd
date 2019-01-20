package jobkickqd

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type JobQueueExecutor interface {
	Run()
}

type PubSubJobQueue struct {
	projectID string
	topicName string
	subscriptionName string
	pubsubClient pubsub.Client
	subscription pubsub.Subscription
}

func NewPubSubJobQueueExecutor(ctx context.Context, projectID, topicName, subscriptionName string) (*PubSubJobQueue, error) {
	qd := new(PubSubJobQueue)
	qd.projectID = projectID
	qd.topicName = topicName
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
	var mu sync.Mutex
	err := jq.subscription.Receive(cctx, func(ctx context.Context, m *pubsub.Message) {
		fmt.Println(string(m.Data), m.Attributes)
		if m.Attributes["app"] != "jobkickqd" {
			return
		}
		m.Ack()
		mu.Lock()
		defer mu.Unlock()
		// TODO: implement execute job
		var jm DefaultJobMessage
		attributes := make(map[string]string)
		if err := json.Unmarshal(m.Data, &jm); err != nil {
			fmt.Println("Err:json.Unmarshal is failed.")
			return
		}
		timeoutDuration := time.Duration(jm.Timeout*100) * time.Second
		j := NewJob(jm.JobID, jm.Command, jm.Environment, timeoutDuration)
		if err := j.Execute(ctx); err != nil {
			fmt.Println("Err:json.Unmarshal() is failed.")
			return
		}
		if err := ld.Write(ctx, j.ExecutionLog, attributes); err != nil {
			fmt.Println("Err:ld.Write() is failed.")
			return
		}

	})
	if err != context.Canceled {
		return err
	}
	return nil
}
