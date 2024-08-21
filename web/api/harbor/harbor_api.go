package harbor

import (
	"github.com/gin-gonic/gin"
	"k8s-web/global"
	"k8s-web/response"
	"strconv"
)

// @Author: morris
type HarborApi struct {
}

func (*HarborApi) GetArtifacts(c *gin.Context) {
	// 接收分页+模糊查询
	curPage := c.DefaultQuery("currentPage", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	newCurPage, _ := strconv.Atoi(curPage)
	newPageSize, _ := strconv.Atoi(pageSize)
	keyword := c.Query("keyword")
	// 调用api
	projectsPage := global.HarborClient.GetArtifacts(c.Param("projectName"),
		c.Param("repositoryName"),
		newCurPage, newPageSize, keyword)
	response.SuccessWithDetailed(c, "获取Artifacts成功", projectsPage)
}

func (*HarborApi) GetRepositories(c *gin.Context) {
	// 接收分页+模糊查询
	curPage := c.DefaultQuery("currentPage", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	newCurPage, _ := strconv.Atoi(curPage)
	newPageSize, _ := strconv.Atoi(pageSize)
	keyword := c.Query("keyword")
	// 调用api
	projectsPage := global.HarborClient.GetRepositories(c.Param("projectName"), newCurPage, newPageSize, keyword)
	response.SuccessWithDetailed(c, "获取Repositories成功", projectsPage)
}

func (*HarborApi) GetProjects(c *gin.Context) {
	// 接收分页+模糊查询
	curPage := c.DefaultQuery("currentPage", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	newCurPage, _ := strconv.Atoi(curPage)
	newPageSize, _ := strconv.Atoi(pageSize)
	keyword := c.Query("keyword")
	// 调用api
	projectsPage := global.HarborClient.GetProjects(newCurPage, newPageSize, keyword)
	response.SuccessWithDetailed(c, "获取Projects成功", projectsPage)
}

// 根据关键词推荐镜像
func (*HarborApi) MatchImage(c *gin.Context) {
	images := global.HarborClient.MatchImage(c.Query("keyword"))
	response.SuccessWithDetailed(c, "查询匹配镜像成功！", images)
}
