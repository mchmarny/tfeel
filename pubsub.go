package main

import (
	"encoding/json"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"

	"golang.org/x/net/context"
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

// ProcessResult represents processed Message and it's content
type ProcessResult struct {
	SourceID         string  `json:"id"`
	SentimentScore   float32 `json:"score"`
	EntityAttributes string  `json:"parts"`
}

func initPubSub() {

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

func process(r chan<- ProcessResult) {

	initLangAPI()

	sub := pubsubClient.Subscription(subTweetsEvents)
	log.Printf("Subscribing: %v", sub)
	err := sub.Receive(appContext, func(c context.Context, m *pubsub.Message) {

		if canceling {
			log.Print("Cancel pending, skipping...")
			m.Nack()
			return
		}

		var mt MiniTweet
		if err := json.Unmarshal(m.Data, &mt); err != nil {
			log.Printf("Error while decoding tweet data: %#v", err)
			m.Nack()
			return
		}

		// score
		score, aErr := scoreSentiment(mt.Body)
		if aErr != nil {
			log.Printf("Error while scoring: %v", aErr)
		}

		// analyze
		args, eErr := analyzeEntities(mt.Body)
		if eErr != nil {
			log.Printf("Error while analyzing: %v", eErr)
		}

		result := ProcessResult{
			SourceID:         mt.ID,
			SentimentScore:   score,
			EntityAttributes: args,
		}

		r <- result

		m.Ack()

		log.Printf("P: %v", result.SourceID)

	})

	if err != nil {
		log.Printf("Error on subcsription: %v", err)
	}

}
