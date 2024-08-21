package metrics

import "github.com/gin-gonic/gin"

// @Author: morris
type MetricsRouter struct {
}

func (MetricsRouter) InitMetricsRouter(r *gin.Engine) {
	group := r.Group("/metrics")
	initMetricsRouter(group)
}
