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

type ingester struct {
	stream *twitter.Stream
}

func newIngester() *ingester {
	return &ingester{}
}

func (i *ingester) stop() {
	fmt.Println("Stopping Ingester...")
	if i.stream != nil {
		i.stream.Stop()
	}
}

func (i *ingester) start(s []string, ch chan<- Message) error {

	consumerKey := os.Getenv("T_CONSUMER_KEY")
	consumerSecret := os.Getenv("T_CONSUMER_SECRET")
	accessToken := os.Getenv("T_ACCESS_TOKEN")
	accessSecret := os.Getenv("T_ACCESS_SECRET")

	if consumerKey == "" || consumerSecret == "" || accessToken == "" || accessSecret == "" {
		return errors.New("Both, consumer key/secret and access token/secret are required")
	}

	// init convif
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)

	// HTTP Client - will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	demux := twitter.NewSwitchDemux()

	//Tweet processor
	demux.Tweet = func(tweet *twitter.Tweet) {
		msg := Message{
			ID:   tweet.IDStr,
			On:   tweet.CreatedAt,
			By:   tweet.User.ScreenName,
			Body: tweet.Text,
		}
		ch <- msg
		//fmt.Print(".")
	}

	// Tweet filter
	filterParams := &twitter.StreamFilterParams{
		Track:         s,
		StallWarnings: twitter.Bool(true),
		Language:      []string{"en"},
	}

	fmt.Println("Starting Ingest For: " + strings.Join(s, ","))

	// Start stream
	stream, err := client.Streams.Filter(filterParams)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// set local stream ref and go to work
	i.stream = stream
	go demux.HandleChan(stream.Messages)

	return nil

}
