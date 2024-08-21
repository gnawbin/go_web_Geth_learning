package configmap

import (
	configmapres "k8s-web/model/configmap/response"
	"k8s-web/utils"
	corev1 "k8s.io/api/core/v1"
)

//@Author: morris

type K82Res struct {
}

func (K82Res) GeCmReqItem(configMap corev1.ConfigMap) configmapres.ConfigMap {
	return configmapres.ConfigMap{
		Name:      configMap.Name,
		Namespace: configMap.Namespace,
		DataNum:   len(configMap.Data),
		Age:       configMap.CreationTimestamp.Unix(),
	}
}

func (this K82Res) GeCmReqDetail(configMap corev1.ConfigMap) configmapres.ConfigMap {
	detail := this.GeCmReqItem(configMap)
	detail.Labels = utils.ToList(configMap.Labels)
	detail.Data = utils.ToList(configMap.Data)
	return detail
}
