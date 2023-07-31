package k8s

import (
	"github.com/gin-gonic/gin"
	"kubmanager/api"
)

func initIngRouteRouter(group *gin.RouterGroup) {
	ingRtGroup := api.ApiGroupApp.K8SApiGroup.IngressRouteApi
	// Ingress
	group.POST("/ingroute", ingRtGroup.CreateOrUpdateIngRoute)
	group.GET("/ingroute/:namespace", ingRtGroup.GetIngRouteDetailOrList)
	group.GET("/ingroute/:namespace/middleware", ingRtGroup.GetIngRouteMiddlewareList)
	group.DELETE("/ingroute/:namespace/:name", ingRtGroup.DeleteIngRoute)
}
