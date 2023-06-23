package secret

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubmanager/global"
	req_secret "kubmanager/model/secret/request"
	res_secret "kubmanager/model/secret/response"
	"strings"
)

type SecretService struct {
}

func (s *SecretService) CreateOrUpdateSecret(secretReq req_secret.Secret) error {
	ctx := context.TODO()
	secretK8s := secretConvert.Req2K8sConvert.SecretReq2K8s(secretReq)
	k8sApi := global.KubeConfigSet.CoreV1().Secrets(secretK8s.Namespace)

	// 查询是否存在
	_, err := k8sApi.Get(ctx,
		secretK8s.Name, metav1.GetOptions{})
	// 存在更新
	if err == nil {
		_, err = k8sApi.Update(ctx, &secretK8s, metav1.UpdateOptions{})
		return err
	}
	// 创建
	_, err = k8sApi.Create(ctx,
		&secretK8s, metav1.CreateOptions{})

	return err
}

func (s *SecretService) DeleteSecret(namespace, name string) error {
	return global.KubeConfigSet.CoreV1().Secrets(namespace).Delete(context.TODO(),
		name, metav1.DeleteOptions{})
}

func (s *SecretService) GetSecretDetail(namespace, name string) (*res_secret.Secret, error) {
	secretK8s, err := global.KubeConfigSet.CoreV1().Secrets(namespace).Get(context.TODO(),
		name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	secretRes := secretConvert.K8sResConvert.SecretK8s2ResDetailConvert(*secretK8s)
	return &secretRes, err
}

func (s *SecretService) GetSecretList(namespace, keyword string) ([]res_secret.Secret, error) {
	list, err := global.KubeConfigSet.CoreV1().Secrets(namespace).List(context.TODO(),
		metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	secretResList := make([]res_secret.Secret, 0)
	for _, item := range list.Items {
		if !strings.Contains(item.Name, keyword) {
			continue
		}
		secretRes := secretConvert.K8sResConvert.SecretK8s2ResItemConvert(item)
		secretResList = append(secretResList, secretRes)
	}
	return secretResList, err
}
