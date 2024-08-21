package statefulset

import (
	"context"
	"k8s-web/global"
	pvc_req "k8s-web/model/pvc/request"
	statefulset_req "k8s-web/model/statefulset/request"
	statefulset_res "k8s-web/model/statefulset/response"
	"k8s-web/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
	"strings"
)

// @Author: morris
type StatefulSetService struct {
}

func (StatefulSetService) GetStatefulSetDetail(namespace, name string) (statefulset_req.StatefulSet, error) {
	var statefulSetReq statefulset_req.StatefulSet
	statefulSetK8s, err := global.KubeConfigSet.AppsV1().StatefulSets(namespace).
		Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return statefulSetReq, err
	}
	pvcReqList := make([]pvc_req.PersistentVolumeClaim, len(statefulSetK8s.Spec.VolumeClaimTemplates))
	for i, template := range statefulSetK8s.Spec.VolumeClaimTemplates {
		pvcReqList[i] = pvc_req.PersistentVolumeClaim{
			Name:             template.Name,
			AccessModes:      template.Spec.AccessModes,
			Capacity:         int32(template.Spec.Resources.Requests.Storage().Value() / (1024 * 1024)),
			StorageClassName: *template.Spec.StorageClassName,
		}
	}
	podReq := podConvert.PodK8s2Req(corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Labels: statefulSetK8s.Spec.Template.Labels,
		},
		Spec: statefulSetK8s.Spec.Template.Spec,
	})
	statefulSetReq = statefulset_req.StatefulSet{
		Base: statefulset_req.StatefulSetBase{
			Name:                 statefulSetK8s.Name,
			Namespace:            statefulSetK8s.Namespace,
			Replicas:             *statefulSetK8s.Spec.Replicas,
			Labels:               utils.ToList(statefulSetK8s.Labels),
			Selector:             utils.ToList(statefulSetK8s.Spec.Selector.MatchLabels),
			ServiceName:          statefulSetK8s.Spec.ServiceName,
			VolumeClaimTemplates: pvcReqList,
		},
		Template: podReq,
	}
	return statefulSetReq, err
}
func (StatefulSetService) GetStatefulSetList(namespace, keyword string) ([]statefulset_res.StatefulSet, error) {
	resList := make([]statefulset_res.StatefulSet, 0)
	list, err := global.KubeConfigSet.AppsV1().
		StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, item := range list.Items {
		if !strings.Contains(item.Name, keyword) {
			continue
		}
		resList = append(resList, statefulset_res.StatefulSet{
			Name:      item.Name,
			Namespace: item.Namespace,
			Ready:     item.Status.ReadyReplicas,
			Replicas:  item.Status.Replicas,
			Age:       item.CreationTimestamp.Unix(),
		})
	}
	return resList, err
}

func (StatefulSetService) DeleteStatefulSet(namespace, name string) error {
	return global.KubeConfigSet.AppsV1().StatefulSets(namespace).
		Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func (StatefulSetService) CreateOrUpdateStatefulSet(statefulSetReq statefulset_req.StatefulSet) error {
	pvcTemplates := make([]corev1.PersistentVolumeClaim, len(statefulSetReq.Base.VolumeClaimTemplates))
	for index, volumeClaimTemplate := range statefulSetReq.Base.VolumeClaimTemplates {
		pvcTemplates[index] = corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:   volumeClaimTemplate.Name,
				Labels: utils.ToMap(volumeClaimTemplate.Labels),
			},
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes: volumeClaimTemplate.AccessModes,
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceStorage: resource.MustParse(strconv.Itoa(int(volumeClaimTemplate.Capacity)) + "Mi"),
					},
				},
				StorageClassName: &volumeClaimTemplate.StorageClassName,
			},
		}
	}
	podK8s := podConvert.PodReq2K8s(statefulSetReq.Template)
	statefulSet := appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      statefulSetReq.Base.Name,
			Namespace: statefulSetReq.Base.Namespace,
			Labels:    utils.ToMap(statefulSetReq.Base.Labels),
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas:    &statefulSetReq.Base.Replicas,
			ServiceName: statefulSetReq.Base.ServiceName,
			Selector: &metav1.LabelSelector{
				MatchLabels: utils.ToMap(statefulSetReq.Base.Selector),
			},
			VolumeClaimTemplates: pvcTemplates,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: podK8s.ObjectMeta,
				Spec:       podK8s.Spec,
			},
		},
	}
	ctx := context.TODO()
	api := global.KubeConfigSet.AppsV1().StatefulSets(statefulSet.Namespace)
	statefulSetK8s, err := api.Get(ctx, statefulSet.Name, metav1.GetOptions{})
	if err != nil {
		_, err = api.
			Create(ctx, &statefulSet, metav1.CreateOptions{})
	} else {
		//防止服务抖动 id序号大的会先停止然后随即启动 -> id序号小的会先停止然后随即启动
		statefulSetK8s.Spec = statefulSet.Spec
		_, err = api.
			Update(ctx, statefulSetK8s, metav1.UpdateOptions{})
	}
	return err
}
