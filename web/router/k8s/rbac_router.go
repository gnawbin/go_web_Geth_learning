package k8s

import (
	"github.com/gin-gonic/gin"
	"k8s-web/api"
)

// @Author: morris
func initRBACRouter(group *gin.RouterGroup) {
	rbacApiGroup := api.ApiGroupApp.K8SApiGroup.RbacApi
	group.GET("/sa/:namespace", rbacApiGroup.GetServiceAccountList)
	group.POST("sa", rbacApiGroup.CreateServiceAccount)
	group.DELETE("sa/:namespace/:name", rbacApiGroup.DeleteServiceAccount)

	//角色管理
	group.GET("/role", rbacApiGroup.GetRoleDetailOrList)
	group.DELETE("/role", rbacApiGroup.DeleteRole)
	group.POST("/role", rbacApiGroup.CreateOrUpdateRole)

	//账号角色绑定
	group.GET("/rb", rbacApiGroup.GetRbDetailOrList)
	group.DELETE("/rb", rbacApiGroup.DeleteRb)
	group.POST("/rb", rbacApiGroup.CreateOrUpdateRb)
}
