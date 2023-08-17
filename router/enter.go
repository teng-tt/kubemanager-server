package router

import (
	"kubemanager.com/router/example"
	"kubemanager.com/router/harbor"
	"kubemanager.com/router/k8s"
	"kubemanager.com/router/metrics"
)

type RouterGroup struct {
	ExampleRouterGroup example.ExampleRouter
	K8SRouterGroup     k8s.K8SRouter
	HarborRouterGroup  harbor.HarborRouter
	MetricsRouterGroup metrics.MetricsRouter
}

var RouterGroupApp = new(RouterGroup)
