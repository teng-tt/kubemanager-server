package router

import (
	"kubmanager/router/example"
	"kubmanager/router/k8s"
)

type RouterGroup struct {
	ExampleRouterGroup example.ExampleRouter
	K8SRouterGroup     k8s.K8SRouter
}

var RouterGroupApp = new(RouterGroup)
