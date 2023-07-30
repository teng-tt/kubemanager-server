package service

import (
	"kubmanager/service/configmap"
	"kubmanager/service/node"
	"kubmanager/service/pod"
	"kubmanager/service/pv"
	"kubmanager/service/pvc"
	"kubmanager/service/sc"
	"kubmanager/service/secret"
)

type ServiceGroup struct {
	PodServiceGroup       pod.PodServiceGroup
	NodeServiceGroup      node.NodeServiceGroup
	ConfigMapServiceGroup configmap.ServiceGroup
	SecretServiceGroup    secret.ServicerGroup
	PvServiceGroup        pv.ServiceGroup
	PvcServiceGroup       pvc.ServiceGroup
	SCServiceGroup        sc.ServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
