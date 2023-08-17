package k8s

import (
	"github.com/gin-gonic/gin"
	"kubemanager.com/api"
)

func initCronjobRouter(group *gin.RouterGroup) {
	cronJobGroup := api.ApiGroupApp.K8SApiGroup.CronJobApi
	group.POST("/cronjob", cronJobGroup.CreateOrUpdateCronjob)
	group.GET("/cronjob/:namespace", cronJobGroup.GetCronjobDetailOrList)
	group.DELETE("/cronjob/:namespace/:name", cronJobGroup.DeleteCronjob)
}
