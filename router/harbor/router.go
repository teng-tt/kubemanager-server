package harbor

import "github.com/gin-gonic/gin"

type HarborRouter struct {
}

func (k *HarborRouter) InitKHarborRouter(r *gin.Engine) {
	group := r.Group("/harbor")
	initHarborRouter(group)
}
