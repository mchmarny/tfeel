package main

import (
	"errors"
	"fmt"
	"os"

	"cloud.google.com/go/pubsub"
)

type publisher struct {
	client *pubsub.Client
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

	/*
		ctx := context.Background()
		client, err := pubsub.NewClient(ctx, projectID)
		if err != nil {
			return nil, fmt.Errorf("Failed to create client: %v", err)
		}
	*/

	pub := &publisher{}
	//pub.client = client
	return pub, nil
}

func (p *publisher) pub(m Message) {
	fmt.Println(m.toString())
}
