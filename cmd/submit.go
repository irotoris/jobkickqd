package cmd

import (
	"context"
	"github.com/irotoris/jobkickqd/jobkickqd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"time"
)

var jobTopicName string
var jobConfigFile string
var jobID string
var command string
var environmentString string
var retry int
var timeout time.Duration

// submitCmd represents the submit command
var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit a job command to a job queue.",
	Long:  `Submit a job command to a job queue.`,
	Run: func(cmd *cobra.Command, args []string) {
		data, err := ioutil.ReadFile(jobConfigFile)
		if err != nil {
			logrus.Errorf("%s:Cannot open jobConfigFile:%s", err, jobConfigFile)
			os.Exit(1)
		}

		// publish a job
		ctx := context.Background()
		kickq, err := jobkickqd.NewPubSubMessageDriver(ctx, projectID, jobTopicName)
		if err != nil {
			logrus.Errorf("%s", err)
			os.Exit(1)
		}

		attribute := map[string]string{
			"app": "jobkickqd",
		}
		err = kickq.Write(ctx, string(data), attribute)
		if err != nil {
			logrus.Errorf("%s", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(submitCmd)
	submitCmd.PersistentFlags().StringVar(&projectID, "projectID", "", "GCP project name")
	submitCmd.PersistentFlags().StringVar(&jobTopicName, "jobTopicName", "", "Colud PubSub topic name")
	submitCmd.PersistentFlags().StringVar(&jobConfigFile, "jobConfigFile", "", "Job config filename")
	submitCmd.PersistentFlags().StringVar(&jobID, "jobID", "", "Job ID")
	submitCmd.PersistentFlags().StringVar(&command, "command", "", "command")
	submitCmd.PersistentFlags().StringVar(&environmentString, "environment", "", "environment")
	submitCmd.PersistentFlags().IntVar(&retry, "retry", 0, "retry count of command")
	submitCmd.PersistentFlags().IntVar(&timeout, "timeout", 0, "timeout of command")
}
