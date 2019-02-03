package cmd

import (
	"context"

	"github.com/irotoris/jobkickqd/jobkickqd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var jobSubscriptionName string
var logTopicName string

// daemonCmd represents start a daemon command
var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Start commands polling.",
	Long:  `Start commands polling.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		l, err := jobkickqd.NewPubSubMessageDriver(ctx, projectID, logTopicName)
		if err != nil {
			logrus.Errorf("Failed to create a pubsub log driver.: %s", err)
		}

		q, err := jobkickqd.NewPubSubJobQueueExecutor(ctx, projectID, jobSubscriptionName)
		if err != nil {
			logrus.Errorf("Failed to create a pubsub job queue executor.: %s", err)
		}

		cctx, _ := context.WithCancel(ctx)

		// start job queue polling
		q.Run(ctx, cctx, *l)

	},
}

func init() {
	rootCmd.AddCommand(daemonCmd)
	daemonCmd.PersistentFlags().StringVar(&projectID, "projectID", "", "GCP project name")
	daemonCmd.PersistentFlags().StringVar(&jobSubscriptionName, "jobSubscriptionName", "", "Colud PubSub subscription name for job queue")
	daemonCmd.PersistentFlags().StringVar(&logTopicName, "logTopicName", "", "Colud PubSub topic name for log stream")
}