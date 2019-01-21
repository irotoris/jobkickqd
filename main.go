package main

import (
	"context"
	"flag"
	"github.com/irotoris/jobkickqd/jobkickqd"
	"github.com/sirupsen/logrus"
)

var (
	project      = flag.String("project", "", "Google Cloud Platform project name")
	topic        = flag.String("topic", "", "PubSub topic name for a log stream of jobs")
	subscription = flag.String("subscription", "", "PubSub subscription name for a job queue.")
)

func init() {
	flag.Parse()
}

func main() {
	ctx := context.Background()
	l, err := jobkickqd.NewPubSubMessageDriver(ctx, *project, *topic)
	if err != nil {
		logrus.Errorf("Failed to create a pubsub log driver.: %s", err)
	}

	q, err := jobkickqd.NewPubSubJobQueueExecutor(ctx, *project, *subscription)
	if err != nil {
		logrus.Errorf("Failed to create a pubsub job queue executor.: %s", err)
	}

	cctx, _ := context.WithCancel(ctx)

	// start job queue polling
	q.Run(ctx, cctx, *l)

}
