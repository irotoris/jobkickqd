package jobkickqd

import (
	"context"
	"fmt"
	"os"
	"testing"
)

func TestPubSubLogWrite(t *testing.T) {
	ctx := context.Background()
	projectID, ok := os.LookupEnv("projectID")
	if !ok {
		t.Error("projectID is required.")
	}

	topicName, ok := os.LookupEnv("topicName")
	if !ok {
		t.Error("topicName is required.")
	}

	ld, err := NewPubSubLogDriver(ctx, projectID, topicName)
	if err != nil {
		fmt.Println("err", err)
		t.Error("NewPubSubLogDriver is failed.")
	}

	logMessages := [5]string{"test logs 1.", "test logs 2.", "test logs 3.", "test logs 4.", "test logs 5."}
	attribute := map[string]string{
		"jobID":          "test-job",
		"jobEcexutionID": "test-job-execution-1",
	}

	for _, msg := range logMessages {
		err = ld.Write(ctx, msg, attribute)
		if err != nil {
			fmt.Println("err", err)
			t.Error("NewPubSubLogDriver is failed.")
		}
	}
}

func TestFileLogWrite(t *testing.T) {
	ctx := context.Background()
	logFilePath := "logs/test-job.log"
	ld, err := NewFileLogDriver(logFilePath)
	if err != nil {
		fmt.Println("err", err)
		t.Error("NewFileLogDriver is failed.")
	}
	defer ld.Close(ctx)

	logMessages := [5]string{"test logs 1.", "test logs 2.", "test logs 3.", "test logs 4.", "test logs 5."}
	for _, msg := range logMessages {
		err = ld.Write(ctx, msg)
		if err != nil {
			fmt.Println("err", err)
			t.Error("NewFileLogWrite() is failed.")
		}
	}

	_, err = os.Stat(logFilePath)
	if err != nil {
		fmt.Println("err", err)
		t.Error("NewFileLogWrite() is failed. A log file not found.")
	}
}
