package k8s

import (
	"github.com/gin-gonic/gin"
	cronjob_req "k8s-web/model/cronjob/request"
	"k8s-web/response"
)

// @Author: morris
type CronJobApi struct {
}

// 创建或更新deloyment
func (CronJobApi) CreateOrUpdateCronJob(c *gin.Context) {
	var cronJobReq cronjob_req.CronJob
	err := c.ShouldBind(&cronJobReq)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	err = cronJobService.CreateOrUpdateCronJob(cronJobReq)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.Success(c)
}

// 查询Deloyment详情或列表
func (CronJobApi) GetCronJobDetailOrList(c *gin.Context) {
	//查询deployment 列表或详情
	namespace := c.Param("namespace")
	keyword := c.Param("keyword")
	name := c.Query("name")
	if name == "" {
		jobResList, err := cronJobService.GetCronJobList(namespace, keyword)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		response.SuccessWithDetailed(c, "获取cronjob列表成功！", jobResList)
	} else {
		jobDetail, err := cronJobService.GetCronJobDetail(namespace, name)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		response.SuccessWithDetailed(c, "获取cronjob详情成功！", jobDetail)
	}
}

// 删除Deloyment
func (CronJobApi) DeleteCronJob(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	err := cronJobService.DeleteCronJob(namespace, name)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.Success(c)
}
