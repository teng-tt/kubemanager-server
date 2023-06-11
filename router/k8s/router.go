package k8s

import (
	"github.com/gin-gonic/gin"
	"kubmanager/api"
)

type K8SRouter struct {
}

func (k *K8SRouter) InitK8SRouter(r *gin.Engine) {
	group := r.Group("/k8s")
	apiGroup := api.ApiGroupApp.K8SApiGroup
	group.GET("/listPod", apiGroup.GetPodList)
}
