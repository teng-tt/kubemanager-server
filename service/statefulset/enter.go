package statefulset

import "kubemanager.com/convert"

type ServiceGroup struct {
	StatefulSetService
}

var podConvert = convert.ConvertGroupApp.PodConvert
