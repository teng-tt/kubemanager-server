package request

import (
	corev1 "k8s.io/api/core/v1"
	"kubemanager.com/model/base"
)

type PersistentVolumeClaim struct {
	Name             string                              `json:"name"`
	Namespace        string                              `json:"namespace"`
	Labels           []base.ListMapItem                  `json:"labels"`
	AccessModes      []corev1.PersistentVolumeAccessMode `json:"accessModes"`
	Capacity         int32                               `json:"capacity"`
	Selector         []base.ListMapItem                  `json:"selector"`
	StorageClassName string                              `json:"storageClassName"`
}
