package convert

import "kubmanager/convert/pod"

type ConvertGroup struct {
	PodConvert pod.PodConvertGroup
}

var ConvertGroupApp = new(ConvertGroup)
