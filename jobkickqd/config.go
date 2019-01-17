package jobkickqd

import "time"

type RunnerConfig struct {
	Loglevel              string
	Logfile               string
	LogDriversConfig      map[string]string
	JobQueueDriversConfig map[string]string
	Concurency            int
}

type ClientConfig struct {
	Loglevel              string
	LogPollingInverval    time.Duration
	LogDriversConfig      map[string]string
	JobQueueDriversConfig map[string]string
}
