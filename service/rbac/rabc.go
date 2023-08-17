package rbac

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubemanager.com/global"
	rbacReqs "kubemanager.com/model/rbac/request"
	rbacRes "kubemanager.com/model/rbac/response"
	"kubemanager.com/utils"
	"strings"
)

type RbacServiceApi struct {
}

func (r *RbacServiceApi) CreateServiceAccount(saReq rbacReqs.ServiceAccount) error {
	saK8s := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      saReq.Name,
			Namespace: saReq.Namespace,
			Labels:    utils.ToMap(saReq.Labels),
		},
	}
	_, err := global.KubeConfigSet.CoreV1().ServiceAccounts(saK8s.Namespace).
		Create(context.TODO(), saK8s, metav1.CreateOptions{})
	return err
}

func (r *RbacServiceApi) DeleteServiceAccount(name, namespace string) error {
	err := global.KubeConfigSet.CoreV1().ServiceAccounts(namespace).
		Delete(context.TODO(), name, metav1.DeleteOptions{})
	return err
}

func (r *RbacServiceApi) GetServiceAccountList(namespace, keyword string) ([]rbacRes.ServiceAccountRes, error) {
	list, err := global.KubeConfigSet.CoreV1().ServiceAccounts(namespace).
		List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	saResList := make([]rbacRes.ServiceAccountRes, 0)
	for _, item := range list.Items {
		if !strings.Contains(item.Name, keyword) {
			continue
		}
		saResList = append(saResList, rbacRes.ServiceAccountRes{
			Name:      item.Name,
			Namespace: item.Namespace,
			Age:       item.CreationTimestamp.UnixMilli(),
		})

	}
	return saResList, err
}

func (r *RbacServiceApi) CreateOrUpdateRole(roleReq rbacReqs.RoleReq) (err error) {
	ctx := context.TODO()
	// 创建ClusterRole
	if roleReq.Namespace == "" {
		clusterRoleApi := global.KubeConfigSet.RbacV1().ClusterRoles()
		clusterRoleK8s := &rbacv1.ClusterRole{
			ObjectMeta: metav1.ObjectMeta{
				Name:      roleReq.Name,
				Namespace: roleReq.Namespace,
				Labels:    utils.ToMap(roleReq.Labels),
			},
			Rules: roleReq.Rules,
		}
		clusterRoleSrc, err := clusterRoleApi.Get(ctx, roleReq.Name, metav1.GetOptions{})
		if err != nil {
			// 创建
			_, err = clusterRoleApi.Create(ctx, clusterRoleK8s, metav1.CreateOptions{})
			return err
		} else {
			// 更新
			clusterRoleSrc.ObjectMeta.Labels = clusterRoleK8s.Labels
			clusterRoleSrc.Rules = clusterRoleK8s.Rules
			_, err = clusterRoleApi.Update(ctx, clusterRoleSrc, metav1.UpdateOptions{})
			return err
		}
	} else {
		// 创建 namespace role
		nsRoleApi := global.KubeConfigSet.RbacV1().Roles(roleReq.Namespace)
		nsRoleK8s := &rbacv1.Role{
			ObjectMeta: metav1.ObjectMeta{
				Name:      roleReq.Name,
				Namespace: roleReq.Namespace,
				Labels:    utils.ToMap(roleReq.Labels),
			},
			Rules: roleReq.Rules,
		}
		roleSrc, err := nsRoleApi.Get(ctx, roleReq.Name, metav1.GetOptions{})
		if err != nil {
			// 创建
			_, err = nsRoleApi.Create(ctx, nsRoleK8s, metav1.CreateOptions{})
			return err
		} else {
			// 更新
			roleSrc.ObjectMeta.Labels = nsRoleK8s.Labels
			roleSrc.Rules = nsRoleK8s.Rules
			_, err = nsRoleApi.Update(ctx, roleSrc, metav1.UpdateOptions{})
			return err
		}
	}
}

func (r *RbacServiceApi) GetRoleDetail(namespace, name string) (roleResInfo rbacReqs.RoleReq, err error) {
	ctx := context.TODO()
	if namespace != "" {
		// role查看详情
		roleApi := global.KubeConfigSet.RbacV1().Roles(namespace)
		roleK8s, err := roleApi.Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return roleResInfo, err
		}
		roleResInfo = rbacReqs.RoleReq{
			Name:      roleK8s.Name,
			Namespace: roleK8s.Namespace,
			Labels:    utils.ToList(roleK8s.Labels),
			Rules:     roleK8s.Rules,
		}
	} else {
		// clusterRole查看详情
		clusterRoleApi := global.KubeConfigSet.RbacV1().ClusterRoles()
		roleK8s, err := clusterRoleApi.Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return roleResInfo, err
		}
		roleResInfo = rbacReqs.RoleReq{
			Name:      roleK8s.Name,
			Namespace: roleK8s.Namespace,
			Labels:    utils.ToList(roleK8s.Labels),
			Rules:     roleK8s.Rules,
		}
	}
	return roleResInfo, err
}

