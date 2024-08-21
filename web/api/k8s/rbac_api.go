package k8s

import (
	"github.com/gin-gonic/gin"
	rbac_req "k8s-web/model/rbac/request"
	"k8s-web/response"
)

// @Author: morris
type RbacApi struct {
}

func (RbacApi) CreateOrUpdateRb(c *gin.Context) {
	var rbReq rbac_req.RoleBinding
	if err := c.ShouldBind(&rbReq); err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	err := rbacService.CreateOrUpdateRb(rbReq)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.Success(c)
}
func (RbacApi) GetRbDetailOrList(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	keyword := c.Query("keyword")
	if name != "" {
		detail, err := rbacService.GetRbDetail(namespace, name)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		response.SuccessWithDetailed(c, "获取RoleBinding成功！", detail)
	} else {
		list, err := rbacService.GetRbList(namespace, keyword)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		response.SuccessWithDetailed(c, "获取RoleBinding列表成功！", list)
	}
}
func (RbacApi) DeleteRb(c *gin.Context) {
	err := rbacService.DeleteRb(c.Query("namespace"), c.Query("name"))
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.Success(c)
}

// 创建或更新Role
func (RbacApi) CreateOrUpdateRole(c *gin.Context) {
	var roleReq rbac_req.Role
	if err := c.ShouldBind(&roleReq); err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	err := rbacService.CreateOrUpdateRole(roleReq)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.Success(c)

}

// 查看Role详情或列表
func (RbacApi) GetRoleDetailOrList(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	keyword := c.Query("keyword")
	if name != "" {
		detail, err := rbacService.GetRoleDetail(namespace, name)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		response.SuccessWithDetailed(c, "查询Role/ClusterRole成功！", detail)
	} else {
		list, err := rbacService.GetRoleList(namespace, keyword)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		response.SuccessWithDetailed(c, "查询Role/ClusterRole列表成功！", list)
	}
}

// 删除Role
func (RbacApi) DeleteRole(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	err := rbacService.DeleteRole(namespace, name)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.Success(c)
}

// 创建Sa
func (RbacApi) CreateServiceAccount(c *gin.Context) {
	var saReq rbac_req.ServiceAccount
	if err := c.ShouldBind(&saReq); err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	err := rbacService.CreateServiceAccount(saReq)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.Success(c)
}

// 删除Sa
func (RbacApi) DeleteServiceAccount(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := rbacService.DeleteServiceAccount(namespace, name)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.Success(c)
}

// 查询Sa列表
func (RbacApi) GetServiceAccountList(c *gin.Context) {
	list, err := rbacService.GetServiceAccountList(c.Param("namespace"), c.Query("name"))
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.SuccessWithDetailed(c, "查询ServiceAccount列表成功！", list)
}
