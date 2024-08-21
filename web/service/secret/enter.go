package secret

import "k8s-web/convert"

// @Author: morris
type SeviceGroup struct {
	SecretService
}

var secretConvert = convert.ConvertGroupApp.SecretConvert
