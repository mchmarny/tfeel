package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"google.golang.org/api/option"

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

	optArgs, errgErr := getClientOptions()
	if errgErr != nil {
		return nil, fmt.Errorf("Error while creating client Args: %v", errgErr)
	}
	fmt.Printf("DEBUG: %v", optArgs)
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

func getClientOptions() (option.ClientOption, error) {

	if _, err := os.Stat(serviceAccountKeyFilePath); err != nil {
		return nil, errors.New("Service account file path not set")
	}

	fmt.Printf("Acquiring service account config using %v", serviceAccountKeyFilePath)
	optArgs := option.WithServiceAccountFile(serviceAccountKeyFilePath)

	return optArgs, nil

}
