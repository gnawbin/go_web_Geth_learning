package metrics

import (
	"github.com/gin-gonic/gin"
	metrics_res "k8s-web/model/metrics/response"
	"k8s-web/response"
)

// @Author: morris
type MetricsApi struct {
}

func (MetricsApi) GetDashboardData(c *gin.Context) {
	cluster := metricsService.GetClusterInfo()
	resource := metricsService.GetResource()
	usage := metricsService.GetClusterUsage()
	usageRange := metricsService.GetClusterUsageRange()
	resultMap := make(map[string][]metrics_res.MetricsItem)
	resultMap["cluster"] = cluster
	resultMap["resource"] = resource
	resultMap["usage"] = usage
	resultMap["usageRange"] = usageRange
	response.SuccessWithDetailed(c, "获取Dashboard数据成功！", resultMap)
}
