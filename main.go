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

	ingestor := newIngest()

	err := ingestor.start(searchTerms)
	if err != nil {
		log.Fatal(err)
	}

	defer ingestor.stop()

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

}
