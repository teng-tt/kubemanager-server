package deployment

import "kubemanager.com/convert"

type ServiceGroup struct {
	DeploymentService
}

var podConvert = convert.ConvertGroupApp.PodConvert
