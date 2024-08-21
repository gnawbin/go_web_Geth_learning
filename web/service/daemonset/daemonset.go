package daemonset

import (
	"context"
	"k8s-web/global"
	daemonset_req "k8s-web/model/daemonset/request"
	daemonset_res "k8s-web/model/daemonset/response"
	"k8s-web/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

// @Author: morris
type DaemonSetService struct {
}

func (DaemonSetService) CreateOrUpdateDaemonSet(daemonSetReq daemonset_req.DaemonSet) error {
	//转换为k8s结构
	podK8s := podConvert.PodReq2K8s(daemonSetReq.Template)
	daemonset := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      daemonSetReq.Base.Name,
			Namespace: daemonSetReq.Base.Namespace,
			Labels:    utils.ToMap(daemonSetReq.Base.Labels),
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: utils.ToMap(daemonSetReq.Base.Selector),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: podK8s.ObjectMeta,
				Spec:       podK8s.Spec,
			},
		},
	}
	ctx := context.TODO()
	daemonsetApi := global.KubeConfigSet.AppsV1().DaemonSets(daemonset.Namespace)
	daemonsetK8s, err := daemonsetApi.Get(ctx, daemonset.Name, metav1.GetOptions{})
	if err == nil {
		daemonsetK8s.Spec = daemonset.Spec
		_, err = daemonsetApi.Update(ctx, daemonsetK8s, metav1.UpdateOptions{})
	} else {
		_, err = daemonsetApi.Create(ctx, daemonset, metav1.CreateOptions{})
	}
	return err
}
func (DaemonSetService) DeleteDaemonSet(namespace, name string) error {
	return global.KubeConfigSet.AppsV1().DaemonSets(namespace).
		Delete(context.TODO(), name, metav1.DeleteOptions{})
}
func (DaemonSetService) GetDaemonSetDetail(namespace, name string) (daemonset_req.DaemonSet, error) {
	var daemonsetReq daemonset_req.DaemonSet
	daemonsetK8s, err := global.KubeConfigSet.AppsV1().DaemonSets(namespace).
		Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return daemonsetReq, err
	}
	podReq := podConvert.PodK8s2Req(corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Labels: daemonsetK8s.Spec.Template.Labels,
		},
		Spec: daemonsetK8s.Spec.Template.Spec,
	})
	daemonsetReq = daemonset_req.DaemonSet{
		Base: daemonset_req.DaemonSetBase{
			Name:      daemonsetK8s.Name,
			Namespace: daemonsetK8s.Namespace,
			Labels:    utils.ToList(daemonsetK8s.Labels),
			Selector:  utils.ToList(daemonsetK8s.Spec.Selector.MatchLabels),
		},
		Template: podReq,
	}
	return daemonsetReq, err
}
func (DaemonSetService) GetDaemonSetList(namespace, keyword string) ([]daemonset_res.DaemonSet, error) {
	daemonsetResList := make([]daemonset_res.DaemonSet, 0)
	list, err := global.KubeConfigSet.AppsV1().DaemonSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return daemonsetResList, err
	}
	for _, item := range list.Items {
		if !strings.Contains(item.Name, keyword) {
			continue
		}
		daemonsetResList = append(daemonsetResList, daemonset_res.DaemonSet{
			Name:      item.Name,
			Namespace: item.Namespace,
			Desired:   item.Status.DesiredNumberScheduled,
			Current:   item.Status.CurrentNumberScheduled,
			Ready:     item.Status.NumberReady,
			Available: item.Status.NumberAvailable,
			UpToDate:  item.Status.UpdatedNumberScheduled,
			Age:       item.CreationTimestamp.Unix(),
		})
	}
	return daemonsetResList, err
}
