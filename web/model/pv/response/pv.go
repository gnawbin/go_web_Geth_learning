package response

import (
	"k8s-web/model/base"
	corev1 "k8s.io/api/core/v1"
)

//@Author: morris

type PersistentVolume struct {
	Name string `json:"name"`
	//pv容量
	Capacity int32 `json:"capacity"`
	//ns 不必传
	//Namespace string             `json:"namespace"`
	Labels []base.ListMapItem `json:"labels"`
	//数据读写权限
	AccessModes []corev1.PersistentVolumeAccessMode `json:"accessModes"`
	//pv回收策略
	ReClaimPolicy corev1.PersistentVolumeReclaimPolicy `json:"reClaimPolicy"`
	//todo 待完善
	Status corev1.PersistentVolumePhase `json:"status"`
	//被具备某个pvc绑定
	Claim string `json:"claim"`
	//创建时间
	Age int64 `json:"age"`
	//状况描述
	Reason string `json:"reason"`
	//sc 名称
	StorageClassName string `json:"storageClassName"`
}
