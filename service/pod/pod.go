package pod

import (
	"context"
	"errors"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"kubmanager/global"
	pod_req "kubmanager/model/pod/request"
	pod_res "kubmanager/model/pod/response"
	"strings"
)

type PodService struct {
}

func (p *PodService) CreateOrUpdate(podReq pod_req.Pod) (msg string, err error) {
	ctx := context.TODO()
	k8sPod := podConvert.Req2K8sConvert.PodReq2K8s(podReq)
	podApi := global.KubeConfigSet.CoreV1().Pods(k8sPod.Namespace)

	// 更新
	if getK8sPod, err := podApi.Get(ctx, k8sPod.Name, metav1.GetOptions{}); err == nil {
		// 更新pod三种更新方式：[no] update | [no] patch  | [yes] delete+create
		// 校验pod 参数是否合理，避免不合法参数创建失败
		k8sPodCopy := *k8sPod
		k8sPodCopy.Name = k8sPod.Name + "-validate"
		_, err := podApi.Create(ctx, &k8sPodCopy, metav1.CreateOptions{
			DryRun: []string{metav1.DryRunAll},
		})
		if err != nil {
			errMsg := fmt.Sprintf("Poc[namespace=%s, name=%s]更新失败，detail：%s", k8sPod.Namespace, k8sPod.Name, err.Error())
			return errMsg, err
		}
		// 删除 --强制删除，减少删除等待时间，防止前端等待删除超时
		background := metav1.DeletePropagationBackground
		var gracePeriodSeconds int64 = 0
		err = podApi.Delete(ctx, k8sPod.Name, metav1.DeleteOptions{
			GracePeriodSeconds: &gracePeriodSeconds,
			PropagationPolicy:  &background,
		})
		if err != nil {
			errMsg := fmt.Sprintf("Poc[namespace=%s, name=%s]更新失败，detail：%s", k8sPod.Namespace, k8sPod.Name, err.Error())
			return errMsg, err
		}
		// 当pod处于terminating状态，监听pod删除完毕之后 才开始创建pod
		var labelSelector []string
		for k, v := range getK8sPod.Labels {
			labelSelector = append(labelSelector, fmt.Sprintf("%s=%s", k, v))
		}
		// label 格式 app=test,app=test2
		watcher, err := podApi.Watch(ctx, metav1.ListOptions{
			LabelSelector: strings.Join(labelSelector, ","),
		})
		if err != nil {
			errMsg := fmt.Sprintf("Poc[namespace=%s, name=%s]更新失败，detail：%s", k8sPod.Namespace, k8sPod.Name, err.Error())
			return errMsg, err
		}
		for event := range watcher.ResultChan() {
			// 获取通道事件，事件可能相同，其它通道也有一样的事件
			k8sPodChan := event.Object.(*corev1.Pod)
			// 查询k8s,判断是否已经删除，那么就不用判断删除事件了,f防止进入死循环
			if _, err := podApi.Get(ctx, k8sPod.Name, metav1.GetOptions{}); k8serror.IsNotFound(err) {
				// 重新创建
				if createdPod, err := podApi.Create(ctx, k8sPod, metav1.CreateOptions{}); err != nil {
					errMsg := fmt.Sprintf("Poc[namespace=%s, name=%s]更新失败，detail：%s", k8sPod.Namespace, k8sPod.Name, err.Error())
					return errMsg, err
				} else {
					successMsg := fmt.Sprintf("Pod[namespace=%s, name=%s]更新成功", createdPod.Namespace, createdPod.Name)
					return successMsg, err
				}
			}
			switch event.Type {
			case watch.Deleted:
				// 防止相同的事件，只有在当前通道的事件才创建
				if k8sPodChan.Name != k8sPod.Name {
					continue
				}
				// 重新创建
				if createdPod, err := podApi.Create(ctx, k8sPod, metav1.CreateOptions{}); err != nil {
					errMsg := fmt.Sprintf("Poc[namespace=%s, name=%s]更新失败，detail：%s", k8sPod.Namespace, k8sPod.Name, err.Error())
					return errMsg, err
				} else {
					successMsg := fmt.Sprintf("Pod[namespace=%s, name=%s]更新成功", createdPod.Namespace, createdPod.Name)
					return successMsg, err
				}
			}
		}
		return "", nil
	} else {
		// 创建Pod
		if createdPod, err := podApi.Create(ctx, k8sPod, metav1.CreateOptions{}); err != nil {
			errMsg := fmt.Sprintf("Poc[namespace=%s, name=%s]创建失败，detail：%s", k8sPod.Namespace, k8sPod.Name, err.Error())
			return errMsg, err
		} else {
			successMsg := fmt.Sprintf("Pod[namespace=%s, name=%s]创建成功", createdPod.Namespace, createdPod.Name)
			return successMsg, err
		}
	}
}

func (p *PodService) GetPodDetail(namespace, name string) (podReq pod_req.Pod, err error) {
	ctx := context.TODO()
	podApi := global.KubeConfigSet.CoreV1().Pods(namespace)
	k8sGetPod, err := podApi.Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		errMsg := fmt.Sprintf("Poc[namespace=%s, name=%s]查询失败，detail：%s", namespace, name, err.Error())
		err = errors.New(errMsg)
		return
	}
	// 将k8s pod 转为 pod request
	podReq = podConvert.K8s2RqeConver.PodK8s2Req(*k8sGetPod)
	return
}

func (p *PodService) GetPodList(namespace, keyword string) (podList []pod_res.PodListItem, err error) {
	ctx := context.TODO()
	list, err := global.KubeConfigSet.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return

	}
	podList = make([]pod_res.PodListItem, 0)
	for _, item := range list.Items {
		if strings.Contains(item.Name, keyword) {
			podItem := podConvert.K8s2RqeConver.PodK8s2ItemRes(item)
			podList = append(podList, podItem)
		}
	}

	return podList, err
}

func (p *PodService) DeletePod(namespace, name string) error {
	// 删除 --强制删除，减少删除等待时间，防止前端等待删除超时
	background := metav1.DeletePropagationBackground // 后台删除
	var gracePeriodSeconds int64 = 0                 // 等待时间
	return global.KubeConfigSet.CoreV1().Pods(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{
		GracePeriodSeconds: &gracePeriodSeconds,
		PropagationPolicy:  &background,
	})
}
