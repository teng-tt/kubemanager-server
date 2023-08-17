package k8s

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

type NodeMetrics struct {
	Metadata  metav1.ObjectMeta   `json:"metadata"`
	Timestamp time.Time           `json:"timestamp"`
	Window    string              `json:"window"`
	Usage     corev1.ResourceList `json:"usage"`
}

type NodeMetricsList struct {
	Kind       string            `json:"kind"`
	ApiVersion string            `json:"apiVersion"`
	Metadata   metav1.ObjectMeta `json:"metadata"`
	Items      []NodeMetrics     `json:"items"`
}
