package pod

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	pod_req "kubmanager/model/pod/request"
)

type PodConvert struct {
}

// PodReq2K8s 将 pod 的请求格式的数据 转换为 k8s 结构数据
func (p *PodConvert) PodReq2K8s(podReq pod_req.Pod) *corev1.Pod {
	labels := podReq.Base.Labels
	k8sLabels := p.getK8sLabels(labels)
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podReq.Base.Name,
			Namespace: podReq.Base.NameSpace,
			Labels:    k8sLabels,
		},
		Spec: corev1.PodSpec{
			InitContainers: nil,
			Containers:     nil,
			Volumes:        nil,
			DNSConfig:      &corev1.PodDNSConfig{},
			DNSPolicy:      "",
			HostAliases:    nil,
			Hostname:       "",
			RestartPolicy:  "",
		},
	}
}

func (p *PodConvert) getK8sContainers(podReqContainers []pod_req.Container) []corev1.Container {
	return nil
}

// Pod 请求 labels 转换为 k8s labels
func (p *PodConvert) getK8sLabels(podReqLabels []pod_req.ListMapItem) map[string]string {
	podK8sLabels := make(map[string]string)
	for _, label := range podReqLabels {
		podK8sLabels[label.Key] = label.Value
	}

	return podK8sLabels
}
