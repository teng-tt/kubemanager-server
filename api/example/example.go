package example

import (
	"github.com/gin-gonic/gin"
	"kubemanager.com/response"
)

type ExampleApi struct {
}

func (e *ExampleApi) ExampleTest(c *gin.Context) {
	response.SuccessWithDetailed(c, "请求数据成功!", map[string]string{
		"message": "pong",
	})
}
