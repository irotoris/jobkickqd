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

// submitCmd represents the submit command
var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit a job command to a job queue.",
	Long:  `Submit a job command to a job queue.`,
	Run: func(cmd *cobra.Command, args []string) {
		var d []byte
		if jobConfigFile != "" {
			data, err := ioutil.ReadFile(jobConfigFile)
			if err != nil {
				logrus.Errorf("%s:Cannot open jobConfigFile:%s", err, jobConfigFile)
				os.Exit(1)
			}
			var jm jobkickqd.DefaultJobMessage
			if err := json.Unmarshal(data, &jm); err != nil {
				logrus.Errorf("json.Unmarshal() failed.: %s", err)
				os.Exit(1)
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
				os.Exit(1)
			}
			jobID = jm.JobID
			command = jm.Command
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
				os.Exit(1)
			}
			var envList []string
			if environmentString != "" {
				envList = strings.Split(environmentString, ",")
			}
			jobMessage := jobkickqd.DefaultJobMessage{JobID: jobID, Command: command, Environment: envList, Timeout: timeoutInt}
			data, err := json.Marshal(jobMessage)
			if err != nil {
				logrus.Errorf("JSON Marshal error in parse parameters:%s", err)
				os.Exit(1)
			}
			d = data
		}

		ctx := context.Background()

		// Initialize to subscribe log messages
		pubsubClient, err := pubsub.NewClient(ctx, projectID)
		if err != nil {
			logrus.Errorf("%s", err)
			os.Exit(1)
		}
		topic := pubsubClient.Topic(logTopicName)
		sub, err := pubsubClient.CreateSubscription(ctx, jobID, pubsub.SubscriptionConfig{
			Topic:       topic,
			AckDeadline: 10 * time.Second,
		})
		defer sub.Delete(ctx)
		if err != nil {
			logrus.Errorf("%s", err)
			os.Exit(1)
		}

		// Publish a job
		kickq, err := jobkickqd.NewPubSubMessageDriver(ctx, projectID, jobTopicName)
		if err != nil {
			logrus.Errorf("%s", err)
			os.Exit(1)
		}

		attribute := map[string]string{
			"app": "jobkickqd",
		}

		id, err := kickq.Write(ctx, string(d), attribute)
		if err != nil {
			logrus.Errorf("%s", err)
			os.Exit(1)
		}
		jobExecutionID := jobID + id

		// Start subscribe log messages
		cctx, cancel := context.WithCancel(ctx)
		var mu sync.Mutex
		err = sub.Receive(cctx, func(ctx context.Context, m *pubsub.Message) {
			if m.Attributes["job_execution_id"] != jobExecutionID {
				return
			}
			m.Ack()
			logrus.Infof(string(m.Data))
			mu.Lock()
			defer mu.Unlock()
			cancel()
		})
		if err != nil {
			logrus.Errorf("%s", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(submitCmd)
	submitCmd.PersistentFlags().StringVar(&projectID, "projectID", "", "GCP project name")
	submitCmd.PersistentFlags().StringVar(&jobTopicName, "jobTopicName", "", "Colud PubSub topic name for job queue")
	submitCmd.PersistentFlags().StringVar(&logTopicName, "logTopicName", "", "Colud PubSub topic name for job logs")
	submitCmd.PersistentFlags().StringVar(&jobConfigFile, "jobConfigFile", "", "Job config filename")
	submitCmd.PersistentFlags().StringVar(&jobID, "jobID", "", "Job ID")
	submitCmd.PersistentFlags().StringVar(&command, "command", "", "command")
	submitCmd.PersistentFlags().StringVar(&environmentString, "environment", "", "environment")
	submitCmd.PersistentFlags().IntVar(&timeoutInt, "timeout", 0, "timeout of command")
}
