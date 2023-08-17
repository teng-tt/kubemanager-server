package api

import (
	"kubmanager/api/example"
	"kubmanager/api/harbor"
	"kubmanager/api/k8s"
	"kubmanager/api/metrics"
)

type ApiGroup struct {
	ExampleApiGroup example.ApiGroup
	K8SApiGroup     k8s.ApiGroup
	HarborApiGroup  harbor.ApiGroup
	MetricsApiGroup metrics.ApiGroup
}

var ApiGroupApp = new(ApiGroup)
