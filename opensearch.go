package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"

	opensearch "github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

type OpenSearch struct {
	Client *opensearch.Client
}

func newOpenSearch() (*OpenSearch, error) {
	client, err := opensearch.NewClient(opensearch.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Addresses: []string{"https://" + os.Getenv("OPENSEARCH_ENDPOINT")},
	})

	if err != nil {
		return nil, err
	}

	return &OpenSearch{
		Client: client,
	}, nil
}

func (opnsrch *OpenSearch) postLogs(ctx context.Context, msg *Message, logs []*SqlLog) error {
	document := strings.NewReader(`{
		"title": "Moneyball",
		"director": "Bennett Miller",
		"year": "2011"
}`)

	req := opensearchapi.IndexRequest{
		Index: "test",
		Body:  document,
	}

	res, err := req.Do(ctx, opnsrch.Client)

	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", res)

	return nil
}
