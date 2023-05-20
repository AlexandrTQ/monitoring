package public

import (
	"monitoring/service/monitoring"
	"time"
)

type service struct{}

func (service) getServerInfo(domain string) (bool, int64, string, time.Time, error) {
	return monitoring.GetInfoByDomain(domain)
}

func (service) getServerInfoWithMaxResponseTime() (string, int64, time.Time, error) {
	return monitoring.GetByMaxResponseTime()
}

func (service) getServerInfoWithMinResponseTime() (string, int64, time.Time, error) {
	return monitoring.GetByMinResponseTime()
}
