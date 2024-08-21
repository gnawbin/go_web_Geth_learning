package k8s

import (
	"github.com/gin-gonic/gin"
	"k8s-web/api"
)

// @Author: morris
func initDaemonSetRouter(group *gin.RouterGroup) {
	daemonbsetApiGroup := api.ApiGroupApp.K8SApiGroup.DaemonsetApi
	group.POST("/daemonset", daemonbsetApiGroup.CreateOrUpdateDaemonSet)
	group.GET("/daemonset/:namespace", daemonbsetApiGroup.GetDaemonSetDetailOrList)
	group.DELETE("/daemonset/:namespace/:name", daemonbsetApiGroup.DeleteDaemonSet)
}
