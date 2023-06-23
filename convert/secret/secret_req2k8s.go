package secret

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	req_secret "kubmanager/model/secret/request"
	"kubmanager/utils"
)

type Req2K8s struct {
}

func (r *Req2K8s) SecretReq2K8s(secretReq req_secret.Secret) corev1.Secret {
	return corev1.Secret{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretReq.Name,
			Namespace: secretReq.Namespace,
			Labels:    utils.ToMap(secretReq.Labels),
		},
		Type:       secretReq.Type,
		StringData: utils.ToMap(secretReq.Data),
	}
}
