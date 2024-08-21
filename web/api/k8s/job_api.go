package k8s

import (
	"github.com/gin-gonic/gin"
	job_req "k8s-web/model/job/request"
	"k8s-web/response"
)

// @Author: morris
type JobApi struct {
}

// 创建或更新deloyment
func (JobApi) CreateOrUpdateJob(c *gin.Context) {
	var jobReq job_req.Job
	err := c.ShouldBind(&jobReq)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	err = jobService.CreateOrUpdateJob(jobReq)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.Success(c)
}

// 查询Deloyment详情或列表
func (JobApi) GetJobDetailOrList(c *gin.Context) {
	//查询deployment 列表或详情
	namespace := c.Param("namespace")
	keyword := c.Param("keyword")
	name := c.Query("name")
	if name == "" {
		jobResList, err := jobService.GetJobList(namespace, keyword)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		response.SuccessWithDetailed(c, "获取job列表成功！", jobResList)
	} else {
		jobDetail, err := jobService.GetJobDetail(namespace, name)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		response.SuccessWithDetailed(c, "获取job详情成功！", jobDetail)
	}
}

// 删除Deloyment
func (JobApi) DeleteJob(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	err := jobService.DeleteJob(namespace, name)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.Success(c)
}
