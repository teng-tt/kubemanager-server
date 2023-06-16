package pod

import "kubmanager/convert"

type PodServiceGroup struct {
	PodService PodService
}

var podConvert = convert.ConvertGroupApp.PodConvert
