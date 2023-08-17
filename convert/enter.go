package convert

import (
	"kubemanager.com/convert/configmap"
	"kubemanager.com/convert/node"
	"kubemanager.com/convert/pod"
	"kubemanager.com/convert/secret"
)

type ConvertGroup struct {
	PodConvert       pod.PodConvertGroup
	NodeConver       node.NodeConvertGroup
	ConfigMapConvert configmap.ConvertGroup
	SecretConvert    secret.ConvertGroup
}

var ConvertGroupApp = new(ConvertGroup)
