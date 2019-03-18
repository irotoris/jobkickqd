package jobkickqd

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

// Job is...
type Job struct {
	JobID          string
	JobExecutionID string
	CommandString  string
	Environment    []string
	JobExitCode    int
	ExecutionLog   string
	SentAt         time.Time
	SubmittedAt    time.Time
	StartedAt      time.Time
	FinishedAt     time.Time
	Timeout        time.Duration
	Cmd            exec.Cmd
}

// NewJob is...
func NewJob(jobID, jobExecutionID, CommandString string, environment []string, timeout time.Duration) *Job {
	j := new(Job)
	j.JobID = jobID
	j.JobExecutionID = jobExecutionID
	j.CommandString = CommandString
	j.Environment = environment
	j.SubmittedAt = time.Now()
	j.Timeout = timeout
	return j
}

// Execute is...
func (j *Job) Execute(ctx context.Context) error {
	j.Cmd = *exec.Command("sh", "-c", j.CommandString)
	j.Cmd.Env = append(os.Environ())
	j.Cmd.Env = append(j.Environment)

	var stdout bytes.Buffer
	j.Cmd.Stderr = &stdout
	j.Cmd.Stdout = &stdout
	j.StartedAt = time.Now()


	j.Cmd.Start()

	// kill command when timeout
	var timer *time.Timer
	timer = time.AfterFunc(j.Timeout, func() {
		timer.Stop()
		j.Kill(ctx)
	})
	if err := j.Cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				logrus.Errorf("Exit Status: %d", status.ExitStatus())
				j.JobExitCode = status.ExitStatus()
			}
		} else {
			logrus.Errorf("cmd.Wait: %v", err)
			j.JobExitCode = 127
		}
	} else {
		j.JobExitCode = 0
	}

	j.FinishedAt = time.Now()
	j.ExecutionLog = stdout.String()

	return nil
}

// Kill is...
func (j *Job) Kill(ctx context.Context) error {
	err := j.Cmd.Process.Kill()
	if err != nil {
		return err
	}
	return nil
}
