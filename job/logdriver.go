package job

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
)

// LogDriver is ...
type LogDriver interface {
	Write()
}

// PubSubLogDriver is ...
type PubSubLogDriver struct {
	projectID            string
	topicName            string
	jsonCredentialString string
}

func (ld *PubSubLogDriver) Write(ctx context.Context, message string) error {
	client, err := pubsub.NewClient(ctx, ld.projectID)
	if err != nil {
		// TODO: Handle error.
	}
	topic := client.Topic(ld.topicName)
	defer topic.Stop()
	var results []*pubsub.PublishResult
	r := topic.Publish(ctx, &pubsub.Message{
		Data: []byte(message),
	})
	results = append(results, r)
	// Do other work ...
	for _, r := range results {
		id, err := r.Get(ctx)
		if err != nil {
			// TODO: Handle error.
		}
		fmt.Printf("Published a message with a message ID: %s\n", id)
	}
	return nil
}

// FileLogDriver is ...
type FileLogDriver struct {
	filename string
}

func (ld *FileLogDriver) Write(ctx context.Context, message string) error {
	return nil
}
