package k8s

import (
	"github.com/gin-gonic/gin"
	deployment_req "k8s-web/model/deployment/request"
	"k8s-web/response"
)

// @Author: morris
type DeploymentApi struct {
}

// 创建或更新deloyment
func (DeploymentApi) CreateOrUpdateDeployment(c *gin.Context) {
	var deploymentReq deployment_req.Deployment
	err := c.ShouldBind(&deploymentReq)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	err = deploymentService.CreateOrUpdateDeployment(deploymentReq)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.Success(c)
}

// 查询Deloyment详情或列表
func (DeploymentApi) GetDeploymentDetailOrList(c *gin.Context) {
	//查询deployment 列表或详情
	namespace := c.Param("namespace")
	keyword := c.Query("keyword")
	name := c.Query("name")
	if name == "" {
		deploymentResList, err := deploymentService.GetDeploymentList(namespace, keyword)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		response.SuccessWithDetailed(c, "获取deployment列表成功！", deploymentResList)
	} else {
		deloymentDetail, err := deploymentService.GetDeploymentDetail(namespace, name)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		response.SuccessWithDetailed(c, "获取deployment详情成功！", deloymentDetail)
	}

}

// 删除Deloyment
func (DeploymentApi) DeleteDeployment(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	err := deploymentService.DeleteDeployment(namespace, name)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.Success(c)
}
