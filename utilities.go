package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
)

const (
	subTweetEvents            = "tweets-events"
	serviceAccountKeyFilePath = "service-account-key.json"
)

type pubSubHelper struct {
	ctx    context.Context
	client *pubsub.Client
}

func newPubSubHelper(projectID string) (*pubSubHelper, error) {

	// ProjectID
	if len(projectID) < 1 {
		return nil, fmt.Errorf("ProjectID required: %v", projectID)
	}

	fmt.Println("Creating pubsub client...")
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client: %v", err)
	}

	pub := &pubSubHelper{
		ctx:    ctx,
		client: client,
	}

	return pub, nil
}
