package k8s

import (
	"github.com/gin-gonic/gin"
	"kubmanager/api"
)

func initRBACRouter(group *gin.RouterGroup) {
	rbacApiGroup := api.ApiGroupApp.K8SApiGroup.RbacApi
	group.GET("/sa/:namespace", rbacApiGroup.GetServiceAccountList)
	group.POST("/sa", rbacApiGroup.CreateServiceAccount)
	group.DELETE("/sa/:namespace/:name", rbacApiGroup.DeleteServiceAccount)

	// 角色管理
	group.GET("/role", rbacApiGroup.GetRoleDetailOrList)
	group.POST("/role", rbacApiGroup.CreateOrUpdateRole)
	group.DELETE("/role", rbacApiGroup.DeleteRole)

	// 账号角色绑定
	group.GET("/rb", rbacApiGroup.GetRbDetailOrList)
	group.POST("/rb", rbacApiGroup.CreateOrUpdateRb)
	group.DELETE("/rb", rbacApiGroup.DeleteRb)
}
