package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"

	"github.com/aws/aws-lambda-go/lambda"
)

type Event struct {
	Awslogs struct {
		Data string `json:"data"`
	} `json:"awslogs"`
}

type Message struct {
	MessageType         string   `json:"messageType"`
	Owner               string   `json:"owner"`
	LogGroup            string   `json:"logGroup"`
	LogStream           string   `json:"logStream"`
	SubscriptionFilters []string `json:"subscriptionFilters"`
	LogEvents           []struct {
		ID        string `json:"id"`
		Timestamp int64  `json:"timestamp"`
		Message   string `json:"message"`
	} `json:"logEvents"`
}

func decodeEvent(event *Event) (*Message, error) {
	data, err := base64.StdEncoding.DecodeString(event.Awslogs.Data)

	if err != nil {
		return nil, err
	}

	zr, err := gzip.NewReader(bytes.NewBuffer(data))

	if err != nil {
		return nil, err
	}

	defer zr.Close()

	buf := bytes.Buffer{}
	io.Copy(&buf, zr)

	msg := &Message{}
	err = json.Unmarshal(buf.Bytes(), msg)

	if err != nil {
		return nil, err
	}

	return msg, nil
}

func HandleRequest(ctx context.Context, event Event) (string, error) {
	msg, err := decodeEvent(&event)

	if err != nil {
		return err.Error(), err
	}

	logs := []*SqlLog{}

	for _, logEvent := range msg.LogEvents {
		sqlLog, err := parseLog(logEvent.Message)

		if err != nil {
			return err.Error(), err
		}

		if sqlLog != nil {
			logs = append(logs, sqlLog)
		}
	}

	if len(logs) > 0 {
		// dd := newDatadog()
		// err = dd.sendLogs(ctx, msg, logs)

		// if err != nil {
		// 	return err.Error(), err
		// }

		opnsrch, err := newOpenSearch()

		if err != nil {
			return err.Error(), err
		}

		err = opnsrch.postLogs(ctx, msg, logs)

		if err != nil {
			return err.Error(), err
		}
	}

	return "", nil
}

func main() {
	lambda.Start(HandleRequest)
}
