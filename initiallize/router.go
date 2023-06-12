package initiallize

import (
	"github.com/gin-gonic/gin"
	"kubmanager/middleware"
	"kubmanager/router"
)

func Routers() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors)
	exampleRouter := router.RouterGroupApp.ExampleRouterGroup
	k8sRouter := router.RouterGroupApp.K8SRouterGroup
	exampleRouter.InitExample(r)
	k8sRouter.InitK8SRouter(r)

	return r
}
