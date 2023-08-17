package k8s

import (
	"github.com/gin-gonic/gin"
	"kubemanager.com/api"
)

type K8SRouter struct {
}

func (k *K8SRouter) InitK8SRouter(r *gin.Engine) {
	group := r.Group("/k8s")
	apiGroup := api.ApiGroupApp.K8SApiGroup
	// Pod
	group.GET("/namespace", apiGroup.GetNamespaceList)
	group.GET("/pod/:namespace", apiGroup.GetPodListOrDetail)
	group.POST("/pod", apiGroup.CreateOrUpdatePod)
	group.DELETE("/pod/:namespace/:name", apiGroup.DeletePod)
	// **************************************************************//
	// nodeScheduling
	group.GET("/node", apiGroup.GetNodeDetailOrList)
	group.POST("/node/label", apiGroup.UpdateNode)
	group.POST("/node/taint", apiGroup.UpdateTaints)
	// **************************************************************//
	// Config configMap/Secret
	group.GET("/configMap/:namespace", apiGroup.GetConfigMApDetailOrList)
	group.POST("/configMap/:namespace", apiGroup.CreateOrUpdateConfigMap)
	group.DELETE("/configMap", apiGroup.DeleteConfigMap)
	group.GET("/secret/:namespace", apiGroup.GetSecretListOrDetail)
	group.POST("/secret/:namespace", apiGroup.CreateOrUpdateSecret)
	group.DELETE("/secret", apiGroup.DeleteSecret)
	// ***************************************************************//
	// PV
	group.POST("/pv", apiGroup.CreatePV)
	group.GET("/pv/:namespace", apiGroup.GetPVList)
	group.DELETE("/pv/:namespace/:name", apiGroup.DeletePV)
	// PVC
	group.POST("/pvc", apiGroup.CreatePVC)
	group.GET("/pvc/:namespace", apiGroup.GetPVCList)
	group.DELETE("/pvc/:namespace/:name", apiGroup.DeletePVC)
	// SC
	// PVC
	group.POST("/sc", apiGroup.CreateSc)
	group.GET("/sc/:namespace", apiGroup.GetScList)
	group.DELETE("/sc/:namespace/:name", apiGroup.DeleteSc)
	// SVC
	initSvcRoute(group)
	// Ingress
	initIngressRouter(group)
	// IngressRouter
	initIngRouteRouter(group)
	// StatefulSetRouter
	initStatefulSetRouter(group)
	// DeploymentRouter
	initDeploymentRouter(group)
	// JobRouter
	initJobRouter(group)
	// CronJobRouter
	initCronjobRouter(group)
	// RBACRouter
	initRBACRouter(group)
}
