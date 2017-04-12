package main

import (
	"encoding/json"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
)

const (
	topicTweets     = "tweets"
	topicResults    = "results"
	subTweetsEvents = "tweets-events"
)

var (
	tweetTopic   *pubsub.Topic
	resultTopic  *pubsub.Topic
	pubsubClient *pubsub.Client
)

func initPublisher() {

	client, err := pubsub.NewClient(appContext, projectID)
	if err != nil {
		log.Panicf("Failed to create client: %v", err)
	}

	pubsubClient = client
	tweetTopic = client.Topic(topicTweets)
	resultTopic = client.Topic(topicResults)
}

func publish(t *pubsub.Topic, m interface{}) error {
	b, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("Error while marshaling object: %v", err)
	}

	msg := &pubsub.Message{Data: b}
	result := t.Publish(appContext, msg)
	id, err := result.Get(appContext)
	if err != nil {
		return fmt.Errorf("Error while publishing message: %v:%v", err, id)
	}

	//log.Printf("Published: %v", id)

	return nil
}
