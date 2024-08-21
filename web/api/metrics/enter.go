package metrics

import "k8s-web/service"

//@Author: morris

type ApiGroup struct {
	MetricsApi
	PrometheusApi
}

var metricsService = service.ServiceGroupApp.MetricsServiceGroup.MetricsService
