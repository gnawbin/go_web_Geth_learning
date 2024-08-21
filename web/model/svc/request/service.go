package request

import (
	"k8s-web/model/base"
	corev1 "k8s.io/api/core/v1"
)

//@Author: morris

type Service struct {
	Name      string             `json:"name"`
	Namespace string             `json:"namespace"`
	Labels    []base.ListMapItem `json:"labels"`
	Type      corev1.ServiceType `json:"type"`
	Selector  []base.ListMapItem `json:"selector"`
	Ports     []ServicePort      `json:"ports"`
}

type ServicePort struct {
	Name       string `json:"name"`
	Port       int32  `json:"port"`
	TargetPort int32  `json:"targetPort"`
	NodePort   int32  `json:"nodePort"`
}
