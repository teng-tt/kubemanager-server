package convert

import (
	"kubmanager/convert/node"
	"kubmanager/convert/pod"
)

type ConvertGroup struct {
	PodConvert pod.PodConvertGroup
	NodeConver node.NodeConvertGroup
}

var ConvertGroupApp = new(ConvertGroup)
