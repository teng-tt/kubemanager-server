package job

import "kubemanager.com/convert"

type ServiceGroup struct {
	JobService
}

var podConvert = convert.ConvertGroupApp.PodConvert
