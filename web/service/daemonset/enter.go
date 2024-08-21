package daemonset

import "k8s-web/convert"

// @Author: morris
type ServiceGroup struct {
	DaemonSetService
}

var podConvert = convert.ConvertGroupApp.PodConvert
