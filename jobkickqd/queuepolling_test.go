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

	topicNameForLog, ok := os.LookupEnv("topicNameForLog")
	if !ok {
		t.Error("topicName is required.")
	}

	topicNameForJobQueue, ok := os.LookupEnv("topicNameForJobQueue")
	if !ok {
		t.Error("topicName is required.")
	}

	subscriptionName, ok := os.LookupEnv("subscriptionName")
	if !ok {
		t.Error("subscriptionName is required.")
	}

	ld, err := NewPubSubLogDriver(ctx, projectID, topicNameForLog)
	if err != nil {
		fmt.Println("err", err)
		t.Error("NewPubSubLogDriver is failed.")
	}

	qd, err := NewPubSubJobQueueExecutor(ctx, projectID, topicNameForJobQueue, subscriptionName)
	if err != nil {
		fmt.Println("err", err)
		t.Error("NewPubSubJobQueue is failed.")
	}

	cctx, cancel := context.WithCancel(ctx)

	qd.Run(ctx, cctx, *ld)


	kickq, err := NewPubSubLogDriver(ctx, projectID, topicNameForJobQueue)
	if err != nil {
		fmt.Println("err", err)
		t.Error("NewPubSubLogDriver is failed.")
	}
	time.Sleep(5 * time.Second)
	msg := "{\"job_id\":\"test-from-queue\",\"job_execution_id\":\"test-00001\",\"command\":\"echo ${ENV}\",\"Environment\":[\"ENV=dev\",\"ROLE=test\"],\"timeout\":60,\"retry\":0}"
	attribute := map[string]string{
		"app": "jobkickqd",
	}

	err = kickq.Write(ctx, msg, attribute)
	if err != nil {
		fmt.Println("err", err)
		t.Error("Write() is failed.")
	}

	time.Sleep(30 * time.Second)
	cancel()

	if err != nil {
		fmt.Println("err", err)
		t.Error("PubSubJobQueueSubscribe() is failed.")
	}

}
