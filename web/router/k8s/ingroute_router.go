package k8s

import (
	"github.com/gin-gonic/gin"
	"k8s-web/api"
)

// @Author: morris
func intIngRouteRouter(group *gin.RouterGroup) {
	ingRtApiGroup := api.ApiGroupApp.K8SApiGroup.IngRouteApi
	group.POST("/ingroute", ingRtApiGroup.CreateOrUpdateIngRoute)
	group.GET("/ingroute/:namespace", ingRtApiGroup.GetIngRouteDetailOrList)
	group.GET("/ingroute/:namespace/middleware", ingRtApiGroup.GetIngRouteMiddlewareList)
	group.DELETE("/ingroute/:namespace/:name", ingRtApiGroup.DeleteIngRoute)
}
