package main

import (
	"errors"
	"fmt"

	"cloud.google.com/go/pubsub"
)

type publisher struct {
	ps *psHelper
}

// publish publishes messahes using precreated client
func (p *publisher) publish(m Message) error {

	if p.ps.topic == nil {
		return errors.New("Null topic error")
	}

	b := []byte(m.Body)
	attrs := map[string]string{"ID": m.ID, "By": m.By, "On": m.On}
	c := &pubsub.Message{Data: b, Attributes: attrs}

	if _, err := p.ps.topic.Publish(p.ps.ctx, c).Get(p.ps.ctx); err != nil {
		return fmt.Errorf("Error while publishing message: %v", err)
	}

	return nil

}
