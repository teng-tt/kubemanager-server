package configmap

import (
	corev1 "k8s.io/api/core/v1"
	conf_res "kubemanager.com/model/configmap/response"
	"kubemanager.com/utils"
)

type K8s2ResConvert struct {
}

func (k *K8s2ResConvert) GetCmReqItem(configMap corev1.ConfigMap) conf_res.ConfigMap {
	return conf_res.ConfigMap{
		Name:      configMap.Name,
		Namespace: configMap.Namespace,
		DataNum:   len(configMap.Data),
		Age:       configMap.CreationTimestamp.Unix(),
	}
}

func (k *K8s2ResConvert) GetCmReqDetail(configMap corev1.ConfigMap) conf_res.ConfigMap {
	detail := k.GetCmReqItem(configMap)
	detail.Data = utils.ToList(configMap.Data)
	detail.Labels = utils.ToList(configMap.Labels)

	return detail
}
