package api

import (
	"k8s-web/api/example"
	"k8s-web/api/harbor"
	"k8s-web/api/k8s"
	"k8s-web/api/metrics"
)

type ApiGroup struct {
	ExampleApiGroup example.ApiGroup
	K8SApiGroup     k8s.ApiGroup
	HarborApiGroup  harbor.ApiGroup
	MetricsApiGroup metrics.ApiGroup
}

var ApiGroupApp = new(ApiGroup)
