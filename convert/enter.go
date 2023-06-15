package convert

import "kubmanager/convert/pod"

type ConvertGroup struct {
	PodConvert pod.PodConvert
}

var ConvertGroupApp = new(ConvertGroup)
