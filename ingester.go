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

// MiniTweet represents simple tweet content
type MiniTweet struct {
	Query string `json:"query"`
	ID    string `json:"id"`
	On    string `json:"on"`
	By    string `json:"by"`
	Body  string `json:"body"`
}

// toString returns readable string representation of the MiniTweet struct
func (m *MiniTweet) toString() string {
	return fmt.Sprintf("ID:%v, On:%v, By:%v, Body:%v", m.ID, m.On, m.By, m.Body)
}

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
	log.Println("Stopping Ingester...")
	if i.stream != nil {
		i.stream.Stop()
	}
}

// start initiates the Tweeter stream subscription and pumps all messages into
// the passed in channel
func (i *ingester) start(s []string, ch chan<- MiniTweet) error {

	consumerKey := os.Getenv("T_CONSUMER_KEY")
	consumerSecret := os.Getenv("T_CONSUMER_SECRET")
	accessToken := os.Getenv("T_ACCESS_TOKEN")
	accessSecret := os.Getenv("T_ACCESS_SECRET")

	if consumerKey == "" || consumerSecret == "" || accessToken == "" || accessSecret == "" {
		return errors.New("Both, consumer key/secret and access token/secret are required")
	}

	query := strings.Join(s, ",")

	// init convif
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)

	// HTTP Client - will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	demux := twitter.NewSwitchDemux()

	//Tweet processor
	demux.Tweet = func(tweet *twitter.Tweet) {
		msg := MiniTweet{
			Query: query,
			ID:    tweet.IDStr,
			On:    tweet.CreatedAt,
			By:    strings.ToLower(tweet.User.ScreenName),
			Body:  tweet.Text,
		}
		ch <- msg
	}

	// Tweet filter
	filterParams := &twitter.StreamFilterParams{
		Track:         s,
		StallWarnings: twitter.Bool(true),
		Language:      []string{"en"},
	}

	fmt.Printf("Starting Ingest For: " + strings.Join(s, ","))

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
