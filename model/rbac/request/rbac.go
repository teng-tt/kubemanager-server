package request

import (
	rbacv1 "k8s.io/api/rbac/v1"
	"kubemanager.com/model/base"
)

type ServiceAccount struct {
	Name      string             `json:"name"`
	Namespace string             `json:"namespace"`
	Labels    []base.ListMapItem `json:"labels"`
}

type RoleReq struct {
	Name      string              `json:"name"`
	Namespace string              `json:"namespace"`
	Labels    []base.ListMapItem  `json:"labels"`
	Rules     []rbacv1.PolicyRule `json:"rules"`
}

type RoleBinding struct {
	Name      string             `json:"name"`
	Namespace string             `json:"namespace"`
	Labels    []base.ListMapItem `json:"labels"`
	Subjects  []ServiceAccount   `json:"subjects"` // 账号
	RoleRef   string             `json:"roleRef"`  // 角色
}
