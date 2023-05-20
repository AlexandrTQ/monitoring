package monitoring

import (
	"monitoring/service/logging"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWorkerUpdateMap(t *testing.T) {
	logging.InitMockLogger()
	monitoringWorker.updateServerStatuses([]server{{Domain: "invalidHost"}, {Domain: "google.com"}, {Domain: "https://google.com"}})
	time.Sleep(time.Second * 6)
	invalidServer, ok := monitoringWorker.serverStatusMap["invalidHost"]
	if !ok {
		t.Fatalf("can't get by key \"invalidHost\"")
	}

	assert.Equal(t, invalidServer.Status, ServerStatusUnavailable)

	google, ok := monitoringWorker.serverStatusMap["google.com"]
	if !ok {
		t.Fatalf("can't get by key \"google.com\"")
	}

	googleWithHttps, ok := monitoringWorker.serverStatusMap["https://google.com"]
	if !ok {
		t.Fatalf("can't get by key \"https://google.com\"")
	}

	assert.Equal(t, google.Status, googleWithHttps.Status)

	if google.Status == ServerStatusAvailable {
		max, ok := monitoringWorker.serverStatusMap[keyForMaxResponseTime]
		if !ok {
			t.Fatalf("can't get by key max when google is available")
		}

		min, ok := monitoringWorker.serverStatusMap[keyForMinResponseTime]
		if !ok {
			t.Fatalf("can't get by key min when google is available")
		}

		assert.True(t, max.ResponseTimeMilliseconds >= min.ResponseTimeMilliseconds, "min timeout bigger then max")
	} else {
		_, ok := monitoringWorker.serverStatusMap[keyForMaxResponseTime]
		if !ok {
			t.Fatalf("get value by max when all services unavailable")
		}

		_, ok = monitoringWorker.serverStatusMap[keyForMinResponseTime]
		if !ok {
			t.Fatalf("get value by min when all services unavailable")
		}
	}
}
