package configmap

import "kubmanager/convert"

type ServiceGroup struct {
	ConfigMapService ConfigMapService
}

var configMapConvert = convert.ConvertGroupApp.ConfigMapConvert
