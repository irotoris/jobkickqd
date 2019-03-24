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

	logTopic, ok := os.LookupEnv("logTopic")
	if !ok {
		t.Error("logTopic is required.")
	}

	jobQueueTopic, ok := os.LookupEnv("jobQueueTopic")
	if !ok {
		t.Error("jobQueueTopic is required.")
	}

	logDriver, err := NewPubSubMessageDriver(ctx, projectID, logTopic)
	if err != nil {
		fmt.Println("err", err)
		t.Error("NewPubSubMessageDriver is failed.")
	}

	queueDriver, err := NewPubSubJobQueueExecutor(ctx, projectID, jobQueueTopic, "test-app", "test-app")
	if err != nil {
		fmt.Println("err", err)
		t.Error("NewPubSubJobQueueExecutor is failed.")
	}

	cctx, cancel := context.WithCancel(ctx)

	// start job queue polling
	go queueDriver.Run(ctx, cctx, *logDriver)

	// publish test job
	kickq, err := NewPubSubMessageDriver(ctx, projectID, jobQueueTopic)
	if err != nil {
		fmt.Println("err", err)
		t.Error("NewPubSubMessageDriver is failed.")
	}

	time.Sleep(1 * time.Second)

	msgs := []string{
		"{\"jobID\":\"test-from-queue\",\"command\":\"sleep 1;echo \\\"env is ${ENV}\\\"\",\"Environment\":[\"ENV=dev1\",\"ROLE=test1\"],\"timeout\":60}",
		"{\"jobID\":\"test-from-queue\",\"command\":\"echo \\\"env is ${ENV}\\\"\",\"Environment\":[\"ENV=dev2\",\"ROLE=test2\"],\"timeout\":60}",
		"{\"jobID\":\"test-from-queue\",\"command\":\"echo \\\"env is ${ENV}\\\"\",\"Environment\":[\"ENV=dev3\",\"ROLE=test3\"],\"timeout\":60}",
		"{\"jobID\":\"test-from-queue\",\"command\":\"sleep 3;echo \\\"env is ${ENV}\\\"\",\"Environment\":[\"ENV=dev4\",\"ROLE=test4\"],\"timeout\":60}",
		"{\"jobID\":\"test-from-queue\",\"command\":\"sleep 300\",\"Environment\":[\"ENV=dev\",\"ROLE=test\"],\"timeout\":3}",
	}
	attribute := map[string]string{
		"app": "test-app",
	}
	msgIDs := []string{}
	for _, msg := range msgs {
		id, err := kickq.Write(ctx, msg, attribute)
		if err != nil {
			fmt.Println("err", err)
			t.Error("Write() is failed.")
		}
		msgIDs = append(msgIDs, id)
	}

	time.Sleep(10 * time.Second)
	cancel()

	for _, msgID := range msgIDs {
		dir := "test-from-queue" + msgID
		if err := os.RemoveAll(dir); err != nil {
			t.Errorf("post script is failed in job_test.go. %v", err)
		}
	}

}
