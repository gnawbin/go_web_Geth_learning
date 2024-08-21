package k8s

import (
	"github.com/gin-gonic/gin"
	"k8s-web/api"
)

type K8sRouter struct {
}

func (*K8sRouter) InitK8SRouter(r *gin.Engine) {
	group := r.Group("/k8s")
	apiGroup := api.ApiGroupApp.K8SApiGroup
	group.POST("/pod", apiGroup.CreateOrUpdatePod)
	group.GET("/pod/:namespace", apiGroup.GetPodListOrDetail)
	group.DELETE("/pod/:namespace/:name", apiGroup.DeletePod)
	group.GET("/namespace", apiGroup.GetNamespaceList)
	///////////////////

	//nodeScheduling
	group.GET("/node", apiGroup.GetNodeDetailOrList)
	group.PUT("/node/label", apiGroup.UpdateNodeLabel)
	group.PUT("/node/taint", apiGroup.UpdateNodeTaint)

	//******************ConfigMap************************//
	group.POST("/configmap", apiGroup.CreateOrUpdateConfigMap)
	group.GET("/configmap/:namespace", apiGroup.GetConfigMapDetailOrList)
	group.DELETE("/configmap/:namespace/:name", apiGroup.DeleteConfigMap)

	//*******************Secret***********************//
	group.POST("/secret", apiGroup.CreateOrUpdateSecret)
	group.GET("/secret/:namespace", apiGroup.GetSecretDetailOrList)
	group.DELETE("/secret/:namespace/:name", apiGroup.DeleteSecret)

	//*******************PV***********************//
	group.POST("/pv", apiGroup.CreatePV)
	group.GET("/pv/:namespace", apiGroup.GetPVList)
	group.DELETE("/pv/:namespace/:name", apiGroup.DeletePV)

	//*******************PVC***********************//
	group.POST("/pvc", apiGroup.CreatePVC)
	group.GET("/pvc/:namespace", apiGroup.GetPVCList)
	group.DELETE("/pvc/:namespace/:name", apiGroup.DeletePVC)

	//*******************SC***********************//
	group.POST("/sc", apiGroup.CreateSC)
	group.GET("/sc/:namespace", apiGroup.GetSCList)
	group.DELETE("/sc/:namespace/:name", apiGroup.DeleteSC)

	initSvcRouter(group)
	initIngressRouter(group)
	intIngRouteRouter(group)
	initStatefulSetRouter(group)
	initDeloymentRouter(group)
	initDaemonSetRouter(group)
	initJobRouter(group)
	initCronJobRouter(group)
	initRBACRouter(group)
}
