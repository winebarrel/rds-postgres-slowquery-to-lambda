package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/DataDog/datadog-api-client-go/api/v2/datadog"
)

type Datadog struct {
	Source         string
	Hostname       string
	Service        string
	APIClient      *datadog.APIClient
	OptionalParams datadog.SubmitLogOptionalParameters
}

type LogMessage struct {
	Fingerprint string  `json:"fingerprint"`
	Duration    float64 `json:"duration"`
}

func newDatadog() *Datadog {
	funcName := os.Getenv("AWS_LAMBDA_FUNCTION_NAME")
	contentEncoding := datadog.ContentEncoding("gzip")
	ddtags := fmt.Sprintf("version:%s", os.Getenv("AWS_LAMBDA_FUNCTION_VERSION"))

	if v, ok := os.LookupEnv("DD_TAGS"); ok {
		ddtags += "," + v
	}

	dd := &Datadog{
		Source:    funcName,
		Hostname:  funcName,
		Service:   funcName,
		APIClient: datadog.NewAPIClient(datadog.NewConfiguration()),
		OptionalParams: datadog.SubmitLogOptionalParameters{
			ContentEncoding: &contentEncoding,
			Ddtags:          &ddtags,
		},
	}

	if v, ok := os.LookupEnv("DD_SOURCE"); ok {
		dd.Source = v
	}

	if v, ok := os.LookupEnv("DD_HOSTNAME"); ok {
		dd.Hostname = v
	}

	if v, ok := os.LookupEnv("DD_SERVICE"); ok {
		dd.Service = v
	}

	return dd
}

func (dd *Datadog) sendLogs(ctx context.Context, msg *Message, logs []*SqlLog) error {
	ctx = datadog.NewDefaultContext(ctx)
	body := []datadog.HTTPLogItem{}

	for _, log := range logs {
		logMsg := &LogMessage{
			Fingerprint: log.Fingerprint,
			Duration:    log.Duration,
		}

		rawJson, err := json.Marshal(logMsg)

		if err != nil {
			return err
		}

		jsonStr := string(rawJson)
		fmt.Println(jsonStr)

		item := datadog.HTTPLogItem{
			Ddsource: &dd.Source,
			Hostname: &dd.Hostname,
			Service:  &dd.Service,
			Message:  &jsonStr,
		}

		body = append(body, item)
	}

	ctx = datadog.NewDefaultContext(ctx)
	_, _, err := dd.APIClient.LogsApi.SubmitLog(ctx, body, dd.OptionalParams)

	if err != nil {
		return err
	}

	return nil
}
