package k8s

import (
	"github.com/gin-gonic/gin"
	ingroute_req "k8s-web/model/ingroute/request"
	"k8s-web/response"
)

// @Author: morris
type IngRouteApi struct {
}

func (IngRouteApi) CreateOrUpdateIngRoute(c *gin.Context) {
	var ingressRouteReq ingroute_req.IngressRoute
	if err := c.ShouldBind(&ingressRouteReq); err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	err := ingRouteService.CreateOrUpdateIngRoute(ingressRouteReq)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.Success(c)
}

func (IngRouteApi) DeleteIngRoute(c *gin.Context) {
	err := ingRouteService.DeleteIngRoute(c.Param("namespace"), c.Param("name"))
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.Success(c)
}

// @https://github.com/kubernetes-client/python/blob/master/kubernetes/README.md
func (IngRouteApi) GetIngRouteDetailOrList(c *gin.Context) {
	namespace := c.Param("namespace")
	keyword := c.Query("keyword")
	name := c.Query("name")
	var data any
	var err error
	if name == "" {
		data, err = ingRouteService.GetIngRouteList(namespace, keyword)
	} else {
		data, err = ingRouteService.GetIngRouteDetail(namespace, name)
	}
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.SuccessWithDetailed(c, "查询成功！", data)
}

func (IngRouteApi) GetIngRouteMiddlewareList(c *gin.Context) {
	list, err := ingRouteService.GetIngRouteMiddlewareList(c.Param("namespace"))
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.SuccessWithDetailed(c, "查询成功！", list)
}
