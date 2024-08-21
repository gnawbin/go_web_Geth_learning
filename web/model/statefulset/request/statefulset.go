package request

import (
	"k8s-web/model/base"
	pod_req "k8s-web/model/pod/request"
	pvc_req "k8s-web/model/pvc/request"
)

//@Author: morris

type StatefulSetBase struct {
	Name                 string                          `json:"name"`
	Namespace            string                          `json:"namespace"`
	Replicas             int32                           `json:"replicas"`
	Labels               []base.ListMapItem              `json:"labels"`
	Selector             []base.ListMapItem              `json:"selector"`
	ServiceName          string                          `json:"serviceName"`
	VolumeClaimTemplates []pvc_req.PersistentVolumeClaim `json:"volumeClaimTemplates"`
}
type StatefulSet struct {
	Base     StatefulSetBase `json:"base"`
	Template pod_req.Pod     `json:"template"`
}
