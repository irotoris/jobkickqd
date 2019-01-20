package jobkickqd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

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

// PubSubMessageDriver is ...
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

func (ld *PubSubMessageDriver) Write(ctx context.Context, message string, attributes map[string]string) error {
	topic := ld.client.Topic(ld.topicName)
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
			fmt.Printf("err:%s", err)
			return err
		}
		fmt.Printf("Published a message with a message ID: %s\n", id)
	}
	return nil
}

// FileMessageDriver is ...
type FileMessageDriver struct {
	filePath string
	file     os.File
}

// NewFileMessageDriver is ///
func NewFileMessageDriver(filePath string) (*FileMessageDriver, error) {
	ld := new(FileMessageDriver)
	ld.filePath = filePath
	logDirectory, _ := filepath.Split(filePath)
	err := os.MkdirAll(logDirectory, 0755)
	if err != nil {
		return nil, err
	}
	logFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	ld.file = *logFile

	return ld, nil
}

func (ld *FileMessageDriver) Write(ctx context.Context, message string) error {
	if _, err := ld.file.Write(([]byte)(message)); err != err {
		return err
	}
	return nil
}

// Close is ...
func (ld *FileMessageDriver) Close(ctx context.Context) error {
	if err := ld.file.Close(); err != nil {
		return err
	}
	return nil
}
