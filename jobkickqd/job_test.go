package jobkickqd

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestExecute(t *testing.T) {
	ctx := context.Background()
	jobID := "testJob1"
	cmdString := "echo \"command success test, env is ${ENV}.\""
	envs := []string{"ENV=test"}
	timeout := 60 * time.Second
	j := NewJob(jobID, jobID, cmdString, envs, timeout)

	err := j.Execute(ctx)
	if err != nil {
		t.Errorf("Execute command is failed.:%v", err)
	}
	if err := os.RemoveAll(jobID); err != nil {
		t.Errorf("post script is failed in job_test.go. %v", err)
	}
}

func TestKill(t *testing.T) {
	ctx := context.Background()

	jobID := "testJob2"
	cmdString := "sleep 1"
	var envs []string
	timeout := 60 * time.Second
	j := NewJob(jobID, jobID, cmdString, envs, timeout)

	go j.Execute(ctx)

	time.Sleep(500 * time.Millisecond)
	err := j.Kill(ctx)
	if err != nil {
		t.Errorf("Kill Process is failed.:%v", err)
	}
	if err := os.RemoveAll(jobID); err != nil {
		t.Errorf("post script is failed in job_test.go. %v", err)
	}
}
