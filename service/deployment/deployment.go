package deployment

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubmanager/global"
	deployReq "kubmanager/model/deployment/request"
	deployResp "kubmanager/model/deployment/response"
	"kubmanager/utils"
	"strings"
)

type DeploymentService struct {
}

func (d *DeploymentService) CreateOrUpdateDeploy(deploymentReq deployReq.Deployment) error {
	// 转换为 K8S 结构
	podK8s := podConvert.Req2K8sConvert.PodReq2K8s(deploymentReq.Template)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentReq.Base.Name,
			Namespace: deploymentReq.Base.Namespace,
			Labels:    utils.ToMap(deploymentReq.Base.Labels),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &deploymentReq.Base.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: utils.ToMap(deploymentReq.Base.Selector),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: podK8s.ObjectMeta,
				Spec:       podK8s.Spec,
			},
		},
	}
	ctx := context.TODO()
	deploymentApi := global.KubeConfigSet.AppsV1().Deployments(deployment.Namespace)
	deploymentK8s, err := deploymentApi.Get(ctx, deployment.Name, metav1.GetOptions{})
	if err == nil {
		// 更新
		deploymentK8s.Spec = deployment.Spec
		_, err = deploymentApi.Update(ctx, deploymentK8s, metav1.UpdateOptions{})
	} else {
		// 创建
		_, err = deploymentApi.Create(ctx, deployment, metav1.CreateOptions{})
	}
	return err
}

func (d *DeploymentService) DeleteDeploy(name, namespace string) error {
	return global.KubeConfigSet.AppsV1().Deployments(namespace).
		Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func (d *DeploymentService) GetDeployList(namespace, keyword string) (deploymentList []deployResp.Deployment, err error) {
	deploymentResList := make([]deployResp.Deployment, 0)
	list, err := global.KubeConfigSet.AppsV1().Deployments(namespace).
		List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return deploymentResList, err
	}
	for _, item := range list.Items {
		if !strings.Contains(item.Name, keyword) {
			continue
		}
		deploymentResList = append(deploymentResList, deployResp.Deployment{
			Name:      item.Name,
			Namespace: item.Namespace,
			Age:       item.CreationTimestamp.Unix(),
			Replicas:  *item.Spec.Replicas,
			Ready:     item.Status.Replicas,
			Available: item.Status.AvailableReplicas,
			UpToDate:  item.Status.UpdatedReplicas,
		})
	}
	return deploymentResList, err
}

func (d *DeploymentService) GetDeployDetail(name, namespace string) (resp deployReq.Deployment, err error) {
	var deloymentReq deployReq.Deployment
	deploymentK8s, err := global.KubeConfigSet.AppsV1().Deployments(namespace).
		Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return deloymentReq, err
	}
	podReq := podConvert.K8s2RqeConver.PodK8s2Req(corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Labels: deploymentK8s.Spec.Template.Labels,
		},
		Spec: deploymentK8s.Spec.Template.Spec,
	})
	deloymentReq = deployReq.Deployment{
		Base: deployReq.DeploymentBase{
			Name:      deploymentK8s.Name,
			Namespace: deploymentK8s.Namespace,
			Replicas:  *deploymentK8s.Spec.Replicas,
			Labels:    utils.ToList(deploymentK8s.Labels),
			Selector:  utils.ToList(deploymentK8s.Spec.Selector.MatchLabels),
		},
		Template: podReq,
	}
	return deloymentReq, err
}
