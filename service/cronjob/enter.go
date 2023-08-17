package cronjob

import "kubemanager.com/convert"

type ServiceGroup struct {
	CronjobService
}

var podConvert = convert.ConvertGroupApp.PodConvert
