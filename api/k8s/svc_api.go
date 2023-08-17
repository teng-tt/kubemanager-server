package k8s

import (
	"github.com/gin-gonic/gin"
	svcReq "kubemanager.com/model/svc/request"
	"kubemanager.com/response"
)

type SVCApi struct {
}

func (s *SVCApi) CreateSOrUpdateVC(c *gin.Context) {
	var serviceReq svcReq.Service
	if err := c.ShouldBind(&serviceReq); err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	namespace := c.Param("namespace")
	err := svcService.CreateSOrUpdateVC(serviceReq, namespace)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.Success(c)
}

func (s *SVCApi) DeleteSVC(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := svcService.DeleteSVC(namespace, name)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.Success(c)
}

func (s *SVCApi) GetSVCListOrDetail(c *gin.Context) {
	namespace := c.Param("namespace")
	keyword := c.Query("keyword")
	name := c.Query("name")
	if name == "" {
		// 查询列表
		svcListRes, err := svcService.GetSVCList(namespace, keyword)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		response.SuccessWithDetailed(c, "查询service列表成功!", svcListRes)
	} else {
		// 查询详情
		serviceRes, err := svcService.GetSVCDetail(namespace, name)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		response.SuccessWithDetailed(c, "查询service成功!", serviceRes)
	}
}
