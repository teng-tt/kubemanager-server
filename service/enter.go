package service

import (
	"kubmanager/service/configmap"
	"kubmanager/service/node"
	"kubmanager/service/pod"
	"kubmanager/service/secret"
)

type ServiceGroup struct {
	PodServiceGroup       pod.PodServiceGroup
	NodeServiceGroup      node.NodeServiceGroup
	ConfigMapServiceGroup configmap.ServiceGroup
	SecretServiceGroup    secret.ServicerGroup
}

var ServiceGroupApp = new(ServiceGroup)
