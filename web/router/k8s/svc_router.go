package k8s

import (
	"github.com/gin-gonic/gin"
	"k8s-web/api"
)

// @Author: morris
func initSvcRouter(group *gin.RouterGroup) {
	svcApiGroup := api.ApiGroupApp.K8SApiGroup.SvcApi
	group.POST("/svc", svcApiGroup.CreateOrUpdateSvc)
	group.GET("/svc/:namespace", svcApiGroup.GetSvcDetailOrList)
	group.DELETE("/svc/:namespace/:name", svcApiGroup.DeleteSvc)
}
