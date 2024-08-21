package cronjob

import "k8s-web/convert"

// @Author: morris
type ServiceGroup struct {
	CronJobSetService
}

var podConvert = convert.ConvertGroupApp.PodConvert
