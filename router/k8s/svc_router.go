package k8s

import (
	"github.com/gin-gonic/gin"
	"kubmanager/api"
)

func initSvcRoute(group *gin.RouterGroup) {
	svcGroup := api.ApiGroupApp.K8SApiGroup.SVCApi
	// svc
	group.POST("/svc", svcGroup.CreateSOrUpdateVC)
	group.GET("/svc/:namespace", svcGroup.GetSVCListOrDetail)
	group.DELETE("/svc/:namespace/:name", svcGroup.DeleteSVC)
}
