package convert

import (
	"k8s-web/convert/configmap"
	"k8s-web/convert/node"
	"k8s-web/convert/pod"
	"k8s-web/convert/secret"
)

//@Author: morris

type ConvertGroup struct {
	PodConvert       pod.PodConvertGroup
	NodeConvert      node.Group
	ConfigMapConvert configmap.ConvertGroup
	SecretConvert    secret.ConvertGroup
}

var ConvertGroupApp = new(ConvertGroup)
