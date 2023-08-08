package k8s

import (
	"kubmanager/service"
	"kubmanager/validate"
)

type ApiGroup struct {
	PodApi
	NamespaceApi
	NodeApi
	ConfigMapApi
	SecretApi
	PVApi
	PVCApi
	SCApi
	SVCApi
	IngressApi
	IngressRouteApi
	StatefulSetApi
	DeploymentApi
	JoApi
}

var podValidate = validate.VaildateGroupApp.PodValidate
var podService = service.ServiceGroupApp.PodServiceGroup.PodService
var nodeService = service.ServiceGroupApp.NodeServiceGroup.NodeService
var configMapService = service.ServiceGroupApp.ConfigMapServiceGroup.ConfigMapService
var secretService = service.ServiceGroupApp.SecretServiceGroup.SecretService
var pvService = service.ServiceGroupApp.PvServiceGroup.PvService
var pvcService = service.ServiceGroupApp.PvcServiceGroup.PVCService
var scService = service.ServiceGroupApp.SCServiceGroup.SCService
var svcService = service.ServiceGroupApp.SvcServiceGroup.Service
var ingressService = service.ServiceGroupApp.IngressServiceGroup.IngresService
var ingrouteService = service.ServiceGroupApp.IngRouteServiceGroup.IngressRouteService
var statefulsetService = service.ServiceGroupApp.StatefulSetServiceGroup.StatefulSetService
var deploymentService = service.ServiceGroupApp.DeploymentService.DeploymentService
var jobService = service.ServiceGroupApp.JobService.JobService
