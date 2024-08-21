package deployment

import "k8s-web/convert"

// @Author: morris
type ServiceGroup struct {
	DeploymentService
}

var podConvert = convert.ConvertGroupApp.PodConvert
