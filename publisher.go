package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

type publisher struct {
	ctx    context.Context
	client *pubsub.Client
	topic  *pubsub.Topic
}

func newPublisher() (*publisher, error) {

	// get the Google CLoud Pub/Sub project ID
	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		return nil, errors.New("GCP Porject ID variable required (GCP_PROJECT_ID)")
	}

	topicName := os.Getenv("PUBSUB_TOPIC")
	if topicName == "" {
		return nil, errors.New("GCP PubSub topic name variable required (PUBSUB_TOPIC)")
	}

	saFilePath := os.Getenv("SERVICE_ACCOUNT_FILE_PATH")
	if saFilePath == "" {
		return nil, errors.New("Service account file path not set")
	}

	fmt.Println("Acquiring service account config...")
	optArgs := option.WithServiceAccountFile(saFilePath)

	fmt.Println("Creating context...")
	ctx := context.Background()

	fmt.Println("Creating pubsub client...")
	client, err := pubsub.NewClient(ctx, projectID, optArgs)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client: %v", err)
	}

	fmt.Println("Creating pubsub topic...")
	topic, topicErr := client.CreateTopic(ctx, topicName)
	if topicErr != nil {
		fmt.Println("Topic already exists")
	}

	pub := &publisher{
		ctx:    ctx,
		client: client,
		topic:  topic,
	}

	fmt.Println("Publisher configured")
	return pub, nil
}

// pub publishes messahes using precreated client
func (p *publisher) pub(m Message) error {

	if p.topic == nil {
		return errors.New("Null topic error")
	}

	b, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("Failed to marhal message: %v", err)
	}

	_, err = p.topic.Publish(p.ctx, &pubsub.Message{Data: b}).Get(p.ctx)
	if err != nil {
		return fmt.Errorf("Error while publishing message: %v", err)
	}

	return nil

}
