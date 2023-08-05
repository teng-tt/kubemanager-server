package deployment

import "kubmanager/convert"

type ServiceGroup struct {
	DeploymentService
}

var podConvert = convert.ConvertGroupApp.PodConvert
