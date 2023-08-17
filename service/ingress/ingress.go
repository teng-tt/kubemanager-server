package ingress

import (
	"context"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubemanager.com/global"
	ingress_req "kubemanager.com/model/ingress/request"
	ingress_resp "kubemanager.com/model/ingress/response"
	"kubemanager.com/utils"
	"strings"
)

type IngresService struct {
}

func (i *IngresService) CreateOrUpdateIngress(ingressReq ingress_req.Ingress) error {
	ingressRules := make([]networkingv1.IngressRule, 0)
	for _, rule := range ingressReq.Rules {
		ingressRules = append(ingressRules, networkingv1.IngressRule{
			Host:             rule.Host,
			IngressRuleValue: rule.Value,
		})
	}
	ingress := networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ingressReq.Name,
			Namespace: ingressReq.Namespace,
		},
		Spec: networkingv1.IngressSpec{
			Rules: ingressRules,
		},
	}
	ingressApi := global.KubeConfigSet.NetworkingV1().Ingresses(ingress.Namespace)
	ctx := context.TODO()
	// 查询是否存在
	ingressK8s, err := ingressApi.Get(ctx, ingress.Name, metav1.GetOptions{})
	if err == nil {
		// 更新
		ingressK8s.Spec = ingress.Spec
		_, err = ingressApi.Update(ctx, ingressK8s, metav1.UpdateOptions{})
	} else {
		// 创建
		_, err = ingressApi.Create(ctx, &ingress, metav1.CreateOptions{})
	}
	return err
}

func (i *IngresService) DeleteIngress(namespace, name string) error {
	err := global.KubeConfigSet.NetworkingV1().Ingresses(namespace).
		Delete(context.TODO(), name, metav1.DeleteOptions{})
	return err
}

func (i *IngresService) GetIngressList(namespace, keyword string) ([]ingress_resp.IngressResp, error) {
	ingressList, err := global.KubeConfigSet.NetworkingV1().Ingresses(namespace).
		List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	ingressListRes := make([]ingress_resp.IngressResp, 0)
	for _, item := range ingressList.Items {
		if !strings.Contains(item.Name, keyword) {
			continue
		}
		hosts := make([]string, 0)
		for _, rule := range item.Spec.Rules {
			hosts = append(hosts, rule.Host)
		}
		ingressListRes = append(ingressListRes, ingress_resp.IngressResp{
			Name:      item.Name,
			Namespace: item.Namespace,
			Class:     *item.Spec.IngressClassName,
			Hosts:     strings.Join(hosts, ","),
			Age:       item.CreationTimestamp.UnixMilli(),
		})

	}
	return ingressListRes, err
}
func (i *IngresService) GetIngressDetail(namespace, name string) (ingresReq ingress_req.Ingress, err error) {
	ingressK8s, err := global.KubeConfigSet.NetworkingV1().Ingresses(namespace).
		Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return
	}
	rules := make([]ingress_req.IngressRule, 0)
	for _, rule := range ingressK8s.Spec.Rules {
		rules = append(rules, ingress_req.IngressRule{
			Host:  rule.Host,
			Value: rule.IngressRuleValue,
		})
	}
	ingressReq := ingress_req.Ingress{
		Name:      ingressK8s.Name,
		Namespace: ingressK8s.Namespace,
		Labels:    utils.ToList(ingressK8s.Labels),
		Rules:     rules,
	}
	return ingressReq, err
}
