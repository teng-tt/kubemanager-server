package request

import (
	corev1 "k8s.io/api/core/v1"
	"kubmanager/model/base"
)

type NfsVolumeSource struct {
	NfsPath     string `json:"nfsPath"`
	NfsServer   string `json:"nfsServer"`
	NfsReadOnly bool   `json:"nfsReadOnly"`
}

type VolumeSource struct {
	Type            string          `json:"type"`
	NfsVolumeSource NfsVolumeSource `json:"nfsVolumeSource"`
}

type PersistentVolume struct {
	Name string `json:"name"`
	// ns 不必传
	// Namespace     string                              `json:"namespace"`
	Labels        []base.ListMapItem                   `json:"labels"`
	Capacity      int32                                `json:"capacity"`      // pv的容量
	AccessModes   []corev1.PersistentVolumeAccessMode  `json:"accessModes"`   // 数据读写权限
	ReClaimPolicy corev1.PersistentVolumeReclaimPolicy `json:"reClaimPolicy"` // pv回收策略
	VolumeSource  VolumeSource                         `json:"volumeSource"`
}
