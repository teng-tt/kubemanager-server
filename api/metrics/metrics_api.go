package metrics

import (
	"github.com/gin-gonic/gin"
	metricsResp "kubemanager.com/model/metrics/response"
	"kubemanager.com/response"
)

type MetricsApi struct {
}

func (m *MetricsApi) GetDashboardData(c *gin.Context) {
	cluster := metricsService.GetClusterInfo()
	resource := metricsService.GetResource()
	usage := metricsService.GetClusterUsage()
	usageRange := metricsService.GetClusterUsageRange()
	resultMap := make(map[string][]metricsResp.MetricsItem)
	resultMap["cluster"] = cluster
	resultMap["resource"] = resource
	resultMap["usage"] = usage
	resultMap["usageRange"] = usageRange
	response.SuccessWithDetailed(c, "获取Dashboard数据成功!", resultMap)
}
