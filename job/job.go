package job

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/dchest/uniuri"
)

// Job is...
type Job struct {
	JobID          string
	JobEcexutionID string
	JobAttributes  map[string]string
	CommandSring   string
	Environment    []string
	JobStatus      string
	ExecutionLog   string
	Cmd            exec.Cmd
}

// NewJob is...
func NewJob(jobID, CommandSring string, environment []string, jobAttributes map[string]string) *Job {
	j := new(Job)
	j.JobID = jobID
	j.JobEcexutionID = jobID + uniuri.New()
	j.CommandSring = CommandSring
	j.Environment = environment
	j.JobAttributes = jobAttributes
	return j
}

// Execute is...
func (job *Job) Execute(ctx context.Context) (bool, error) {
	cmd := exec.Command("sh", "-c", job.CommandSring)
	cmd.Env = append(os.Environ())
	cmd.Env = append(job.Environment)

	job.Cmd = *cmd

	// TODO:implement asyn
	out, err := cmd.Output()
	if err != nil {
		return false, err
	}

	// TODO:implement stdout/stderr/logdriver
	fmt.Print(string(out))

	// TODO:implement return PID or something infomation
	return true, nil
}

// Kill is...
func (job *Job) Kill(ctx context.Context) (bool, error) {
	err := job.Cmd.Process.Kill()
	if err != nil {
		return false, err
	}
	return true, nil
}
