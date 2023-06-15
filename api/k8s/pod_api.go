package k8s

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubmanager/global"
	pod_req "kubmanager/model/pod/request"
	"kubmanager/response"
	"net/http"
)

type PodApi struct {
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
	k8sPod := podConvert.PodReq2K8s(podReq)
	ctx := context.TODO()
	createdPod, err := global.KubeConfigSet.CoreV1().Pods(k8sPod.Namespace).Create(ctx, k8sPod, metav1.CreateOptions{})
	if err != nil {
		errMsg := fmt.Sprintf("Poc[%s-%s]创建失败，detail：%s", k8sPod.Namespace, k8sPod.Name, err.Error())
		response.FailWithMessage(c, errMsg)
		return
	}
	successMsg := fmt.Sprintf("Pod[%s-%s]创建成功", createdPod.Namespace, createdPod.Name)
	response.SuccessWithMessage(c, successMsg)

}

func (p *PodApi) GetPodList(c *gin.Context) {
	ctx := context.TODO()
	list, err := global.KubeConfigSet.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Println(err.Error())
	}

	for _, item := range list.Items {
		fmt.Println(item.Namespace, item.Name)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
