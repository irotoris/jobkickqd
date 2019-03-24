package jobkickqd

// DefaultJobMessage is ...
type DefaultJobMessage struct {
	JobID          string   `json:"jobID"`
	JobExecutionID string   `json:"jobExecutionID"`
	Command        string   `json:"command"`
	Environment    []string `json:"environment"`
	Timeout        int      `json:"timeout"`
}

// DaemonConfig is ...
type DaemonConfig struct {
	ProjectID     string
	JobQueueTopic string
	LogTopic      string
	App           string
	WorkDir       string
}
