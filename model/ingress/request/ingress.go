package request

import (
	networkingv1 "k8s.io/api/networking/v1"
	"kubmanager/model/base"
)

type IngressRule struct {
	Host  string                        `json:"host"`
	Value networkingv1.IngressRuleValue `json:"value"`
}
type Ingress struct {
	Name      string             `json:"name"`
	Namespace string             `json:"namespace"`
	Labels    []base.ListMapItem `json:"labels"`
	Rules     []IngressRule      `json:"rules"`
}
