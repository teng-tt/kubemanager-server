package k8s

import (
	"github.com/gin-gonic/gin"
	"kubemanager.com/api"
)

func initStatefulSetRouter(group *gin.RouterGroup) {
	ingStatefulSetGroup := api.ApiGroupApp.K8SApiGroup.StatefulSetApi
	group.POST("/statefulset", ingStatefulSetGroup.CreateOrUpdateStatefulSet)
	group.GET("/statefulset/:namespace", ingStatefulSetGroup.GetStatefulSetDetailOrList)
	group.DELETE("/statefulset/:namespace/:name", ingStatefulSetGroup.DeleteStatefulSet)
}
