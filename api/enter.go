package api

import (
	"kubemanager.com/api/example"
	"kubemanager.com/api/harbor"
	"kubemanager.com/api/k8s"
	"kubemanager.com/api/metrics"
)

type ApiGroup struct {
	ExampleApiGroup example.ApiGroup
	K8SApiGroup     k8s.ApiGroup
	HarborApiGroup  harbor.ApiGroup
	MetricsApiGroup metrics.ApiGroup
}

var ApiGroupApp = new(ApiGroup)
