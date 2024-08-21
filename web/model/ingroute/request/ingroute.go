package request

import (
	"k8s-web/model/base"
	ingroute_k8s "k8s-web/model/ingroute/k8s"
)

//@Author: morris

type IngressRoute struct {
	Name      string             `json:"name"`
	Namespace string             `json:"namespace"`
	Labels    []base.ListMapItem `json:"labels"`
	ingroute_k8s.IngressRouteSpec
}
