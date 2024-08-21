package request

import (
	"k8s-web/model/base"
	corev1 "k8s.io/api/core/v1"
)

// @Author: morris
type Secret struct {
	Name      string             `json:"name"`
	Namespace string             `json:"namespace"`
	Type      corev1.SecretType  `json:"type"`
	Labels    []base.ListMapItem `json:"labels"`
	Data      []base.ListMapItem `json:"data"`
}
