package k8s

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"k8s-web/global"
	pod_req "k8s-web/model/pod/request"
	"k8s-web/response"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type PodApi struct {
}

// 因为update的字段属性有限 而我们实际更新过程当中 会修改定义的任一字段
func UpdatePod(ctx context.Context, pod *corev1.Pod) error {
	//update
	_, err := global.KubeConfigSet.CoreV1().Pods(pod.Namespace).Update(ctx, pod, metav1.UpdateOptions{})
	return err
}

func PatchPod(patchData map[string]interface{}, k8sPod *corev1.Pod, ctx context.Context) error {
	patchDataBytes, _ := json.Marshal(&patchData)
	_, err := global.KubeConfigSet.CoreV1().Pods(k8sPod.Namespace).Patch(
		ctx,
		k8sPod.Name,
		types.StrategicMergePatchType,
		patchDataBytes,
		metav1.PatchOptions{},
	)
	return err
}

func (*PodApi) DeletePod(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := podService.DeletePod(namespace, name)
	if err != nil {
		response.FailWithMessage(c, "删除Pod失败，detail："+err.Error())
	} else {
		response.Success(c)
	}
}

func (*PodApi) CreateOrUpdatePod(c *gin.Context) {
	var podReq pod_req.Pod
	if err := c.ShouldBind(&podReq); err != nil {
		response.FailWithMessage(c, "参数解析失败，detail："+err.Error())
		return
	}
	//校验必填项
	if err := podValidate.Validate(&podReq); err != nil {
		response.FailWithMessage(c, "参数验证失败，detail："+err.Error())
		return
	}
	if msg, err := podService.CreateOrUpdatePod(podReq); err != nil {
		response.FailWithMessage(c, msg)
	} else {
		response.SuccessWithMessage(c, msg)
	}
}

func (*PodApi) GetPodListOrDetail(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Query("name")
	keyword := c.Query("keyword")
	if name != "" {
		detail, err := podService.GetPodDetail(namespace, name)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		response.SuccessWithDetailed(c, "获取Pod详情成功", detail)
	} else {
		err, items := podService.GetPodList(namespace, keyword, c.Query("nodeName"))
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		response.SuccessWithDetailed(c, "获取Pod列表成功", items)
	}
}
