package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

/*
	Keys: https://apps.twitter.com/app/13608793/keys
*/

type ingest struct {
	stream *twitter.Stream
}

func newIngest() *ingest {
	return &ingest{}
}

func (i *ingest) stop() {
	fmt.Println("Stopping Stream...")
	if i.stream != nil {
		i.stream.Stop()
	}
}

func (i *ingest) start(searchTerms []string) error {

	consumerKey := os.Getenv("T_CONSUMER_KEY")
	consumerSecret := os.Getenv("T_CONSUMER_SECRET")
	accessToken := os.Getenv("T_ACCESS_TOKEN")
	accessSecret := os.Getenv("T_ACCESS_SECRET")

	if consumerKey == "" || consumerSecret == "" || accessToken == "" || accessSecret == "" {
		return errors.New("Consumer key/secret and Access token/secret required")
	}

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)

	// HTTP Client - will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter Client
	client := twitter.NewClient(httpClient)

	// Convenience Demux demultiplexed stream messages
	demux := twitter.NewSwitchDemux()

	//Tweet
	demux.Tweet = func(tweet *twitter.Tweet) {
		msg := Message{
			ID:   tweet.IDStr,
			On:   tweet.CreatedAt,
			By:   tweet.User.ScreenName,
			Body: tweet.Text,
		}
		fmt.Println(msg.toString())
	}

	// Tweet filter
	filterParams := &twitter.StreamFilterParams{
		Track:         searchTerms,
		StallWarnings: twitter.Bool(true),
		Language:      []string{"en"},
	}

	fmt.Println("Starting Stream...")
	fmt.Println("Search Term: " + strings.Join(searchTerms, ","))

	// Start stream
	stream, err := client.Streams.Filter(filterParams)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// set local for stop
	i.stream = stream

	// go to work
	go demux.HandleChan(stream.Messages)

	return nil

}
