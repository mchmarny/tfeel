package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"cloud.google.com/go/pubsub"

	"google.golang.org/api/option"
)

type psHelper struct {
	ctx    context.Context
	client *pubsub.Client
	topic  *pubsub.Topic
}

func newPSHelper() (*psHelper, error) {

	// get the Google CLoud Pub/Sub project ID
	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		return nil, errors.New("GCP Porject ID variable required (GCP_PROJECT_ID)")
	}

	topicName := os.Getenv("PUBSUB_TOPIC")
	if topicName == "" {
		topicName = "messages"
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
		fmt.Printf("Topic already exists: %v", topicErr)
	}

	pub := &psHelper{
		ctx:    ctx,
		client: client,
		topic:  topic,
	}

	fmt.Println("PubSub Helper configured")
	return pub, nil
}
