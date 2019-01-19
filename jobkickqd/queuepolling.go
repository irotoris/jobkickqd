package jobkickqd

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"sync"
)

type JobQueueExecutor interface {
	Run()
	Stop()
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

func (qd *PubSubJobQueue)Run(ctx, cctx context.Context ) error {
	var mu sync.Mutex
	err := qd.subscription.Receive(cctx, func(ctx context.Context, m *pubsub.Message) {
		fmt.Println(string(m.Data), m.Attributes)
		m.Ack()
		mu.Lock()
		defer mu.Unlock()
	})
	if err != context.Canceled {
		return err
	}
	return nil
}
