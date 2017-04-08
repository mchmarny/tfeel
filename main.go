package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

/*
	Keys: https://apps.twitter.com/app/13608793/keys
*/

func main() {

	projectID := flag.String("p", "", "GCP Project ID")
	queryFlag := flag.String("q", "", "Query args (e.g. golang, code, cloud)")
	flag.Parse()

	if len(*projectID) < 1 {
		log.Fatal("ProjectID required")
	}

	if len(*queryFlag) < 1 {
		log.Fatal("Query required")
	}

	conf := &Config{
		Query:     strings.Split(*queryFlag, ","),
		ProjectID: *projectID,
	}

	// messages channel with some buffer
	messages := make(chan MiniTweet)
	results := make(chan ProcessResult)

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
