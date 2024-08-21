package request

import (
	"k8s-web/model/base"
	pod_req "k8s-web/model/pod/request"
)

// @Author: morris
type DeploymentBase struct {
	Name      string             `json:"name"`
	Namespace string             `json:"namespace"`
	Replicas  int32              `json:"replicas"`
	Labels    []base.ListMapItem `json:"labels"`
	Selector  []base.ListMapItem `json:"selector"`
}
type Deployment struct {
	Base     DeploymentBase `json:"base"`
	Template pod_req.Pod    `json:"template"`
}
