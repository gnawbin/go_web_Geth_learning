package request

import (
	"k8s-web/model/base"
	corev1 "k8s.io/api/core/v1"
)

// @Author: morris
type UpdatedLabel struct {
	Name   string             `json:"name"`
	Labels []base.ListMapItem `json:"labels"`
}

type UpdatedTaint struct {
	Name   string         `json:"name"`
	Taints []corev1.Taint `json:"taints"`
}
