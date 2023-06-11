package api

import (
	"kubmanager/api/example"
	"kubmanager/api/k8s"
)

type ApiGroup struct {
	ExampleApiGroup example.ApiGroup
	K8SApiGroup     k8s.ApiGroup
}

var ApiGroupApp = new(ApiGroup)
