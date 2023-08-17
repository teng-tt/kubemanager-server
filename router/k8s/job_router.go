package k8s

import (
	"github.com/gin-gonic/gin"
	"kubemanager.com/api"
)

func initJobRouter(group *gin.RouterGroup) {
	jobGroup := api.ApiGroupApp.K8SApiGroup.JoApi
	group.POST("/job", jobGroup.CreateOrUpdateJob)
	group.GET("/job/:namespace", jobGroup.GetJobDetailOrList)
	group.DELETE("/job/:namespace/:name", jobGroup.DeleteJob)
}
