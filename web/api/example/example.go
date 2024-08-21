package example

import (
	"github.com/gin-gonic/gin"
	"k8s-web/response"
)

type ExampleApi struct {
}

func (*ExampleApi) ExampleTest(c *gin.Context) {
	response.SuccessWithDetailed(c, "请求数据成功！", map[string]string{
		"message": "pong-v1.1",
	})
}
