package k8s

import (
	"kubmanager/service"
	"kubmanager/validate"
)

type ApiGroup struct {
	PodApi
	NamespaceApi
}

var podValidate = validate.VaildateGroupApp.PodValidate
var podService = service.ServiceGroupApp.PodServiceGroup.PodService
