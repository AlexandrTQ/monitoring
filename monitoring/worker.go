package monitoring

import (
	"context"
	"fmt"
	"monitoring/service/logging"
	"net/http"
	"strings"
	"sync"
	"time"
)

type worker struct {
	active          bool
	activeMutex     *sync.Mutex
	cancel          context.CancelFunc
	serverStatusMap map[string]serverStatus
	mapLastUpdate   time.Time
	mapMutex        *sync.RWMutex
}

var monitoringWorker = worker{activeMutex: &sync.Mutex{}, mapMutex: &sync.RWMutex{}, serverStatusMap: make(map[string]serverStatus)}

func (w *worker) startWork(serverList []server) error {
	if monitoringWorker.getActive() {
		return fmt.Errorf("worker already is active")
	}

	w.setActive(true)

	ctx, cancel := context.WithCancel(context.Background())
	w.cancel = cancel

	w.updateServerStatuses(serverList)

	go func() {
		ticker := time.NewTicker(60 * time.Second)
		for {
			select {
			case <-ticker.C:
				go w.updateServerStatuses(serverList)
			case <-ctx.Done():
				logging.Errorf("worker was stopped by context")
				return
			}
		}
	}()

	return nil
}

func (w *worker) updateServerStatuses(serverList []server) {
	newServerStatusesMap := make(map[string]serverStatus)

	semaphore := make(chan struct{}, 30)
	mutex := &sync.Mutex{}
	wg := &sync.WaitGroup{}

	httpClient := http.Client{Timeout: 5 * time.Second}

	start := time.Now()

	for _, server := range serverList {
		semaphore <- struct{}{}
		wg.Add(1)
		go func(domain string) {
			defer func() {
				wg.Done()
				<-semaphore
			}()

			url := domain
			if !strings.Contains(url, "http") {
				url = "https://" + url
			}

			result := serverStatus{Domain: domain}

			start := time.Now()
			response, err := httpClient.Get(url)
			responseTime := time.Since(start).Milliseconds()

			if err != nil {
				result.Status = ServerStatusUnavailable
				result.Error = err.Error()
			} else if response.StatusCode >= 500 {
				result.Status = ServerStatusUnavailable
				result.Error = fmt.Sprintf("server responded with statusCode %d", response.StatusCode)
			} else {
				result.Status = ServerStatusAvailable
				result.ResponseTimeMilliseconds = responseTime
			}

			mutex.Lock()
			newServerStatusesMap[domain] = result
			mutex.Unlock()
		}(server.Domain)
	}
	wg.Wait()

	var maxTimeoutStatus serverStatus
	var minTimeoutStatus serverStatus

	for _, serverStatus := range newServerStatusesMap {
		if serverStatus.Status == ServerStatusAvailable && (maxTimeoutStatus.Domain == "" || maxTimeoutStatus.ResponseTimeMilliseconds < serverStatus.ResponseTimeMilliseconds) {
			maxTimeoutStatus = serverStatus
		}

		if serverStatus.Status == ServerStatusAvailable && (minTimeoutStatus.Domain == "" || minTimeoutStatus.ResponseTimeMilliseconds > serverStatus.ResponseTimeMilliseconds) {
			minTimeoutStatus = serverStatus
		}
	}

	if maxTimeoutStatus.Domain != "" {
		newServerStatusesMap[keyForMaxResponseTime] = maxTimeoutStatus
	}
	if minTimeoutStatus.Domain != "" {
		newServerStatusesMap[keyForMinResponseTime] = minTimeoutStatus
	}

	w.mapMutex.Lock()
	w.serverStatusMap = newServerStatusesMap
	w.mapLastUpdate = time.Now()
	w.mapMutex.Unlock()

	logging.Debugf("ServerStatuses map was updated. take %.2fs.", time.Since(start).Seconds())
}

func (w *worker) stopWork() {
	if !w.getActive() {
		return
	}
	w.setActive(false)
	w.cancel()
}

func (w *worker) getServerInfo(key string) (*serverStatus, bool, time.Time) {
	w.mapMutex.RLock()
	defer w.mapMutex.RUnlock()

	result, ok := w.serverStatusMap[key]
	return &result, ok, w.mapLastUpdate
}

func (w *worker) setActive(active bool) {
	w.activeMutex.Lock()
	defer w.activeMutex.Unlock()
	w.active = active
}

func (w *worker) getActive() bool {
	w.activeMutex.Lock()
	defer w.activeMutex.Unlock()
	return w.active
}
