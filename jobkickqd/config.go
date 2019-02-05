package jobkickqd

import "time"

type RunnerConfig struct {
	LogLevel         string
	Logfile          string
	LogDriversConfig map[string]string
	JobQueueConfig   map[string]string
	Concurrency       int
}

type ClientConfig struct {
	LogLevel           string
	LogPollingInterval time.Duration
	LogDriversConfig   map[string]string
	JobQueueConfig     map[string]string
}

type DefaultJobMessage struct {
	JobID          string   `json:"job_id"`
	JobExecutionID string   `json:"job_execution_id"`
	Command        string   `json:"command"`
	Environment    []string `json:"environment"`
	Timeout        int      `json:"timeout"`
	Retry          int      `json:"retry"`
}
