package monitoring

import (
	"fmt"
	"monitoring/service/util"
	"time"
)

func StartMonitoring(path string) error {
	content, err := util.GetFileAsLinesSlice(path)
	if err != nil {
		return fmt.Errorf("error on getting info from file - %s", err.Error())
	}

	serverList := getServersFromLines(content)

	if len(serverList) == 0 {
		return fmt.Errorf("can't get servers info from file")
	}

	err = monitoringWorker.startWork(serverList)
	if err != nil {
		return err
	}

	return nil
}

func StopMonitoring() {
	monitoringWorker.stopWork()
}

func GetInfoByDomain(domain string) (available bool, responseTimeMilliseconds int64, serverError string, lastUpdate time.Time, err error) {
	status, ok, lastUpdate := monitoringWorker.getServerInfo(domain)
	if !ok {
		return false, 0, "", time.Time{}, fmt.Errorf("can't find server by %s", domain)
	}

	if status.Status == ServerStatusAvailable {
		available = true
	}

	return available, status.ResponseTimeMilliseconds, status.Error, lastUpdate, nil
}

func GetByMaxResponseTime() (domain string, responseTimeMilliseconds int64, lastUpdate time.Time, err error) {
	status, ok, lastUpdate := monitoringWorker.getServerInfo(keyForMaxResponseTime)
	if !ok {
		return "", 0, time.Time{}, fmt.Errorf("can't find server by max response time. probably all services unavailable")
	}

	return status.Domain, status.ResponseTimeMilliseconds, lastUpdate, nil
}

func GetByMinResponseTime() (domain string, responseTimeMilliseconds int64, lastUpdate time.Time, err error) {
	status, ok, lastUpdate := monitoringWorker.getServerInfo(keyForMinResponseTime)
	if !ok {
		return "", 0, time.Time{}, fmt.Errorf("can't find server by min response time. probably all services unavailable")
	}

	return status.Domain, status.ResponseTimeMilliseconds, lastUpdate, nil
}
