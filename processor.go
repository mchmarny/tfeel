package main

import (
	"fmt"
	"os"

	"golang.org/x/net/context"

	"cloud.google.com/go/pubsub"
)

func process(ps *psHelper) {

	subName := os.Getenv("PUBSUB_SUBSCRIPTION_NAME")
	if subName == "" {
		subName = "events"
	}

	sub, err := ps.client.CreateSubscription(ps.ctx, subName, ps.topic, 0, nil)
	if err != nil {
		fmt.Printf("Subscription already exists: %v", err)
	}

	// TODO: atach to the kill envent
	waitCancel := false
	cctx, cancel := context.WithCancel(ps.ctx)
	subErr := sub.Receive(cctx, func(c context.Context, m *pubsub.Message) {
		content := string(m.Data)
		fmt.Printf("Body: %v", content)
		fmt.Printf("Ack: %v on %v", m.ID, m.PublishTime)
		if waitCancel {
			m.Nack()
			cancel()
		} else {
			m.Ack()
		}
	})

	if subErr != nil {
		fmt.Printf("Error on subcsription: %v", subErr)
	}

}
