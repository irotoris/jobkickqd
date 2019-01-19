package jobkickqd

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestPubSubJobQueueSubscribe(t *testing.T) {

	ctx := context.Background()
	projectID, ok := os.LookupEnv("projectID")
	if !ok {
		t.Error("projectID is required.")
	}

	topicName, ok := os.LookupEnv("topicName")
	if !ok {
		t.Error("topicName is required.")
	}

	subscriptionName, ok := os.LookupEnv("subscriptionName")
	if !ok {
		t.Error("subscriptionName is required.")
	}

	qd, err := NewPubSubJobQueueExecutor(ctx, projectID, topicName, subscriptionName)
	if err != nil {
		fmt.Println("err", err)
		t.Error("NewPubSubJobQueue is failed.")
	}

	cctx, cancel := context.WithCancel(ctx)

	go func(queue PubSubJobQueue){
		qd.Run(ctx, cctx)
	}(*qd)

	time.Sleep(5 * time.Second)
	cancel()

	if err != nil {
		fmt.Println("err", err)
		t.Error("PubSubJobQueueSubscribe() is failed.")
	}

}
