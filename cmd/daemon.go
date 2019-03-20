package cmd

import (
	"context"
	"os"

	"github.com/irotoris/jobkickqd/jobkickqd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var workDir string

// daemonCmd represents start a daemon command
var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Start commands polling.",
	Long:  `Start commands polling.`,
	Run: func(cmd *cobra.Command, args []string) {
		if workDir != "" {
			if _, err := os.Stat(workDir); os.IsNotExist(err) {
				err := os.Mkdir(workDir, 0777)
				if err != nil {
					logrus.Fatalf("cannot create directory to %s: %v", workDir, err)
				}
			}

			err := os.Chdir(workDir)
			if err != nil {
				logrus.Fatalf("cannot change directory to %s: %v", workDir, err)
			}
		}
		ctx := context.Background()
		l, err := jobkickqd.NewPubSubMessageDriver(ctx, projectID, logTopic)
		if err != nil {
			logrus.Errorf("Failed to create a pubsub log driver.: %s", err)
		}

		q, err := jobkickqd.NewPubSubJobQueueExecutor(ctx, projectID, jobQueueTopic, app, app)
		if err != nil {
			logrus.Errorf("Failed to create a pubsub job queue executor.: %s", err)
		}

		cctx, _ := context.WithCancel(ctx)

		// start job queue polling
		q.Run(ctx, cctx, *l)

	},
}

func init() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
	rootCmd.AddCommand(daemonCmd)
	daemonCmd.PersistentFlags().StringVar(&projectID, "projectID", "", "GCP project name")
	daemonCmd.PersistentFlags().StringVar(&jobQueueTopic, "jobQueueTopic", "", "Colud PubSub topic name for job queue")
	daemonCmd.PersistentFlags().StringVar(&logTopic, "logTopic", "", "Colud PubSub topic name for log stream")
	daemonCmd.PersistentFlags().StringVar(&app, "app", "", "key of application of daemon.")
	daemonCmd.PersistentFlags().StringVar(&workDir, "workDir", "", "daemon work directory.")


}
