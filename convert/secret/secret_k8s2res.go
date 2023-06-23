package secret

import (
	corev1 "k8s.io/api/core/v1"
	res_secret "kubmanager/model/secret/response"
	"kubmanager/utils"
)

type K8s2Res struct {
}

func (k *K8s2Res) SecretK8s2ResItemConvert(secretK8s corev1.Secret) res_secret.Secret {
	return res_secret.Secret{
		Name:      secretK8s.Name,
		Namespace: secretK8s.Namespace,
		Type:      secretK8s.Type,
		DataNum:   len(secretK8s.Data),
		Age:       secretK8s.CreationTimestamp.Unix(),
	}
}

func (k *K8s2Res) SecretK8s2ResDetailConvert(secretK8s corev1.Secret) res_secret.Secret {
	detail := k.SecretK8s2ResItemConvert(secretK8s)
	detail.Labels = utils.ToList(secretK8s.Labels)
	detail.Data = utils.ToListWithMapByte(secretK8s.Data)

	return detail
}
