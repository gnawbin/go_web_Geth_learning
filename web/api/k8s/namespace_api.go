package k8s

import (
	"context"
	"github.com/gin-gonic/gin"
	"k8s-web/global"
	namespace_res "k8s-web/model/namespace/response"
	"k8s-web/response"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//@Author: morris

type NamespaceApi struct {
}

func (*NamespaceApi) GetNamespaceList(c *gin.Context) {
	ctx := context.TODO()
	list, err := global.KubeConfigSet.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	namespaceList := make([]namespace_res.Namespace, 0)
	for _, item := range list.Items {
		namespaceList = append(namespaceList, namespace_res.Namespace{
			Name:              item.Name,
			CreationTimestamp: item.CreationTimestamp.Unix(),
			Status:            string(item.Status.Phase),
		})
	}
	response.SuccessWithDetailed(c, "获取成功", namespaceList)
}
