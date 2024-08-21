package statefulset

import "k8s-web/convert"

// @Author: morris
type ServiceGroup struct {
	StatefulSetService
}

var podConvert = convert.ConvertGroupApp.PodConvert
