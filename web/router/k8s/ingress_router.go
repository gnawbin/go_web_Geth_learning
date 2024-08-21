package k8s

import (
	"github.com/gin-gonic/gin"
	"k8s-web/api"
)

// @Author: morris
func initIngressRouter(group *gin.RouterGroup) {
	ingressApiGroup := api.ApiGroupApp.K8SApiGroup.IngressApi
	group.POST("/ingress", ingressApiGroup.CreateOrUpdateIngress)
	group.GET("/ingress/:namespace", ingressApiGroup.GetIngressDetailOrList)
	group.DELETE("/ingress/:namespace/:name", ingressApiGroup.DeleteIngress)
}
