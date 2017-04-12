package main

import (
	"bytes"
	"log"

	lang "cloud.google.com/go/language/apiv1"
	langpb "google.golang.org/genproto/googleapis/cloud/language/v1"
)

var (
	langClient *lang.Client
)

func initLangAPI() {
	client, err := lang.NewClient(appContext)
	if err != nil {
		log.Panicf("Failed to create client: %v", err)
	}
	langClient = client
}

func scoreSentiment(s string) (float32, error) {

	result, err := langClient.AnalyzeSentiment(appContext, &langpb.AnalyzeSentimentRequest{
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

func analyzeEntities(s string) (string, error) {
	result, err := langClient.AnalyzeEntities(appContext, &langpb.AnalyzeEntitiesRequest{
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
