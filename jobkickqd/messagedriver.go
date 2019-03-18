package jobkickqd

import (
	"context"

	"github.com/sirupsen/logrus"

	"cloud.google.com/go/pubsub"
)

// MessageDriver is ...
type MessageDriver interface {
	Write()
}

// PubSubMessageDriver is ...
type PubSubMessageDriver struct {
	projectID string
	topicName string
	client    pubsub.Client
}

// NewPubSubMessageDriver is ...
func NewPubSubMessageDriver(ctx context.Context, projectID, topicName string) (*PubSubMessageDriver, error) {
	ld := new(PubSubMessageDriver)
	ld.projectID = projectID
	ld.topicName = topicName
	client, err := pubsub.NewClient(ctx, ld.projectID)
	if err != nil {
		return nil, err
	}
	ld.client = *client
	return ld, nil
}

func (ld *PubSubMessageDriver) Write(ctx context.Context, message string, attributes map[string]string) (string, error) {
	topic := ld.client.Topic(ld.topicName)
	defer topic.Stop()
	r := topic.Publish(ctx, &pubsub.Message{
		Data:       []byte(message),
		Attributes: attributes,
	})

	id, err := r.Get(ctx)
	if err != nil {
		logrus.Errorf("Failed to publish a message: %s", err)
		return "", err
	}
	logrus.Infof("Published a message to pubsub[%s] with a message ID: %s", ld.topicName, id)

	return id, nil
}
