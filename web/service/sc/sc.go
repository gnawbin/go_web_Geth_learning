package sc

import (
	"context"
	"fmt"
	"k8s-web/global"
	sc_req "k8s-web/model/sc/request"
	sc_res "k8s-web/model/sc/response"
	"k8s-web/utils"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

// @Author: morris
type SCService struct {
}

func (SCService) GetSCList(keyword string) ([]sc_res.StorageClass, error) {
	list, err := global.KubeConfigSet.StorageV1().StorageClasses().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	scResList := make([]sc_res.StorageClass, 0)
	for _, item := range list.Items {
		if !strings.Contains(item.Name, keyword) {
			continue
		}
		//item -> response
		var allowVolumeExpansion bool
		if item.AllowVolumeExpansion != nil {
			allowVolumeExpansion = *item.AllowVolumeExpansion
		}
		mountOptions := make([]string, 0)
		if item.MountOptions != nil {
			mountOptions = item.MountOptions
		}
		scResItem := sc_res.StorageClass{
			Name:                 item.Name,
			Labels:               utils.ToList(item.Labels),
			Provisioner:          item.Provisioner,
			MountOptions:         mountOptions,
			Parameters:           utils.ToList(item.Parameters),
			ReclaimPolicy:        *item.ReclaimPolicy,
			AllowVolumeExpansion: allowVolumeExpansion,
			Age:                  item.CreationTimestamp.UnixMilli(),
			VolumeBindingMode:    *item.VolumeBindingMode,
		}
		scResList = append(scResList, scResItem)
	}
	return scResList, err
}

func (SCService) DeleteSC(_ string, name string) error {
	return global.KubeConfigSet.StorageV1().StorageClasses().
		Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func (SCService) CreateSC(scReq sc_req.StorageClass) error {
	//判断Provisioner是否在系统支持
	provisionerList := strings.Split(global.CONF.System.Provisioner, ",")
	var flag bool
	for _, val := range provisionerList {
		if scReq.Provisioner == val {
			flag = true
			break
		}
	}
	if !flag {
		err := fmt.Errorf("%s 当前K8S未支持！ ", scReq.Provisioner)
		return err
	}
	sc := storagev1.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name:   scReq.Name,
			Labels: utils.ToMap(scReq.Labels),
		},
		Provisioner:          scReq.Provisioner,
		MountOptions:         scReq.MountOptions,
		VolumeBindingMode:    &scReq.VolumeBindingMode,
		ReclaimPolicy:        &scReq.ReclaimPolicy,
		AllowVolumeExpansion: &scReq.AllowVolumeExpansion,
		Parameters:           utils.ToMap(scReq.Parameters),
	}
	ctx := context.TODO()
	_, err := global.KubeConfigSet.StorageV1().StorageClasses().
		Create(ctx, &sc, metav1.CreateOptions{})
	return err
}
