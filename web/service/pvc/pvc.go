package pvc

import (
	"context"
	"k8s-web/global"
	"k8s-web/model/base"
	pvc_req "k8s-web/model/pvc/request"
	pvc_res "k8s-web/model/pvc/response"
	"k8s-web/utils"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
	"strings"
)

// @Author: morris
type PVCService struct {
}

func (PVCService) GetPVCList(namespace string, keyword string) ([]pvc_res.PersistentVolumeClaim, error) {
	pvcResList := make([]pvc_res.PersistentVolumeClaim, 0)
	list, err := global.KubeConfigSet.CoreV1().PersistentVolumeClaims(namespace).
		List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, item := range list.Items {
		if !strings.Contains(item.Name, keyword) {
			continue
		}
		//item -> response
		matchLabels := make([]base.ListMapItem, 0)
		if item.Spec.Selector != nil {
			matchLabels = utils.ToList(item.Spec.Selector.MatchLabels)
		}
		pvcResItem := pvc_res.PersistentVolumeClaim{
			Name:      item.Name,
			Namespace: item.Namespace,
			Status:    item.Status.Phase,
			//转换为Mi
			Capacity:         int32(item.Spec.Resources.Requests.Storage().Value() / (1024 * 1024)),
			AccessModes:      item.Spec.AccessModes,
			StorageClassName: *item.Spec.StorageClassName,
			Age:              item.CreationTimestamp.UnixMilli(),
			Volume:           item.Spec.VolumeName,
			Labels:           utils.ToList(item.Labels),
			Selector:         matchLabels,
		}
		pvcResList = append(pvcResList, pvcResItem)
	}
	return pvcResList, err
}

func (PVCService) DeletePVC(namespace string, name string) error {
	err := global.KubeConfigSet.CoreV1().PersistentVolumeClaims(namespace).
		Delete(context.TODO(), name, metav1.DeleteOptions{})
	return err
}

func (PVCService) CreatePVC(pvcReq pvc_req.PersistentVolumeClaim) error {
	pvc := corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pvcReq.Name,
			Namespace: pvcReq.Namespace,
			Labels:    utils.ToMap(pvcReq.Labels),
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: utils.ToMap(pvcReq.Selector),
			},
			AccessModes: pvcReq.AccessModes,
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(strconv.Itoa(int(pvcReq.Capacity)) + "Mi"),
				},
			},
			StorageClassName: &pvcReq.StorageClassName,
		},
	}
	if pvc.Spec.StorageClassName != nil {
		pvc.Spec.Selector = nil
	}
	ctx := context.TODO()
	_, err := global.KubeConfigSet.CoreV1().PersistentVolumeClaims(pvc.Namespace).
		Create(ctx, &pvc, metav1.CreateOptions{})
	return err
}
