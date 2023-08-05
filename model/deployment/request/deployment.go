package request

import (
	"kubmanager/model/base"
	podReq "kubmanager/model/pod/request"
)

type DeploymentBase struct {
	Name      string             `json:"name"`
	Namespace string             `json:"namespace"`
	Replicas  int32              `json:"replicas"`
	Labels    []base.ListMapItem `json:"labels"`
	Selector  []base.ListMapItem `json:"selector"`
}

type Deployment struct {
	Base     DeploymentBase `json:"base"`
	Template podReq.Pod
}
