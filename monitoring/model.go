package monitoring

import (
	"strings"
)

type server struct {
	Domain string
}

func getServersFromLines(lines []string) []server {
	var result []server
	for _, line := range lines {
		domain := strings.TrimSpace(line)
		if domain != "" {
			result = append(result, server{Domain: domain})
		}
	}

	return result
}

type serverStatus struct {
	Domain                   string
	Status                   ServerStatusEnum
	ResponseTimeMilliseconds int64
	Error                    string
}

type ServerStatusEnum string

const ServerStatusAvailable ServerStatusEnum = "available"
const ServerStatusUnavailable ServerStatusEnum = "unavailable"

const (
	keyForMinResponseTime = "min"
	keyForMaxResponseTime = "max"
)
