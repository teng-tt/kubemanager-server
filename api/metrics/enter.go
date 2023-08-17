package metrics

import "kubemanager.com/service"

type ApiGroup struct {
	MetricsApi
	PrometheusApi
}

var metricsService = service.ServiceGroupApp.MetricsServiceGroup.MetricsService
