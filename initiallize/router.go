package initiallize

import (
	"github.com/gin-gonic/gin"
	"kubemanager.com/middleware"
	"kubemanager.com/router"
)

func Routers() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors)
	exampleRouter := router.RouterGroupApp.ExampleRouterGroup
	k8sRouter := router.RouterGroupApp.K8SRouterGroup
	harborRouter := router.RouterGroupApp.HarborRouterGroup
	metricsRouter := router.RouterGroupApp.MetricsRouterGroup
	exampleRouter.InitExample(r)
	k8sRouter.InitK8SRouter(r)
	harborRouter.InitKHarborRouter(r)
	metricsRouter.InitKMetricsRouter(r)

	return r
}
