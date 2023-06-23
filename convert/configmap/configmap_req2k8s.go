package configmap

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	conf_req "kubmanager/model/configmap/request"
	"kubmanager/utils"
)

type Req2K8sConvert struct {
}

func (r *Req2K8sConvert) CmReq2K8sConvert(configMap conf_req.ConfigMap) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMap.Name,
			Namespace: configMap.Namespace,
			Labels:    utils.ToMap(configMap.Labels),
		},
		Data: utils.ToMap(configMap.Data),
	}
}
