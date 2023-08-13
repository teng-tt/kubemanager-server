package cronjob

import "kubmanager/convert"

type ServiceGroup struct {
	CronjobService
}

var podConvert = convert.ConvertGroupApp.PodConvert
