package service

import (
	"kubmanager/service/configmap"
	"kubmanager/service/deployment"
	"kubmanager/service/ingress"
	"kubmanager/service/ingroute"
	"kubmanager/service/node"
	"kubmanager/service/pod"
	"kubmanager/service/pv"
	"kubmanager/service/pvc"
	"kubmanager/service/sc"
	"kubmanager/service/secret"
	"kubmanager/service/statefulset"
	"kubmanager/service/svc"
)

type ServiceGroup struct {
	PodServiceGroup         pod.PodServiceGroup
	NodeServiceGroup        node.NodeServiceGroup
	ConfigMapServiceGroup   configmap.ServiceGroup
	SecretServiceGroup      secret.ServicerGroup
	PvServiceGroup          pv.ServiceGroup
	PvcServiceGroup         pvc.ServiceGroup
	SCServiceGroup          sc.ServiceGroup
	SvcServiceGroup         svc.ServiceGroup
	IngressServiceGroup     ingress.ServiceGroup
	IngRouteServiceGroup    ingroute.ServiceGroup
	StatefulSetServiceGroup statefulset.ServiceGroup
	DeploymentService       deployment.ServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
