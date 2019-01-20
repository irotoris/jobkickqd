package jobkickqd

import (
	"context"
	"github.com/dchest/uniuri"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

// Job is...
type Job struct {
	JobID          string
	JobExecutionID string
	CommandString  string
	Environment    []string
	JobState       string
	ExecutionLog   string
	SentAt         time.Time
	SubmittedAt    time.Time
	StartedAt      time.Time
	FinishedAt     time.Time
	Timeout        time.Duration
	Cmd            exec.Cmd
}

// NewJob is...
func NewJob(jobID, CommandString string, environment []string, timeout time.Duration) *Job {
	j := new(Job)
	j.JobID = jobID
	j.JobExecutionID = jobID + "-" + uniuri.New()
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

	logFilename := "logs/" + j.JobExecutionID + ".log"
	logFile, err := NewFileLogDriver(logFilename)
	if err != nil {
		return err
	}
	defer logFile.Close(ctx)
	j.Cmd.Stderr = &logFile.file
	j.Cmd.Stdout = &logFile.file
	j.StartedAt = time.Now()

	j.Cmd.Start()
	j.JobState = "RUNNING"

	// TODO: implement streaming log output. finally put end mark log.
	// TODO: implement update job state to datastore or other KVS.
	// TODO: implement stop commands when daemon process stop.(Or this responsibility is queue daemon.)
	// TODO: implement timeout job cancel
	// TODO: implement retry in fail
	j.Cmd.Wait()
	j.FinishedAt = time.Now()

	// TODO: implement bulk all log output
	data, err := ioutil.ReadFile(logFilename)
	if err != nil {
		j.ExecutionLog = "[jobkickqd][daemon]ERROR:Cannot open a log file." + err.Error()
	} else {
		j.ExecutionLog = string(data)
	}
	j.ExecutionLog = string(data)
	j.changeJobStateAtEnd(ctx)
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

// changeJobStateAtEnd is...
func (j *Job) changeJobStateAtEnd(ctx context.Context) {
	state := j.Cmd.ProcessState
	if state.Exited() && state.Success() {
		j.JobState = "SUCCEEDED"
	} else if state.Exited() && !state.Success() {
		j.JobState = "FAILED"
	}
}
