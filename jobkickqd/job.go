package jobkickqd

import (
	"context"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/dchest/uniuri"
)

// Job is...
type Job struct {
	JobID           string
	JobEcexutionID  string
	ComamndString   string
	Environment     []string
	JobState        string
	ExecutionLog    string
	SubimitedAt     time.Time
	StartedAt       time.Time
	FinishedAt      time.Time
	TimeoutDuration time.Duration
	Cmd             exec.Cmd
}

// NewJob is...
func NewJob(jobID, comamndString string, environment []string, timeoutDuration time.Duration) *Job {
	j := new(Job)
	j.JobID = jobID
	j.JobEcexutionID = jobID + "-" + uniuri.New()
	j.ComamndString = comamndString
	j.Environment = environment
	j.SubimitedAt = time.Now()
	j.TimeoutDuration = timeoutDuration
	return j
}

// Execute is...
func (j *Job) Execute(ctx context.Context) error {
	j.Cmd = *exec.Command("sh", "-c", j.ComamndString)
	j.Cmd.Env = append(os.Environ())
	j.Cmd.Env = append(j.Environment)
	// TODO: via log driver
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs/", 0755)
	}
	logFilename := "logs/" + j.JobEcexutionID + ".log"
	logFile, err := os.OpenFile(logFilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer logFile.Close()
	j.Cmd.Stderr = logFile
	j.Cmd.Stdout = logFile
	j.StartedAt = time.Now()

	j.Cmd.Start()
	j.JobState = "RUNNING"

	// TODO: implement streaming log output. finaly put end mark log.
	// TODO: implement update job state to datastore or other KVS.
	// TODO: implement stop commands when daemon process stop.(Or this responsibility is queue daemon.)
	j.Cmd.Wait()
	j.FinishedAt = time.Now()
	logFile.Close()

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
