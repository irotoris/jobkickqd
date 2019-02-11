package main

import (
	"os"
	"github.com/irotoris/jobkickqd/cmd"
	"github.com/sirupsen/logrus"
)

var (
	// Version is ...
	Version string
	// Revision is ...
	Revision string
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {
	cmd.Execute()
}
