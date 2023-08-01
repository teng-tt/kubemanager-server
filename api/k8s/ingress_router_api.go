package k8s

import (
	"github.com/gin-gonic/gin"
	ingrouteReq "kubmanager/model/ingroute/request"
	"kubmanager/response"
)

type IngressRouteApi struct {
}

func (i *IngressRouteApi) CreateOrUpdateIngRoute(c *gin.Context) {
	var ingressRouteReq ingrouteReq.IngressRouteReq
	if err := c.ShouldBind(&ingressRouteReq); err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	err := ingrouteService.CreateOrUpdateIngRoute(ingressRouteReq)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.Success(c)
}

func (i *IngressRouteApi) DeleteIngRoute(c *gin.Context) {
	err := ingrouteService.DeleteIngRoute(c.Param("namespace"), c.Param("name"))
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.Success(c)

}

func (i *IngressRouteApi) GetIngRouteDetailOrList(c *gin.Context) {
	namespace := c.Param("namespace")
	keyword := c.Query("keyword")
	name := c.Query("name")
	if name == "" {
		ingRtResp, err := ingrouteService.GetIngRouteList(namespace, keyword)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		response.SuccessWithDetailed(c, "获取IngressRoute列表成功!", ingRtResp)
	} else {
		ingRtReq, err := ingrouteService.GetIngRouteDetail(namespace, name)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		response.SuccessWithDetailed(c, "获取Ingress详情成功!", ingRtReq)
	}
}

func (i *IngressRouteApi) GetIngRouteMiddlewareList(c *gin.Context) {
	list, err := ingrouteService.GetIngRouteMiddlewareList(c.Param("namespace"))
	if err != nil {
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
	}
	response.SuccessWithDetailed(c, "查询成功!", list)
}
