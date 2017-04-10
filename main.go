package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

/*
	Keys: https://apps.twitter.com/app/13608793/keys
*/

func main() {

	conf, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}

	// messages channel with some buffer
	messages := make(chan MiniTweet, conf.MessageChannelBuffer)
	results := make(chan ProcessResult, conf.ResultChannelBuffer)

	// configure PubSub Helper
	ps, err := newPubSubHelper(conf.ProjectID)
	if err != nil {
		log.Fatal(err)
	}

	// Instantiate publisher
	pub := getPublisher(ps)

	// configure ingester
	ingester := newIngester()
	iErr := ingester.start(conf.Query, messages)
	if iErr != nil {
		log.Fatal(iErr)
	}
	defer ingester.stop()

	// start processing
	go process(ps, results)

	// wait for signal
	go func() {
		for {
			select {
			case m := <-messages:
				pub.publishTweet(m)
			case r := <-results:
				pub.publishResult(r)
			}
		}
	}()

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

}
