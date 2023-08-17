package example

import (
	"github.com/gin-gonic/gin"
	"kubemanager.com/api"
)

type ExampleRouter struct {
}

func (e *ExampleRouter) InitExample(r *gin.Engine) {
	group := r.Group("example")
	apiGroup := api.ApiGroupApp.ExampleApiGroup
	group.GET("/ping", apiGroup.ExampleTest)
}
