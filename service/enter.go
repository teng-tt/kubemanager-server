package service

import (
	"kubmanager/service/node"
	"kubmanager/service/pod"
)

type ServiceGroup struct {
	PodServiceGroup  pod.PodServiceGroup
	NodeServiceGroup node.NodeServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
