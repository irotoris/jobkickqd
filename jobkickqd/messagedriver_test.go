package jobkickqd

import (
	"context"
	"fmt"
	"os"
	"testing"
)

func TestPubSubMessageDriver_Write(t *testing.T) {
	ctx := context.Background()
	projectID, ok := os.LookupEnv("projectID")
	if !ok {
		t.Error("projectID is required.")
	}

	topicNameForLog, ok := os.LookupEnv("topicNameForLog")
	if !ok {
		t.Error("topicName is required.")
	}

	ld, err := NewPubSubMessageDriver(ctx, projectID, topicNameForLog)
	if err != nil {
		fmt.Println("err", err)
		t.Error("NewPubSubLogDriver is failed.")
	}

	logMessages := [5]string{"test logs 1.", "test logs 2.", "test logs 3.", "test logs 4.", "test logs 5."}
	attribute := map[string]string{
		"jobID":          "test-job",
		"jobExecutionID": "test-job-execution-1",
	}

	for _, msg := range logMessages {
		_, err = ld.Write(ctx, msg, attribute)
		if err != nil {
			fmt.Println("err", err)
			t.Error("NewPubSubLogDriver is failed.")
		}
	}
}
