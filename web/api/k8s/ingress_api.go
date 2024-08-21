package k8s

import (
	"github.com/gin-gonic/gin"
	ingress_req "k8s-web/model/ingress/request"
	"k8s-web/response"
)

// @Author: morris
type IngressApi struct {
}

func (IngressApi) CreateOrUpdateIngress(c *gin.Context) {
	ingressReq := ingress_req.Ingress{}
	err := c.ShouldBind(&ingressReq)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	err = ingressService.CreateOrUpdateIngress(ingressReq)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.Success(c)
}

func (IngressApi) DeleteIngress(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := ingressService.DeleteIngress(namespace, name)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.Success(c)
}

func (IngressApi) GetIngressDetailOrList(c *gin.Context) {
	//查询列表
	name := c.Query("name")
	namespace := c.Param("namespace")
	keyword := c.Query("keyword")
	if name == "" {
		ingressResList, err := ingressService.GetIngressList(namespace, keyword)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		response.SuccessWithDetailed(c, "获取Ingress列表成功", ingressResList)
	} else {
		detail, err := ingressService.GetIngressDetail(namespace, name)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		response.SuccessWithDetailed(c, "获取Ingress成功", detail)
	}

}
