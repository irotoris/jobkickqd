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
	projectID string
	topicName string
	clinet    pubsub.Client
}

// NewPubSubLogDriver is ...
func NewPubSubLogDriver(ctx context.Context, projectID, topicName string) (*PubSubLogDriver, error) {
	ld := new(PubSubLogDriver)
	ld.projectID = projectID
	ld.topicName = topicName
	client, err := pubsub.NewClient(ctx, ld.projectID)
	if err != nil {
		return nil, err
	}
	ld.clinet = *client
	return ld, nil
}

func (ld *PubSubLogDriver) Write(ctx context.Context, message string, attributes map[string]string) error {
	client, err := pubsub.NewClient(ctx, ld.projectID)
	if err != nil {
		return err
	}
	topic := client.Topic(ld.topicName)
	defer topic.Stop()
	var results []*pubsub.PublishResult
	r := topic.Publish(ctx, &pubsub.Message{
		Data:       []byte(message),
		Attributes: attributes,
	})
	results = append(results, r)

	for _, r := range results {
		id, err := r.Get(ctx)
		if err != nil {
			return err
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
