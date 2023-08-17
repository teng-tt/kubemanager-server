package pod

import "kubemanager.com/convert"

type PodServiceGroup struct {
	PodService PodService
}

var podConvert = convert.ConvertGroupApp.PodConvert
