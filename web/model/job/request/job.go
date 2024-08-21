package request

import (
	"k8s-web/model/base"
	pod_req "k8s-web/model/pod/request"
)

// @Author: morris
type JobBase struct {
	Name      string             `json:"name"`
	Namespace string             `json:"namespace"`
	Labels    []base.ListMapItem `json:"labels"`
	//Selector    []base.ListMapItem `json:"selector"`
	//jod的pod副本数，全部副本数运行成功，才能代表job运行成功
	Completions int32 `json:"completions"`
}
type Job struct {
	Base     JobBase     `json:"base"`
	Template pod_req.Pod `json:"template"`
}
