package response

import (
	corev1 "k8s.io/api/core/v1"
	"kubmanager/model/base"
)

type PersistentVolume struct {
	Name             string                               `json:"name"`
	Capacity         int32                                `json:"capacity"` // pv容量
	Labels           []base.ListMapItem                   `json:"labels"`
	AccessModes      []corev1.PersistentVolumeAccessMode  `json:"accessModes"`      // 数据读写权限
	ReClaimPolicy    corev1.PersistentVolumeReclaimPolicy `json:"reClaimPolicy"`    // pv回收策略
	Status           corev1.PersistentVolumePhase         `json:"status"`           // 待完善
	Claim            string                               `json:"claim"`            // 被某个具体的pvc绑定
	Age              int64                                `json:"age"`              // 创建时间戳
	Reason           string                               `json:"reason"`           // 状况描述
	StorageClassName string                               `json:"storageClassName"` // sc名称
}
