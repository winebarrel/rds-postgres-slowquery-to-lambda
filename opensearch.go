package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	opensearch "github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

type OpenSearch struct {
	Client *opensearch.Client
}

type Doc struct {
	Fingerprint string  `json:"fingerprint"`
	Duration    float64 `json:"duration"`
	Timestamp   string  `json:"@timestamp"`
}

func newOpenSearch() (*OpenSearch, error) {
	client, err := opensearch.NewClient(opensearch.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Addresses: []string{fmt.Sprintf("https://%s:443", os.Getenv("OPENSEARCH_ENDPOINT"))},
	})

	if err != nil {
		return nil, err
	}

	return &OpenSearch{
		Client: client,
	}, nil
}

func (opnsrch *OpenSearch) postLogs(ctx context.Context, msg *Message, logs []*SqlLog) error {
	for _, log := range logs {
		doc := &Doc{
			Fingerprint: log.Fingerprint,
			Duration:    log.Duration,
			Timestamp:   time.Now().Format(time.RFC3339),
		}

		rawJson, err := json.Marshal(doc)

		if err != nil {
			return err
		}

		req := opensearchapi.IndexRequest{
			Index: "slowquery",
			Body:  bytes.NewBuffer(rawJson),
		}

		res, err := req.Do(ctx, opnsrch.Client)

		if err != nil {
			return err
		}

		fmt.Printf("%+v\n", res)
	}

	return nil
}
