package private

import (
	"monitoring/service/util"
	"net/http"

	"github.com/go-chi/chi"
)

const (
	PrivateRouterPathPrefix = "/private"
)

var r *chi.Mux

func GetPrivateRouter() *chi.Mux {
	return r
}

func init() {
	r = chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Get("/metric", getMetric)
	})
}

var s service

func getMetric(rw http.ResponseWriter, r *http.Request) {
	counterByDomain, counterByMax, counterByMin := s.getMetric()

	resp := struct {
		CounterByDomain int `json:"counterByDomain"`
		CounterByMax    int `json:"counterByMaxResponseTime"`
		CounterByMin    int `json:"counterByMinResponseTime"`
	}{
		CounterByDomain: counterByDomain,
		CounterByMax:    counterByMax,
		CounterByMin:    counterByMin,
	}

	util.WriteResult(rw, &resp, http.StatusOK)
}
