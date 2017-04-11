package main

import (
	"encoding/json"
	"log"

	"golang.org/x/net/context"

	"cloud.google.com/go/pubsub"
)

const (
	subTweetsEvents = "tweets-events"
)

// ProcessResult represents processed Message and it's content
type ProcessResult struct {
	SourceID         string  `json:"id"`
	SentimentScore   float32 `json:"score"`
	EntityAttributes string  `json:"parts"`
}

func process(ps *pubSubHelper, r chan<- ProcessResult) {

	sh, shErr := newSentimentHelper()
	if shErr != nil {
		log.Printf("Error while creating sentiment helper: %v", shErr)
	}

	// TODO: atach to the kill envent
	waitCancel := false
	cctx, cancel := context.WithCancel(ps.ctx)
	sub := ps.client.Subscription(subTweetsEvents)
	subErr := sub.Receive(cctx, func(c context.Context, m *pubsub.Message) {

		//fmt.Printf("Processing: %#v", m.ID)

		var mt MiniTweet
		if err := json.Unmarshal(m.Data, &mt); err != nil {
			log.Printf("Error while decoding tweet data: %#v", err)
			m.Nack()
			return
		}

		// p := Process{}
		//fmt.Printf("Body: %v", content)
		//fmt.Printf("Ack: %v on %v", m.ID, m.PublishTime)
		if waitCancel {
			m.Nack()
			cancel()
		} else {

			// score
			score, aErr := sh.scoreSentiment(mt.Body)
			if aErr != nil {
				log.Printf("Error while scoring: %v", aErr)
			}
			//fmt.Printf("ID:%v Score:%v\n", m.ID, score)

			// analyze
			args, eErr := sh.analyzeEntities(mt.Body)
			if eErr != nil {
				log.Printf("Error while analyzing: %v", eErr)
			}
			//fmt.Printf("ID:%v Attributes:%v\n", m.ID, args)

			result := ProcessResult{
				SourceID:         mt.ID,
				SentimentScore:   score,
				EntityAttributes: args,
			}

			r <- result

			m.Ack()
		}
	})

	if subErr != nil {
		log.Printf("Error on subcsription: %v", subErr)
	}

}
