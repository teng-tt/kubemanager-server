package svc

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"kubemanager.com/global"
	svcReq "kubemanager.com/model/svc/request"
	svcRes "kubemanager.com/model/svc/response"
	"kubemanager.com/utils"
	"strings"
)

type Service struct {
}

func (s *Service) CreateSOrUpdateVC(serviceReq svcReq.Service, namespace string) error {
	servicePorts := make([]corev1.ServicePort, 0)
	for _, item := range serviceReq.Ports {
		servicePorts = append(servicePorts, corev1.ServicePort{
			Name: item.Name,
			Port: item.Port,
			TargetPort: intstr.IntOrString{
				Type:   intstr.Int,
				IntVal: item.TargetPort,
			},
			NodePort: item.NodePort,
		})
	}
	service := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceReq.Name,
			Namespace: serviceReq.Namespace,
			Labels:    utils.ToMap(serviceReq.Labels),
		},
		Spec: corev1.ServiceSpec{
			Type:     serviceReq.Type,
			Selector: utils.ToMap(serviceReq.Selector),
			Ports:    servicePorts,
		},
	}
	// 提交K8S
	ctx := context.TODO()
	serviceApi := global.KubeConfigSet.CoreV1().Services(namespace)
	serviceK8s, err := serviceApi.Get(ctx, service.Name, metav1.GetOptions{})
	if err == nil {
		// 存在更新
		serviceK8s.Spec = service.Spec
		_, err = serviceApi.Update(ctx, serviceK8s, metav1.UpdateOptions{})
	} else {
		// 创建
		_, err = serviceApi.Create(ctx, &service, metav1.CreateOptions{})
	}
	return err
}

func (s *Service) DeleteSVC(namespace, name string) error {
	err := global.KubeConfigSet.CoreV1().Services(namespace).
		Delete(context.TODO(), name, metav1.DeleteOptions{})
	return err
}

func (s *Service) GetSVCList(namespace, keyword string) ([]svcRes.ServiceRes, error) {
	// 查询列表
	svcList, err := global.KubeConfigSet.CoreV1().Services(namespace).
		List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	serviceResList := make([]svcRes.ServiceRes, 0)
	for _, item := range svcList.Items {
		if !strings.Contains(item.Name, keyword) {
			continue
		}
		serviceResList = append(serviceResList, svcRes.ServiceRes{
			Name:       item.Name,
			Namespace:  item.Namespace,
			Type:       item.Spec.Type,
			ClusterIP:  item.Spec.ClusterIP,
			ExternalIP: item.Spec.ExternalIPs,
			Age:        item.CreationTimestamp.UnixMilli(),
		})
	}
	return serviceResList, err
}

func (s *Service) GetSVCDetail(namespace, name string) (svcReq.Service, error) {
	// 查询详情
	var serviceReq svcReq.Service
	serviceK8s, err := global.KubeConfigSet.CoreV1().Services(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return serviceReq, err
	}
	servicePorts := make([]svcReq.ServicePort, 0)
	for _, item := range serviceK8s.Spec.Ports {
		servicePorts = append(servicePorts, svcReq.ServicePort{
			Name:       item.Name,
			Port:       item.Port,
			TargetPort: item.TargetPort.IntVal,
			NodePort:   item.NodePort,
		})
	}
	serviceReq = svcReq.Service{
		Name:      serviceK8s.Name,
		Namespace: serviceK8s.Namespace,
		Labels:    utils.ToList(serviceK8s.Labels),
		Type:      serviceK8s.Spec.Type,
		Selector:  utils.ToList(serviceK8s.Spec.Selector),
		Ports:     servicePorts,
	}
	return serviceReq, err
}
