package service

import "kubmanager/service/pod"

type ServiceGroup struct {
	PodServiceGroup pod.PodServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
