package k8s

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"kubmanager/global"
	pod_req "kubmanager/model/pod/request"
	"kubmanager/response"
)

type PodApi struct {
}

// UpdatePod 因为update的字段属性有限，而实际更新过程当中 会修改定义的任意字段
func (p *PodApi) UpdatePod(ctx context.Context, pod *corev1.Pod) error {
	_, err := global.KubeConfigSet.CoreV1().Pods(pod.Namespace).Update(ctx, pod, metav1.UpdateOptions{})
	return err
}

// PatchPod 打补丁，需要传指定更新的地方进行操作，比较麻烦， 打补丁：没有就添加，有就执行更新
func (p *PodApi) PatchPod(ctx context.Context, pod *corev1.Pod, patchData map[string]interface{}) error {
	//patchData["metadata"] = map[string]interface{}{
	//	"labels": map[string]string{
	//		"foo": "bar2",
	//		"app": "testnginx",
	//	},
	//}

	patchDataByte, _ := json.Marshal(patchData)
	_, err := global.KubeConfigSet.CoreV1().Pods(pod.Namespace).Patch(
		ctx,
		pod.Name,
		types.StrategicMergePatchType,
		patchDataByte,
		metav1.PatchOptions{},
	)

	return err
}

func (p *PodApi) CreateOrUpdatePod(c *gin.Context) {
	var podReq pod_req.Pod
	if err := c.ShouldBind(&podReq); err != nil {
		response.FailWithMessage(c, "参数解析失败， detail: "+err.Error())
		return
	}
	// 校验必填项
	if err := podValidate.Validate(&podReq); err != nil {
		response.FailWithMessage(c, "参数验证失败， detail: "+err.Error())
		return
	}
	if msg, err := podService.CreateOrUpdate(podReq); err != nil {
		response.FailWithMessage(c, msg)
	} else {
		response.SuccessWithMessage(c, msg)
	}

}

func (p *PodApi) GetPodListOrDetail(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Query("name")
	keyword := c.Query("keyword")
	if name != "" {
		detail, err := podService.GetPodDetail(namespace, name)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		response.SuccessWithDetailed(c, "获取Pod详情成功", detail)
	} else {
		podList, err := podService.GetPodList(namespace, keyword, c.Query("nodeName"))
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		response.SuccessWithDetailed(c, "获取Pod列表成功", podList)
	}

}

func (*PodApi) DeletePod(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := podService.DeletePod(namespace, name)
	if err != nil {
		response.FailWithMessage(c, "删除Pod失败，detail："+err.Error())
	} else {
		response.Success(c)
	}
}
