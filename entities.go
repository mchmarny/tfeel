package main

import "fmt"

// Config represents app runtime configuration
type Config struct {
	Query     []string
	ProjectID string
}

// MiniTweet represents generic authored content
type MiniTweet struct {
	Query string `json:"query"`
	ID    string `json:"id"`
	On    string `json:"on"`
	By    string `json:"by"`
	Body  string `json:"body"`
}

func (m *MiniTweet) toString() string {
	return fmt.Sprintf("ID:%v, On:%v, By:%v, Body:%v", m.ID, m.On, m.By, m.Body)
}

// ProcessResult represents processed Message and it's content
type ProcessResult struct {
	SourceID         string  `json:"id"`
	SentimentScore   float32 `json:"score"`
	EntityAttributes string  `json:"parts"`
}
