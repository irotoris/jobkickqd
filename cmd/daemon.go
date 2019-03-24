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
			daemonConfig.WorkDir = workDir
		}
		if daemonConfig.WorkDir != "" {
			if _, err := os.Stat(daemonConfig.WorkDir); os.IsNotExist(err) {
				err := os.Mkdir(daemonConfig.WorkDir, 0777)
				if err != nil {
					logrus.Fatalf("cannot create directory to %s: %v", daemonConfig.WorkDir, err)
				}
			}

			err := os.Chdir(daemonConfig.WorkDir)
			if err != nil {
				logrus.Fatalf("cannot change directory to %s: %v", daemonConfig.WorkDir, err)
			}
		}
		ctx := context.Background()
		l, err := jobkickqd.NewPubSubMessageDriver(ctx, daemonConfig.ProjectID, daemonConfig.LogTopic)
		if err != nil {
			logrus.Errorf("Failed to create a pubsub log driver.: %s", err)
		}

		q, err := jobkickqd.NewPubSubJobQueueExecutor(ctx, daemonConfig.ProjectID, daemonConfig.JobQueueTopic, daemonConfig.App, daemonConfig.App)
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
	daemonCmd.PersistentFlags().StringVar(&workDir, "workDir", "", "daemon work directory.")
}
