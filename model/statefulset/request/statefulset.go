package request

import (
	"kubemanager.com/model/base"
	podReq "kubemanager.com/model/pod/request"
	pvcReq "kubemanager.com/model/pvc/request"
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
