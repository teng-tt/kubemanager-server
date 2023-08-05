package k8s

import (
	"github.com/gin-gonic/gin"
	deployReq "kubmanager/model/deployment/request"
	"kubmanager/response"
)

type DeploymentApi struct {
}

func (d *DeploymentApi) CreateOrUpdateDeploy(c *gin.Context) {
	var deploymentReq deployReq.Deployment
	if err := c.ShouldBind(&deploymentReq); err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	err := deploymentService.CreateOrUpdateDeploy(deploymentReq)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.Success(c)
}

func (d *DeploymentApi) DeleteDeploy(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	err := deploymentService.DeleteDeploy(name, namespace)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.Success(c)
}

func (d *DeploymentApi) GetDeployDetailOrList(c *gin.Context) {
	namespace := c.Param("namespace")
	keyword := c.Query("keyword")
	name := c.Query("name")
	if name == "" || len(name) == 0 {
		// 查询列表
		deploymentRespList, err := deploymentService.GetDeployList(namespace, keyword)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		response.SuccessWithDetailed(c, "查询deployment列表成功!", deploymentRespList)
	} else {
		// 查询指定deployment
		deploymentRespInfo, err := deploymentService.GetDeployDetail(name, namespace)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		response.SuccessWithDetailed(c, "deployment!", deploymentRespInfo)
	}
}
