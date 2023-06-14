package k8s

import "kubmanager/validate"

type ApiGroup struct {
	PodApi
	NamespaceApi
}

var podValidate = validate.VaildateGroupApp.PodValidate
