package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"golang.org/x/net/context"
)

var (
	appContext context.Context
	projectID  string
	canceling  bool
)

func main() {

	// START CONFIG
	defaultProjectID := os.Getenv("GCLOUD_PROJECT")
	projectFlag := flag.String("projectID", defaultProjectID, "GCP Project ID")
	queryFlag := flag.String("query", "", "Query args (e.g. golang, code, cloud)")
	flag.Parse()

	if len(*projectFlag) < 1 {
		log.Panic("ProjectID argument required")
	}

	if len(*queryFlag) < 1 {
		log.Panic("Query argument required")
	}
	queryArgs := strings.Split(*queryFlag, ",")
	projectID = *projectFlag
	// END CONFIG

	ctx, cancel := context.WithCancel(context.Background())
	appContext = ctx
	go func() {
		// Wait for SIGINT and SIGTERM (HIT CTRL-C)
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		log.Println(<-ch)
		canceling = true
		cancel()
		os.Exit(0)
	}()

	// messages channel with some buffer
	messages := make(chan MiniTweet)
	results := make(chan ProcessResult)

	// initialize PubSub client
	initPubSub()

	// start processing
	go func() {
		process(results)
	}()

	// configure ingester
	ingester := newIngester()
	err := ingester.start(queryArgs, messages)
	if err != nil {
		log.Fatal(err)
	}
	defer ingester.stop()

	for {
		select {
		case <-appContext.Done():
			break
		case m := <-messages:
			publish(tweetTopic, m)
		case r := <-results:
			publish(resultTopic, r)
		}
	}
}
