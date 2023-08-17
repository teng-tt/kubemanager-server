package metrics

import "kubmanager/service"

type ApiGroup struct {
	MetricsApi
	PrometheusApi
}

var metricsService = service.ServiceGroupApp.MetricsServiceGroup.MetricsService
