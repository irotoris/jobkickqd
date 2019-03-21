package jobkickqd

// DefaultJobMessage is ...
type DefaultJobMessage struct {
	JobID          string   `json:"job_id"`
	JobExecutionID string   `json:"job_execution_id"`
	Command        string   `json:"command"`
	Environment    []string `json:"environment"`
	Timeout        int      `json:"timeout"`
}
