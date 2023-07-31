package k8s

import (
	"github.com/gin-gonic/gin"
	ingress_req "kubmanager/model/ingress/request"
	"kubmanager/response"
)

type IngressApi struct {
}

func (i *IngressApi) CreateOrUpdateIngress(c *gin.Context) {
	ingressReq := ingress_req.Ingress{}
	if err := c.ShouldBindUri(&ingressReq); err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	err := ingressService.CreateOrUpdateIngress(ingressReq)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.Success(c)
}

func (i *IngressApi) DeleteIngress(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := ingressService.DeleteIngress(namespace, name)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.Success(c)
}

func (i *IngressApi) GetIngressListOrDetail(c *gin.Context) {
	namespace := c.Param("namespace")
	keyword := c.Query("keyword")
	name := c.Query("name")
	if name == "" {
		// 查询列表
		ingressListRes, err := ingressService.GetIngressList(namespace, keyword)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		response.SuccessWithDetailed(c, "查询ingress列表成功!", ingressListRes)
	} else {
		// 查询详情
		ingressRes, err := ingressService.GetIngressDetail(namespace, name)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		response.SuccessWithDetailed(c, "查询ingress成功!", ingressRes)
	}

}
