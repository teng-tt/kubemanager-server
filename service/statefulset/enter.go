package statefulset

import "kubmanager/convert"

type ServiceGroup struct {
	StatefulSetService
}

var podConvert = convert.ConvertGroupApp.PodConvert
