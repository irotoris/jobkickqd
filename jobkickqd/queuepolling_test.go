package jobkickqd

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestPubSubJobQueue_Run(t *testing.T) {

	ctx := context.Background()
	projectID, ok := os.LookupEnv("projectID")
	if !ok {
		t.Error("projectID is required.")
	}

	topicNameForLog, ok := os.LookupEnv("topicNameForLog")
	if !ok {
		t.Error("topicNameForLog is required.")
	}

	topicNameForJobQueue, ok := os.LookupEnv("topicNameForJobQueue")
	if !ok {
		t.Error("topicNameForJobQueue is required.")
	}

	subscriptionName, ok := os.LookupEnv("subscriptionName")
	if !ok {
		t.Error("subscriptionName is required.")
	}

	logDriver, err := NewPubSubMessageDriver(ctx, projectID, topicNameForLog)
	if err != nil {
		fmt.Println("err", err)
		t.Error("NewPubSubMessageDriver is failed.")
	}

	queueDriver, err := NewPubSubJobQueueExecutor(ctx, projectID, subscriptionName)
	if err != nil {
		fmt.Println("err", err)
		t.Error("NewPubSubJobQueueExecutor is failed.")
	}

	cctx, cancel := context.WithCancel(ctx)

	// start job queue polling
	go queueDriver.Run(ctx, cctx, *logDriver)

	// publish test job
	kickq, err := NewPubSubMessageDriver(ctx, projectID, topicNameForJobQueue)
	if err != nil {
		fmt.Println("err", err)
		t.Error("NewPubSubMessageDriver is failed.")
	}

	time.Sleep(1 * time.Second)

	msgs := []string{
		"{\"job_id\":\"test-from-queue\",\"job_execution_id\":\"test-00001\",\"command\":\"sleep 1;echo \\\"env is ${ENV}\\\"\",\"Environment\":[\"ENV=dev1\",\"ROLE=test1\"],\"timeout\":60,\"retry\":0}",
		"{\"job_id\":\"test-from-queue\",\"job_execution_id\":\"test-00002\",\"command\":\"echo \\\"env is ${ENV}\\\"\",\"Environment\":[\"ENV=dev2\",\"ROLE=test2\"],\"timeout\":60,\"retry\":0}",
		"{\"job_id\":\"test-from-queue\",\"job_execution_id\":\"test-00003\",\"command\":\"echo \\\"env is ${ENV}\\\"\",\"Environment\":[\"ENV=dev3\",\"ROLE=test3\"],\"timeout\":60,\"retry\":0}",
		"{\"job_id\":\"test-from-queue\",\"job_execution_id\":\"test-00004\",\"command\":\"sleep 3;echo \\\"env is ${ENV}\\\"\",\"Environment\":[\"ENV=dev4\",\"ROLE=test4\"],\"timeout\":60,\"retry\":0}",
		"{\"job_id\":\"test-from-queue\",\"job_execution_id\":\"test-00005\",\"command\":\"sleep 300\",\"Environment\":[\"ENV=dev\",\"ROLE=test\"],\"timeout\":3,\"retry\":0}",
	}
	attribute := map[string]string{
		"app": "jobkickqd",
	}
	for _, msg := range msgs {
		_, err = kickq.Write(ctx, msg, attribute)
		if err != nil {
			fmt.Println("err", err)
			t.Error("Write() is failed.")
		}
	}

	time.Sleep(10 * time.Second)
	cancel()

}
