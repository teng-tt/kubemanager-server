package k8s

import (
	"kubmanager/convert"
	"kubmanager/validate"
)

type ApiGroup struct {
	PodApi
	NamespaceApi
}

var podValidate = validate.VaildateGroupApp.PodValidate
var podConvert = convert.ConvertGroupApp.PodConvert
