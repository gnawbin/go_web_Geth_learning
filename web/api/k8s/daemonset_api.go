package k8s

import (
	"github.com/gin-gonic/gin"
	daemonset_req "k8s-web/model/daemonset/request"
	"k8s-web/response"
)

// @Author: morris
type DaemonsetApi struct {
}

// 创建或更新deloyment
func (DaemonsetApi) CreateOrUpdateDaemonSet(c *gin.Context) {
	var daemonsetReq daemonset_req.DaemonSet
	err := c.ShouldBind(&daemonsetReq)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	err = daemonsetService.CreateOrUpdateDaemonSet(daemonsetReq)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.Success(c)
}

// 查询Deloyment详情或列表
func (DaemonsetApi) GetDaemonSetDetailOrList(c *gin.Context) {
	//查询deployment 列表或详情
	namespace := c.Param("namespace")
	keyword := c.Param("keyword")
	name := c.Query("name")
	if name == "" {
		daemonsetResList, err := daemonsetService.GetDaemonSetList(namespace, keyword)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		response.SuccessWithDetailed(c, "获取daemonset列表成功！", daemonsetResList)
	} else {
		daemonsetDetail, err := daemonsetService.GetDaemonSetDetail(namespace, name)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		response.SuccessWithDetailed(c, "获取daemonset详情成功！", daemonsetDetail)
	}

}

// 删除Deloyment
func (DaemonsetApi) DeleteDaemonSet(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	err := daemonsetService.DeleteDaemonSet(namespace, name)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.Success(c)
}
