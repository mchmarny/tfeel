package main

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
)

const (
	topicTweets  = "tweets"
	topicResults = "results"
)

type publisher struct {
	ps *pubSubHelper
	tt *pubsub.Topic
	tr *pubsub.Topic
}

func getPublisher(ps *pubSubHelper) *publisher {

	return &publisher{
		ps: ps,
		tt: ps.client.Topic(topicTweets),
		tr: ps.client.Topic(topicResults),
	}

}

// publish publishes messahes using precreated client
func (p *publisher) publishResult(m ProcessResult) error {
	return publish(p.ps.ctx, p.tr, m)
}

func (p *publisher) publishTweet(m MiniTweet) error {
	return publish(p.ps.ctx, p.tt, m)
}

func publish(ctx context.Context, t *pubsub.Topic, m interface{}) error {

	b, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("Error while marshaling object: %v", err)
	}

	msg := &pubsub.Message{Data: b}
	result := t.Publish(ctx, msg)
	id, err := result.Get(ctx)
	if err != nil {
		return fmt.Errorf("Error while publishing message: %v:%v", err, id)
	}

	//fmt.Printf("Published %v ID:%v\n", t.String(), id)
	return nil

}
