package job

import "k8s-web/convert"

// @Author: morris
type ServiceGroup struct {
	JobSetService
}

var podConvert = convert.ConvertGroupApp.PodConvert
