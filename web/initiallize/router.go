package initiallize

import (
	"github.com/gin-gonic/gin"
	"k8s-web/middleware"
	"k8s-web/router"
)

func Routers() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors)
	examplGroup := router.RouterGroupApp.ExampleRouterGroup
	k8sGroup := router.RouterGroupApp.K8SRouterGroup
	harborRouterGroup := router.RouterGroupApp.HarborRouterGroup
	metricsRouterGroup := router.RouterGroupApp.MetricsRouterGroup
	examplGroup.InitExample(r)
	k8sGroup.InitK8SRouter(r)
	harborRouterGroup.InitHarborRouter(r)
	metricsRouterGroup.InitMetricsRouter(r)
	return r
	//r.Run(":8081") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
