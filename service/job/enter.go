package job

import "kubmanager/convert"

type ServiceGroup struct {
	JobService
}

var podConvert = convert.ConvertGroupApp.PodConvert
