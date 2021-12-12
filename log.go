package main

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/percona/go-mysql/query"
)

type SqlLog struct {
	Timestamp    string
	RemoteHost   string
	UserDatabase string
	ProcessId    string
	ErrorLevel   string
	Duration     float64
	Statement    string
	Fingerprint  string
}

var rePrefix = regexp.MustCompile(`(?s)^(\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2}\s+[^:]+):([^:]*):([^:]*):([^:]*):([^:]*):(.*)`)
var reLog = regexp.MustCompile(`(?s)^\s+duration:\s+(\d+\.\d+)\s+ms\s+(?:statement|execute\s+[^:]+):(.*)`)

func parseLog(line string) (*SqlLog, error) {
	prefixMatches := rePrefix.FindStringSubmatch(line)

	if prefixMatches == nil {
		return nil, nil
	}

	errorLevel := string(prefixMatches[5])
	log := prefixMatches[6]

	if errorLevel != "LOG" {
		return nil, nil
	}

	logMatches := reLog.FindStringSubmatch(log)

	if logMatches == nil {
		return nil, nil
	}

	durationStr := string(logMatches[1])
	duration, err := strconv.ParseFloat(durationStr, 64)

	if err != nil {
		return nil, err
	}

	stmt := string(logMatches[2])

	return &SqlLog{
		Timestamp:    string(prefixMatches[1]),
		RemoteHost:   string(prefixMatches[2]),
		UserDatabase: string(prefixMatches[3]),
		ProcessId:    string(prefixMatches[4]),
		ErrorLevel:   string(prefixMatches[5]),
		Duration:     duration,
		Statement:    stmt,
		Fingerprint:  query.Fingerprint(strings.ReplaceAll(stmt, `"`, "")),
	}, nil
}
