package metrics

import "github.com/gin-gonic/gin"

type MetricsRouter struct {
}

func (m *MetricsRouter) InitKMetricsRouter(r *gin.Engine) {
	group := r.Group("/metrics")
	initMetricsRouter(group)
}
