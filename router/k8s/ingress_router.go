package k8s

import (
	"github.com/gin-gonic/gin"
	"kubemanager.com/api"
)

func initIngressRouter(group *gin.RouterGroup) {
	ingressGroup := api.ApiGroupApp.K8SApiGroup.IngressApi
	// Ingress
	group.POST("/ingress", ingressGroup.CreateOrUpdateIngress)
	group.GET("/ingress/:namespace", ingressGroup.GetIngressListOrDetail)
	group.DELETE("/ingress/:namespace/:name", ingressGroup.DeleteIngress)
}
