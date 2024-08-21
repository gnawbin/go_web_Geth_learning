package harbor

import "github.com/gin-gonic/gin"

// @Author: morris
type HarborRouter struct {
}

func (HarborRouter) InitHarborRouter(r *gin.Engine) {
	group := r.Group("/harbor")
	initHarborRouter(group)
}
