package public

import (
	"fmt"
	"monitoring/service/util"
	"net/http"
	"net/url"

	"github.com/go-chi/chi"
)

const (
	PublicRouterPathPrefix    = "/api/v1"
	ByDomainEndpoint          = "/server/status"
	ByMaxResponseTimeEndpoint = "/server/status/timeout/max"
	ByMinResponseTimeEndpoint = "/server/status/timeout/min"
)

var r *chi.Mux

func GetPublicRouter() *chi.Mux {
	return r
}

func init() {
	r = chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Get(ByDomainEndpoint, getServerStatus)
		r.Get(ByMaxResponseTimeEndpoint, getServerStatusByMaxResponseTime)
		r.Get(ByMinResponseTimeEndpoint, getServerStatusByMinResponseTime)
	})
}

var s service

func getServerStatus(rw http.ResponseWriter, r *http.Request) {
	queryParams, _ := url.ParseQuery(r.URL.RawQuery)
	domain := queryParams.Get("domain")

	if domain == "" {
		util.WriteResult(rw, fmt.Errorf("domain is empty"), http.StatusBadRequest)
		return
	}

	available, respTime, serverError, lastUpdate, err := s.getServerInfo(domain)
	if err != nil {
		util.WriteResult(rw, err, http.StatusBadRequest)
		return
	}

	resp := serverStatusDto{
		Domain:     domain,
		Available:  available,
		LastUpdate: lastUpdate.UTC().Format("02-01-2006 15:04:05 UTC"),
	}

	if available {
		resp.ResponseTime = fmt.Sprintf("%.3fs", float64(respTime)/1000)
	} else {
		resp.ServerError = serverError
	}

	util.WriteResult(rw, &resp, http.StatusOK)
}

func getServerStatusByMaxResponseTime(rw http.ResponseWriter, r *http.Request) {
	domain, respTime, lastUpdate, err := s.getServerInfoWithMaxResponseTime()
	if err != nil {
		util.WriteResult(rw, err, http.StatusNotFound)
		return
	}

	resp := serverStatusDto{
		Domain:       domain,
		Available:    true,
		ResponseTime: fmt.Sprintf("%.3fs", float64(respTime)/1000),
		LastUpdate:   lastUpdate.UTC().Format("02-01-2006 15:04:05 UTC"),
	}

	util.WriteResult(rw, &resp, http.StatusOK)
}

func getServerStatusByMinResponseTime(rw http.ResponseWriter, r *http.Request) {
	domain, respTime, lastUpdate, err := s.getServerInfoWithMinResponseTime()
	if err != nil {
		util.WriteResult(rw, err, http.StatusNotFound)
		return
	}

	resp := serverStatusDto{
		Domain:       domain,
		Available:    true,
		ResponseTime: fmt.Sprintf("%.3fs", float64(respTime)/1000),
		LastUpdate:   lastUpdate.UTC().Format("02-01-2006 15:04:05 UTC"),
	}

	util.WriteResult(rw, &resp, http.StatusOK)
}
