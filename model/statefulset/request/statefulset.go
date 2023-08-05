package request

import (
	"kubmanager/model/base"
	podReq "kubmanager/model/pod/request"
	pvcReq "kubmanager/model/pvc/request"
)

type StatefulSetBase struct {
	Name                 string                         `json:"name"`
	Namespace            string                         `json:"namespace"`
	Replicas             int32                          `json:"replicas"`
	Labels               []base.ListMapItem             `json:"labels"`
	Selector             []base.ListMapItem             `json:"selector"`
	ServiceName          string                         `json:"serviceName"`
	VolumeClaimTemplates []pvcReq.PersistentVolumeClaim `json:"volumeClaimTemplates"`
}

type StatefulSte struct {
	Base     StatefulSetBase `json:"base"`
	Template podReq.Pod      `json:"template"`
}
