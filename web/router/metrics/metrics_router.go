package metrics

import (
	"github.com/gin-gonic/gin"
	"k8s-web/api"
)

// @Author: morris
func initMetricsRouter(group *gin.RouterGroup) {
	metricsApiGroup := api.ApiGroupApp.MetricsApiGroup.MetricsApi
	prometheusApiGroup := api.ApiGroupApp.MetricsApiGroup.PrometheusApi
	group.GET("/dashboard", metricsApiGroup.GetDashboardData)
	group.GET("/prometheus", prometheusApiGroup.GetMetrics)
}
