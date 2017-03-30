package main

import "fmt"

// Message represents generic authored content
type Message struct {
	ID   string `json:"id"`
	On   string `json:"on"`
	By   string `json:"by"`
	Body string `json:"body"`
}

func (m *Message) toString() string {
	return fmt.Sprintf("ID:%v, On:%v, By:%v, Body:%v", m.ID, m.On, m.By, m.Body)
}
