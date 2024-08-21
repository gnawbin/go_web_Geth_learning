package configmap

import "k8s-web/convert"

// @Author: morris
type ServiceGroup struct {
	ConfigMapService
}

var configConvert = convert.ConvertGroupApp.ConfigMapConvert
