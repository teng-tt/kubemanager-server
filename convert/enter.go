package convert

import (
	"kubmanager/convert/configmap"
	"kubmanager/convert/node"
	"kubmanager/convert/pod"
	"kubmanager/convert/secret"
)

type ConvertGroup struct {
	PodConvert       pod.PodConvertGroup
	NodeConver       node.NodeConvertGroup
	ConfigMapConvert configmap.ConvertGroup
	SecretConvert    secret.ConvertGroup
}

var ConvertGroupApp = new(ConvertGroup)