func (r *RbacServiceApi) GetRoleList(namespace, keyword string) ([]rbacRes.RoleRes, error) {
	ctx := context.TODO()
	resRoleList := make([]rbacRes.RoleRes, 0)
	if namespace != "" {
		roleApi := global.KubeConfigSet.RbacV1().Roles(namespace)
		roleList, err := roleApi.List(ctx, metav1.ListOptions{})
		if err != nil {
			return resRoleList, err
		}
		for _, item := range roleList.Items {
			if !strings.Contains(item.Name, keyword) {
				continue
			}
			resRoleList = append(resRoleList, rbacRes.RoleRes{
				Name:      item.Name,
				Namespace: item.Namespace,
				Age:       item.CreationTimestamp.UnixMilli(),
			})
		}
	} else {
		clusterRoleApi := global.KubeConfigSet.RbacV1().ClusterRoles()
		roleList, err := clusterRoleApi.List(ctx, metav1.ListOptions{})
		if err != nil {
			return resRoleList, err
		}
		for _, item := range roleList.Items {
			if !strings.Contains(item.Name, keyword) {
				continue
			}
			resRoleList = append(resRoleList, rbacRes.RoleRes{
				Name:      item.Name,
				Namespace: item.Namespace,
				Age:       item.CreationTimestamp.UnixMilli(),
			})
		}
	}
	return resRoleList, nil
}

func (r *RbacServiceApi) DeleteRole(namespace, name string) (err error) {
	ctx := context.TODO()
	if namespace != "" {
		err = global.KubeConfigSet.RbacV1().Roles(namespace).
			Delete(ctx, name, metav1.DeleteOptions{})
	} else {
		err = global.KubeConfigSet.RbacV1().ClusterRoles().
			Delete(ctx, name, metav1.DeleteOptions{})
	}
	return err
}

func (r *RbacServiceApi) CreateOrUpdateRb(rbReq rbacReqs.RoleBinding) error {
	ctx := context.TODO()
	//创建 cluster role binding
	if rbReq.Namespace == "" {
		rbK8sReq := &rbacv1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      rbReq.Name,
				Namespace: rbReq.Namespace,
				Labels:    utils.ToMap(rbReq.Labels),
			},
			Subjects: func(saList []rbacReqs.ServiceAccount) []rbacv1.Subject {
				subjects := make([]rbacv1.Subject, len(saList))
				for index, item := range saList {
					subjects[index] = rbacv1.Subject{
						Name:      item.Name,
						Kind:      "User",
						Namespace: item.Namespace,
					}
				}
				return subjects
			}(rbReq.Subjects),
			RoleRef: rbacv1.RoleRef{
				Name:     rbReq.RoleRef,
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "ClusterRole",
			},
		}
		clusterRbApi := global.KubeConfigSet.RbacV1().ClusterRoleBindings()
		if cluterRoleSrc, err := clusterRbApi.Get(ctx, rbReq.Name, metav1.GetOptions{}); err != nil {
			_, err := clusterRbApi.Create(ctx, rbK8sReq, metav1.CreateOptions{})
			if err != nil {
				return err
			}
		} else {
			cluterRoleSrc.ObjectMeta.Labels = rbK8sReq.Labels
			cluterRoleSrc.Subjects = rbK8sReq.Subjects
			cluterRoleSrc.RoleRef = rbK8sReq.RoleRef
			_, err := clusterRbApi.Update(ctx, cluterRoleSrc, metav1.UpdateOptions{})
			if err != nil {
				return err
			}
		}
	} else {
		rbK8sReq := &rbacv1.RoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      rbReq.Name,
				Namespace: rbReq.Namespace,
				Labels:    utils.ToMap(rbReq.Labels),
			},
			Subjects: func(saList []rbacReqs.ServiceAccount) []rbacv1.Subject {
				subjects := make([]rbacv1.Subject, len(saList))
				for index, item := range saList {
					subjects[index] = rbacv1.Subject{
						Name:      item.Name,
						Kind:      "User",
						Namespace: item.Namespace,
					}
				}
				return subjects
			}(rbReq.Subjects),
			RoleRef: rbacv1.RoleRef{
				Name:     rbReq.RoleRef,
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "Role",
			},
		}
		rbApi := global.KubeConfigSet.RbacV1().RoleBindings(rbK8sReq.Namespace)
		if rbK8sSrc, err := rbApi.Get(ctx, rbReq.Name, metav1.GetOptions{}); err != nil {
			_, err = rbApi.Create(ctx, rbK8sReq, metav1.CreateOptions{})
			if err != nil {
				return err
			}
		} else {
			rbK8sSrc.ObjectMeta.Labels = rbK8sReq.Labels
			rbK8sSrc.Subjects = rbK8sReq.Subjects
			rbK8sSrc.RoleRef = rbK8sReq.RoleRef
			_, err = rbApi.Update(ctx, rbK8sSrc, metav1.UpdateOptions{})
			if err != nil {
				return err
			}
		}
	}
	return nil

}

