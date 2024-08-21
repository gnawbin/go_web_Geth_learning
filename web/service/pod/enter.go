package pod

import "k8s-web/convert"

// @Author: morris
type PodServiceGroup struct {
	PodService
}

var podConvert = convert.ConvertGroupApp.PodConvert
