package main

import (
	"errors"
	"flag"
	"strings"
)

// Config represents app runtime configuration
type Config struct {
	Query                []string
	ProjectID            string
	MessageChannelBuffer int
	ResultChannelBuffer  int
}

func getConfig() (*Config, error) {
	projectID := flag.String("p", "", "GCP Project ID")
	queryFlag := flag.String("q", "", "Query args (e.g. golang, code, cloud)")
	messageChBuf := flag.Int("channel-message", 1, "Message channel buffer (default: 1)")
	resultChBuf := flag.Int("channel-result", 1, "Result channel buffer (default: 1)")
	flag.Parse()

	if len(*projectID) < 1 {
		return nil, errors.New("ProjectID argument required")
	}

	if len(*queryFlag) < 1 {
		return nil, errors.New("Query argument required")
	}

	return &Config{
		Query:                strings.Split(*queryFlag, ","),
		ProjectID:            *projectID,
		MessageChannelBuffer: *messageChBuf,
		ResultChannelBuffer:  *resultChBuf,
	}, nil
}
