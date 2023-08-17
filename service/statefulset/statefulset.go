package statefulset

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubemanager.com/global"
	pvcReq "kubemanager.com/model/pvc/request"
	statefulsetReqs "kubemanager.com/model/statefulset/request"
	statefulsetResp "kubemanager.com/model/statefulset/response"
	"kubemanager.com/utils"
	"strconv"
	"strings"
)

type StatefulSetService struct {
}

func (s *StatefulSetService) CreateOrUpdateStatefulSet(statefulsetReq statefulsetReqs.StatefulSte) error {
	persistentVolumeClaim := make([]corev1.PersistentVolumeClaim, len(statefulsetReq.Base.VolumeClaimTemplates))
	for index, volumeClaimTemplate := range statefulsetReq.Base.VolumeClaimTemplates {
		persistentVolumeClaim[index] = corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:   volumeClaimTemplate.Name,
				Labels: utils.ToMap(volumeClaimTemplate.Labels),
			},
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes: volumeClaimTemplate.AccessModes,
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceStorage: resource.MustParse(strconv.Itoa(int(volumeClaimTemplate.Capacity)) + "Mi"),
					},
				},
				StorageClassName: &volumeClaimTemplate.StorageClassName,
			},
		}
	}
	podK8s := podConvert.Req2K8sConvert.PodReq2K8s(statefulsetReq.Template)
	statefulSet := appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      statefulsetReq.Base.Name,
			Namespace: statefulsetReq.Base.Namespace,
			Labels:    utils.ToMap(statefulsetReq.Base.Labels),
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas:    &statefulsetReq.Base.Replicas,
			ServiceName: statefulsetReq.Base.ServiceName,
			Selector: &metav1.LabelSelector{
				MatchLabels: utils.ToMap(statefulsetReq.Base.Selector),
			},
			VolumeClaimTemplates: persistentVolumeClaim,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: podK8s.ObjectMeta,
				Spec:       podK8s.Spec,
			},
		},
	}
	ctx := context.TODO()
	statefulApi := global.KubeConfigSet.AppsV1().StatefulSets(statefulSet.Namespace)
	statefulK8s, err := statefulApi.Get(ctx, statefulSet.Name, metav1.GetOptions{})
	if err != nil {
		// 创建
		_, err = statefulApi.Create(ctx, &statefulSet, metav1.CreateOptions{})
	} else {
		// 更新
		// 为防止服务抖动，id序号大的会先停止然后随即启动 -> id序号小的会先停止然后随即启动
		statefulK8s.Spec = statefulSet.Spec
		_, err = statefulApi.Update(ctx, statefulK8s, metav1.UpdateOptions{})
	}
	return err
}

func (s *StatefulSetService) DeleteStatefulSet(name, namespace string) error {
	err := global.KubeConfigSet.AppsV1().StatefulSets(namespace).
		Delete(context.TODO(), name, metav1.DeleteOptions{})
	return err
}

func (s *StatefulSetService) GetStatefulSetDetail(namespace, name string) (resp statefulsetReqs.StatefulSte, err error) {
	ctx := context.TODO()
	statefulApi := global.KubeConfigSet.AppsV1().StatefulSets(namespace)
	// 查询指定statefulSet
	statefulK8s, err := statefulApi.Get(ctx, name, metav1.GetOptions{})
	var statefulSetRespInfo statefulsetReqs.StatefulSte
	if err != nil {
		return statefulSetRespInfo, err
	}
	pvcReqList := make([]pvcReq.PersistentVolumeClaim, len(statefulK8s.Spec.VolumeClaimTemplates))
	for index, VolumeClaimTemplate := range statefulK8s.Spec.VolumeClaimTemplates {
		pvcReqList[index] = pvcReq.PersistentVolumeClaim{
			Name:             VolumeClaimTemplate.Name,
			AccessModes:      VolumeClaimTemplate.Spec.AccessModes,
			Capacity:         int32(VolumeClaimTemplate.Spec.Resources.Requests.Storage().Value() / (1024 * 1024)),
			StorageClassName: *VolumeClaimTemplate.Spec.StorageClassName,
		}
	}
	podReq := podConvert.K8s2RqeConver.PodK8s2Req(corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Labels: statefulK8s.Spec.Template.Labels,
		},
		Spec: statefulK8s.Spec.Template.Spec,
	})
	statefulSetRespInfo = statefulsetReqs.StatefulSte{
		Base: statefulsetReqs.StatefulSetBase{
			Name:                 statefulK8s.Name,
			Namespace:            statefulK8s.Namespace,
			Replicas:             *statefulK8s.Spec.Replicas,
			Labels:               utils.ToList(statefulK8s.Labels),
			Selector:             utils.ToList(statefulK8s.Spec.Selector.MatchLabels),
			ServiceName:          statefulK8s.Spec.ServiceName,
			VolumeClaimTemplates: pvcReqList,
		},
		Template: podReq,
	}
	return statefulSetRespInfo, err
}

func (s *StatefulSetService) GetStatefulSetList(namespace, keyword string) ([]statefulsetResp.StatefulSetResp, error) {
	ctx := context.TODO()
	statefulApi := global.KubeConfigSet.AppsV1().StatefulSets(namespace)
	// 查询列表
	list, err := statefulApi.List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	statefulRespList := make([]statefulsetResp.StatefulSetResp, 0)
	for _, item := range list.Items {
		if !strings.Contains(item.Name, keyword) {
			continue
		}
		statefulRespList = append(statefulRespList, statefulsetResp.StatefulSetResp{
			Name:      item.Name,
			Namespace: item.Namespace,
			Ready:     item.Status.ReadyReplicas,
			Replicas:  item.Status.Replicas,
			Age:       item.CreationTimestamp.UnixMilli(),
		})
	}
	return statefulRespList, err
}
