package service

import (
	"kubemanager.com/service/configmap"
	"kubemanager.com/service/cronjob"
	"kubemanager.com/service/deployment"
	"kubemanager.com/service/ingress"
	"kubemanager.com/service/ingroute"
	"kubemanager.com/service/job"
	"kubemanager.com/service/metrics"
	"kubemanager.com/service/node"
	"kubemanager.com/service/pod"
	"kubemanager.com/service/pv"
	"kubemanager.com/service/pvc"
	"kubemanager.com/service/rbac"
	"kubemanager.com/service/sc"
	"kubemanager.com/service/secret"
	"kubemanager.com/service/statefulset"
	"kubemanager.com/service/svc"
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
	DeploymentServiceGroup  deployment.ServiceGroup
	JobServiceGroup         job.ServiceGroup
	CronJobServiceGroup     cronjob.ServiceGroup
	RbacServiceGroup        rbac.ServiceGroup
	MetricsServiceGroup     metrics.ServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
