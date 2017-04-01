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

	// TODO: externalize
	searchTerms := []string{"google", "cloud", "data", "partner"}

	// messages channel with some buffer
	// TODO: externalize channel size
	messages := make(chan Message, 10)

	// configure PubSub Helper
	ps, err := newPSHelper()
	if err != nil {
		log.Fatal(err)
	}

	pub := &publisher{ps: ps}

	// configure ingester
	ingester := newIngester()
	iErr := ingester.start(searchTerms, messages)
	if iErr != nil {
		log.Fatal(iErr)
	}
	defer ingester.stop()

	go process(ps)

	go func() {
		for {
			select {
			case m := <-messages:
				pub.publish(m)
			}
		}
	}()

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

}
