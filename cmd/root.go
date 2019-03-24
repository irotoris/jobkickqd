package cmd

import (
	"fmt"
	"os"

	"github.com/irotoris/jobkickqd/jobkickqd"

	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string
var daemonConfig jobkickqd.DaemonConfig
var verbose bool
var projectID string
var jobQueueTopic string
var logTopic string
var app string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "jobkickqd",
	Short: "A command kicker via job queue(Cloud Pub/Sub) and this client.",
	Long:  ``,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if verbose {
			logrus.SetLevel(logrus.DebugLevel)
		} else {
			logrus.SetLevel(logrus.InfoLevel)
		}
	},
}

// Execute is...
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Errorf("%v", err)
		os.Exit(1)
	}
}

func init() {

	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "daemon config file")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "log option")
	rootCmd.PersistentFlags().StringVar(&projectID, "projectID", "", "GCP project name")
	rootCmd.PersistentFlags().StringVar(&jobQueueTopic, "jobQueueTopic", "", "Colud PubSub topic name for job queue")
	rootCmd.PersistentFlags().StringVar(&logTopic, "logTopic", "", "Colud PubSub topic name for log stream")
	rootCmd.PersistentFlags().StringVar(&app, "app", "", "key of application of daemon.")
	cobra.OnInitialize(initConfig)
	if projectID != "" {
		daemonConfig.ProjectID = projectID
	}
	if jobQueueTopic != "" {
		daemonConfig.JobQueueTopic = jobQueueTopic
	}
	if logTopic != "" {
		daemonConfig.LogTopic = logTopic
	}
	if app != "" {
		daemonConfig.App = app
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if configFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(configFile)

		if err := viper.ReadInConfig(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		viper.AutomaticEnv() // read in environment variables that match
		if err := viper.Unmarshal(&daemonConfig); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
