package main

import (
	"bytes"
	"log"

	lang "cloud.google.com/go/language/apiv1"
	"golang.org/x/net/context"
	langpb "google.golang.org/genproto/googleapis/cloud/language/v1"
)

type sentimentHelper struct {
	ctx    context.Context
	client *lang.Client
}

func newSentimentHelper() (*sentimentHelper, error) {
	ctx := context.Background()
	client, err := lang.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	return &sentimentHelper{
		ctx:    ctx,
		client: client,
	}, nil

}

func (h *sentimentHelper) scoreSentiment(s string) (float32, error) {

	result, err := h.client.AnalyzeSentiment(h.ctx, &langpb.AnalyzeSentimentRequest{
		Document: &langpb.Document{
			Source: &langpb.Document_Content{
				Content: s,
			},
			Type: langpb.Document_PLAIN_TEXT,
		},
		EncodingType: langpb.EncodingType_UTF8,
	})
	if err != nil {
		return 0, err
	}

	return result.DocumentSentiment.Score, nil

}

func (h *sentimentHelper) analyzeEntities(s string) (string, error) {

	result, err := h.client.AnalyzeEntities(h.ctx, &langpb.AnalyzeEntitiesRequest{
		Document: &langpb.Document{
			Source: &langpb.Document_Content{
				Content: s,
			},
			Type: langpb.Document_PLAIN_TEXT,
		},
		EncodingType: langpb.EncodingType_UTF8,
	})
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	for _, e := range result.Entities {
		b.WriteString("(" + e.Name + ":" + e.Type.String() + ") ")
	}

	return b.String(), nil

}
