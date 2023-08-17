package k8s

import (
	"github.com/gin-gonic/gin"
	"kubemanager.com/api"
)

func initDeploymentRouter(group *gin.RouterGroup) {
	deploymentGroup := api.ApiGroupApp.K8SApiGroup.DeploymentApi
	group.POST("/deployment", deploymentGroup.CreateOrUpdateDeploy)
	group.GET("/deployment/:namespace", deploymentGroup.GetDeployDetailOrList)
	group.DELETE("/deployment/:namespace/:name", deploymentGroup.DeleteDeploy)
}
