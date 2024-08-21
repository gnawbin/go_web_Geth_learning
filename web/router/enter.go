package router

import (
	"k8s-web/router/example"
	"k8s-web/router/harbor"
	"k8s-web/router/k8s"
	"k8s-web/router/metrics"
)

type RouterGroup struct {
	ExampleRouterGroup example.ExampleRouter
	K8SRouterGroup     k8s.K8sRouter
	HarborRouterGroup  harbor.HarborRouter
	MetricsRouterGroup metrics.MetricsRouter
}

var RouterGroupApp = new(RouterGroup)
