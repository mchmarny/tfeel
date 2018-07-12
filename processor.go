package main

import (
	"encoding/json"
	"log"

	"bytes"

	lang "cloud.google.com/go/language/apiv1"
	langpb "google.golang.org/genproto/googleapis/cloud/language/v1"

	"cloud.google.com/go/pubsub"

	"golang.org/x/net/context"
)

var (
	langClient *lang.Client
)

// ProcessResult represents processed Message and it's content
type ProcessResult struct {
	SourceID         string  `json:"id"`
	SentimentScore   float32 `json:"score"`
	EntityAttributes string  `json:"parts"`
}

func process(r chan<- ProcessResult) {

	// START INIT LANG API
	client, err := lang.NewClient(appContext)
	if err != nil {
		log.Panicf("Failed to create client: %v", err)
	}
	langClient = client
	// END INIT LANG API

	sub := pubsubClient.Subscription(subTweetsEvents)
	log.Printf("Subscribing: %v", sub)
	err = sub.Receive(appContext, func(c context.Context, m *pubsub.Message) {

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

		//log.Printf("P: %v", result.SourceID)

	})

	if err != nil {
		log.Printf("Error on subscription: %v", err)
	}

}

func scoreSentiment(s string) (float32, error) {

	result, err := langClient.AnalyzeSentiment(appContext, &langpb.AnalyzeSentimentRequest{
		Document: &langpb.Document{
			Source: &langpb.Document_Content{
				Content: s,
			},
			Type: langpb.Document_PLAIN_TEXT,
		},
		EncodingType: langpb.EncodingType_UTF8,
	})
	if err != nil {
		return 0, err
	}
	return result.DocumentSentiment.Score, nil
}

func analyzeEntities(s string) (string, error) {
	result, err := langClient.AnalyzeEntities(appContext, &langpb.AnalyzeEntitiesRequest{
		Document: &langpb.Document{
			Source: &langpb.Document_Content{
				Content: s,
			},
			Type: langpb.Document_PLAIN_TEXT,
		},
		EncodingType: langpb.EncodingType_UTF8,
	})
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	for _, e := range result.Entities {
		b.WriteString("(" + e.Name + ":" + e.Type.String() + ") ")
	}
	return b.String(), nil
}
