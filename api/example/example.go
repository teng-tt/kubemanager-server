package example

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ExampleApi struct {
}

func (e *ExampleApi) ExampleTest(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
