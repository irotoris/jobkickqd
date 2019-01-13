package job

import (
	"context"
	"os"
	"os/exec"

	"github.com/dchest/uniuri"
)

// Job is...
type Job struct {
	JobID          string
	JobEcexutionID string
	JobAttributes  map[string]string
	ComamndString  string
	Environment    []string
	JobStatus      string
	ExecutionLog   string
	Cmd            exec.Cmd
}

// NewJob is...
func NewJob(jobID, comamndString string, environment []string, jobAttributes map[string]string) *Job {
	j := new(Job)
	j.JobID = jobID
	j.JobEcexutionID = jobID + uniuri.New()
	j.ComamndString = comamndString
	j.Environment = environment
	j.JobAttributes = jobAttributes
	return j
}

// Execute is...
func (j *Job) Execute(ctx context.Context) error {
	j.Cmd = *exec.Command("sh", "-c", j.ComamndString)
	j.Cmd.Env = append(os.Environ())
	j.Cmd.Env = append(j.Environment)
	j.Cmd.Stderr = os.Stderr
	j.Cmd.Stdout = os.Stdout

	// TODO: implement asyn
	j.Cmd.Start()
	// TODO: implement streaming log output
	j.Cmd.Wait()
	// TODO: implement bulk all log output

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
