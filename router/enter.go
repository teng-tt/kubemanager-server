package router

import (
	"kubmanager/router/example"
	"kubmanager/router/harbor"
	"kubmanager/router/k8s"
	"kubmanager/router/metrics"
)

type RouterGroup struct {
	ExampleRouterGroup example.ExampleRouter
	K8SRouterGroup     k8s.K8SRouter
	HarborRouterGroup  harbor.HarborRouter
	MetricsRouterGroup metrics.MetricsRouter
}

var RouterGroupApp = new(RouterGroup)
