package cmd

import (
	"context"
	"github.com/irotoris/jobkickqd/jobkickqd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

var jobTopicName string
var jobConfigFile string

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
}
