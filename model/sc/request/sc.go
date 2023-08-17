package request

import (
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	"kubemanager.com/model/base"
)

type StorageClass struct {
	Name                 string                               `json:"name"`
	Namespace            string                               `json:"namespace"`
	Labels               []base.ListMapItem                   `json:"labels"`
	Provisioner          string                               `json:"provisioner"`          // 制备器
	MountOptions         []string                             `json:"mountOptions"`         // 卷绑定参数配置
	Parameters           []base.ListMapItem                   `json:"parameters"`           // 制备器入参
	ReclaimPolicy        corev1.PersistentVolumeReclaimPolicy `json:"reclaimPolicy"`        // 卷回收策略
	AllowVolumeExpansion bool                                 `json:"allowVolumeExpansion"` // 是否允许卷扩充
	VolumeBindIngMode    storagev1.VolumeBindingMode          `json:"volumeBindIngMode"`    // 卷绑定模式
}
