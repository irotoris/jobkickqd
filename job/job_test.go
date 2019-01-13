package job

import (
	"context"
	"testing"
	"time"
)

func TestExecute(t *testing.T) {
	ctx := context.Background()
	jobID := "testJob"
	cmdString := "echo \"command success test, env is ${ENV}.\""
	envs := []string{"ENV=test"}
	attr := make(map[string]string)
	j := NewJob(jobID, cmdString, envs, attr)

	err := j.Execute(ctx)
	if err != nil {
		t.Error("Execute command is failed.")
	}
}

func TestKill(t *testing.T) {
	ctx := context.Background()

	jobID := "testJob"
	cmdString := "sleep 1"
	var envs []string
	attr := make(map[string]string)
	j := NewJob(jobID, cmdString, envs, attr)

	go j.Execute(ctx)

	time.Sleep(500 * time.Millisecond)
	err := j.Kill(ctx)
	if err != nil {
		t.Error("Kill Process is failed.")
	}
}
