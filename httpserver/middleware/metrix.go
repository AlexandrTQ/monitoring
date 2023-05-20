package middleware

import (
	"net/http"
	"sync"
)

var statisticMap = make(map[string]int)

var m = &sync.RWMutex{}

func GetMetricValue(endpoint string) (int, bool) {
	m.RLock()
	defer m.RUnlock()
	value, ok := statisticMap[endpoint]
	return value, ok
}

func Metric(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		m.Lock()
		value, ok := statisticMap[r.URL.Path]
		if !ok {
			statisticMap[r.URL.Path] = 1
		} else {
			statisticMap[r.URL.Path] = value + 1
		}
		m.Unlock()

		next.ServeHTTP(rw, r)
	})
}
