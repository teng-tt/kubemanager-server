package response

import (
	corev1 "k8s.io/api/core/v1"
	"kubemanager.com/model/base"
)

type Node struct {
	Name             string             `json:"name"`
	Status           string             `json:"status"`
	Age              int64              `json:"age"`
	InternalIp       string             `json:"internalIp"`
	ExternalIp       string             `json:"externalIp"`
	Version          string             `json:"version"` // kubelet 版本
	OsImage          string             `json:"osImage"`
	KernelVersion    string             `json:"kernelVersion"`
	ContainerRuntime string             `json:"containerRuntime"` // 容器运行时
	Labels           []base.ListMapItem `json:"labels"`
	Taints           []corev1.Taint     `json:"taints"`
}
