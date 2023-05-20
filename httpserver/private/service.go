package private

import (
	"monitoring/service/httpserver/middleware"
	"monitoring/service/httpserver/public"
	"path"
)

type service struct{}

func (service) getMetric() (int, int, int) {
	byDomainCount, ok := middleware.GetMetricValue(path.Join(public.PublicRouterPathPrefix, public.ByDomainEndpoint))
	if !ok {
		byDomainCount = 0
	}

	byMaxTimeCount, ok := middleware.GetMetricValue(path.Join(public.PublicRouterPathPrefix, public.ByMaxResponseTimeEndpoint))
	if !ok {
		byMaxTimeCount = 0
	}

	byMixTimeCount, ok := middleware.GetMetricValue(path.Join(public.PublicRouterPathPrefix, public.ByMinResponseTimeEndpoint))
	if !ok {
		byMixTimeCount = 0
	}

	return byDomainCount, byMaxTimeCount, byMixTimeCount
}