func (r *RbacServiceApi) DeleteRb(namespace, name string) error {
	ctx := context.TODO()
	if namespace != "" {
		err := global.KubeConfigSet.RbacV1().RoleBindings(namespace).
			Delete(ctx, name, metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	} else {
		err := global.KubeConfigSet.RbacV1().ClusterRoleBindings().
			Delete(ctx, name, metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

// GetRbDetail 查看RoleBinding详情
func (r *RbacServiceApi) GetRbDetail(namespace, name string) (rbacReqs.RoleBinding, error) {
	ctx := context.TODO()
	rbReq := rbacReqs.RoleBinding{}
	if namespace != "" {
		rbK8s, err := global.KubeConfigSet.RbacV1().RoleBindings(namespace).
			Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return rbReq, err
		}
		rbReq.Name = rbK8s.Name
		rbReq.Namespace = rbK8s.Namespace
		rbReq.Labels = utils.ToList(rbK8s.Labels)
		rbReq.RoleRef = rbK8s.RoleRef.Name
		rbReq.Subjects = func(subjects []rbacv1.Subject) []rbacReqs.ServiceAccount {
			saList := make([]rbacReqs.ServiceAccount, len(subjects))
			for i, subject := range subjects {
				saList[i] = rbacReqs.ServiceAccount{
					Name:      subject.Name,
					Namespace: subject.Namespace,
				}
			}
			return saList
		}(rbK8s.Subjects)
	} else {
		rbK8s, err := global.KubeConfigSet.RbacV1().ClusterRoleBindings().
			Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return rbReq, err
		}
		rbReq.Name = rbK8s.Name
		rbReq.Namespace = rbK8s.Namespace
		rbReq.Labels = utils.ToList(rbK8s.Labels)
		rbReq.RoleRef = rbK8s.RoleRef.Name
		rbReq.Subjects = func(subjects []rbacv1.Subject) []rbacReqs.ServiceAccount {
			saList := make([]rbacReqs.ServiceAccount, len(subjects))
			for i, subject := range subjects {
				saList[i] = rbacReqs.ServiceAccount{
					Name:      subject.Name,
					Namespace: subject.Namespace,
				}
			}
			return saList
		}(rbK8s.Subjects)
	}
	return rbReq, nil
}

// GetRbList 查看RoleBinding列表
func (r *RbacServiceApi) GetRbList(namespace, keyword string) ([]rbacRes.RoleBindingRes, error) {
	ctx := context.TODO()
	rbResList := make([]rbacRes.RoleBindingRes, 0)
	if namespace != "" {
		list, err := global.KubeConfigSet.RbacV1().RoleBindings(namespace).
			List(ctx, metav1.ListOptions{})
		if err != nil {
			return rbResList, err
		}
		for _, item := range list.Items {
			if !strings.Contains(item.Name, keyword) {
				continue
			}
			rbResList = append(rbResList, rbacRes.RoleBindingRes{
				Name:      item.Name,
				Namespace: item.Namespace,
				Age:       item.CreationTimestamp.Unix(),
			})
		}
	} else {
		list, err := global.KubeConfigSet.RbacV1().ClusterRoleBindings().
			List(ctx, metav1.ListOptions{})
		if err != nil {
			return rbResList, err
		}
		for _, item := range list.Items {
			if !strings.Contains(item.Name, keyword) {
				continue
			}
			rbResList = append(rbResList, rbacRes.RoleBindingRes{
				Name:      item.Name,
				Namespace: item.Namespace,
				Age:       item.CreationTimestamp.Unix(),
			})
		}
	}
	return rbResList, nil
}
