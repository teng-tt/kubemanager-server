package configmap

import "kubemanager.com/convert"

type ServiceGroup struct {
	ConfigMapService ConfigMapService
}

var configMapConvert = convert.ConvertGroupApp.ConfigMapConvert
