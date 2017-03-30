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

	// configure publisher
	p, pErr := newPublisher()
	if pErr != nil {
		log.Fatal(pErr)
	}

	// configure ingester
	ingester := newIngester()
	iErr := ingester.start(searchTerms, messages)
	if iErr != nil {
		log.Fatal(iErr)
	}
	defer ingester.stop()

	go func() {
		for {
			select {
			case m := <-messages:
				p.pub(m)
			}
		}
	}()

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

}
