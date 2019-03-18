package cmd

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/irotoris/jobkickqd/jobkickqd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var jobTopicName string
var jobConfigFile string
var jobID string
var command string
var environmentString string
var timeoutInt int
var app string

// submitCmd represents the submit command
var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit a job command to a job queue.",
	Long:  `Submit a job command to a job queue.`,
	Run: func(cmd *cobra.Command, args []string) {
		exitCode, err := submit(args)
		if err != nil {
			logrus.Errorf("%s", err)
		}
		os.Exit(exitCode)
	},
}

func init() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
	rootCmd.AddCommand(submitCmd)
	submitCmd.PersistentFlags().StringVar(&projectID, "projectID", "", "GCP project name")
	submitCmd.PersistentFlags().StringVar(&jobTopicName, "jobTopicName", "", "Colud PubSub topic name for job queue")
	submitCmd.PersistentFlags().StringVar(&logTopicName, "logTopicName", "", "Colud PubSub topic name for job logs")
	submitCmd.PersistentFlags().StringVar(&jobConfigFile, "jobConfigFile", "", "Job config filename")
	submitCmd.PersistentFlags().StringVar(&jobID, "jobID", "", "Job ID")
	submitCmd.PersistentFlags().StringVar(&command, "command", "", "command")
	submitCmd.PersistentFlags().StringVar(&environmentString, "environment", "", "environment")
	submitCmd.PersistentFlags().IntVar(&timeoutInt, "timeout", 300, "timeout of command")
	submitCmd.PersistentFlags().StringVar(&app, "app", "default", "key of application")
}

func submit(args []string) (int, error) {
	var d []byte
	if jobConfigFile != "" {
		data, err := ioutil.ReadFile(jobConfigFile)
		if err != nil {
			logrus.Errorf("%s:Cannot open jobConfigFile:%s", err, jobConfigFile)
			return 1, err
		}
		var jm jobkickqd.DefaultJobMessage
		if err := json.Unmarshal(data, &jm); err != nil {
			logrus.Errorf("json.Unmarshal() failed.: %s", err)
			return 1, err
		}
		validParamFlag := true
		if jm.JobID == "" {
			logrus.Errorf("jobID is required in --jobConfigFile.")
			validParamFlag = false
		}
		if jm.Command == "" {
			logrus.Errorf("command is required in --jobConfigFile.")
			validParamFlag = false
		}
		if !validParamFlag {
			return 1, nil
		}
		jobID = jm.JobID
		command = jm.Command
		timeoutInt = jm.Timeout
		d = data
	} else {
		validParamFlag := true
		if jobID == "" {
			logrus.Errorf("--jobID is required if no --jobConfigFile.")
			validParamFlag = false
		}
		if command == "" {
			logrus.Errorf("--command is required if no --jobConfigFile.")
			validParamFlag = false
		}
		if !validParamFlag {
			return 1, nil
		}
		var envList []string
		if environmentString != "" {
			envList = strings.Split(environmentString, ",")
		}
		jobMessage := jobkickqd.DefaultJobMessage{JobID: jobID, Command: command, Environment: envList, Timeout: timeoutInt}
		data, err := json.Marshal(jobMessage)
		if err != nil {
			logrus.Errorf("JSON Marshal error in parse parameters:%s", err)
			return 1, err
		}
		d = data
	}

	if app == "" {
		logrus.Errorf("--app is required.")
		return 1, nil
	}

	ctx := context.Background()

	// Initialize to subscribe log messages
	pubsubClient, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		logrus.Errorf("%s", err)
		return 1, nil
	}
	topic := pubsubClient.Topic(logTopicName)
	sub, err := pubsubClient.CreateSubscription(ctx, jobID, pubsub.SubscriptionConfig{
		Topic:       topic,
		AckDeadline: 10 * time.Second,
	})
	if err != nil {
		logrus.Warnf("%s", err)
		sub = pubsubClient.Subscription(jobID)
	}
	defer sub.Delete(ctx)

	// Publish a job
	kickq, err := jobkickqd.NewPubSubMessageDriver(ctx, projectID, jobTopicName)
	if err != nil {
		logrus.Errorf("%s", err)
		return 1, err
	}

	attribute := map[string]string{
		"app": app,
	}

	id, err := kickq.Write(ctx, string(d), attribute)
	if err != nil {
		logrus.Errorf("%s", err)
		return 1, err
	}
	jobExecutionID := jobID + id

	// Start subscribe log messages

	// add interval 5 seconds for timeout
	cctx, cancel := context.WithTimeout(ctx, time.Duration((timeoutInt+5)*1)*time.Second)
	var jobExitCodeString string
	var mu sync.Mutex
	sub.Receive(cctx, func(ctx context.Context, m *pubsub.Message) {
		logrus.Debugf("message id:%s", m.ID)
		logrus.Debugf("message body:%s", string(m.Data))
		logrus.Debugf("message attr:%s", m.Attributes)
		m.Ack()
		if m.Attributes["job_execution_id"] != jobExecutionID {
			return
		}
		logrus.Infof("Job stdout/stderr:\n%s", string(m.Data))
		mu.Lock()
		defer mu.Unlock()
		jobExitCodeString = m.Attributes["job_exit_code"]
		cancel()
	})

	// waiting pubsub receive
	if jobExitCodeString == "" {
		logrus.Errorf("A command might be timeout...")
		return 1, nil
	} else if jobExitCodeString != "0" {
		logrus.Errorf("Job exit code: %s", jobExitCodeString)
		return 1, nil
	}
	return 0, nil

}
