package service

import (
	"k8s-web/service/configmap"
	"k8s-web/service/cronjob"
	"k8s-web/service/daemonset"
	"k8s-web/service/deployment"
	"k8s-web/service/ingress"
	"k8s-web/service/ingroute"
	"k8s-web/service/job"
	"k8s-web/service/metrics"
	"k8s-web/service/node"
	"k8s-web/service/pod"
	"k8s-web/service/pv"
	"k8s-web/service/pvc"
	"k8s-web/service/rbac"
	"k8s-web/service/sc"
	"k8s-web/service/secret"
	"k8s-web/service/statefulset"
	"k8s-web/service/svc"
)

// @Author: morris
type ServiceGroup struct {
	PodServiceGroup         pod.PodServiceGroup
	NodeServiceGroup        node.Group
	ConfigMapServiceGroup   configmap.ServiceGroup
	SecretServiceGroup      secret.SeviceGroup
	PVServiceGroup          pv.ServiceGroup
	PVCServiceGroup         pvc.ServiceGroup
	SCServiceGroup          sc.SCServiceGroup
	SvcServiceGroup         svc.ServiceGroup
	IngressServiceGroup     ingress.ServiceGroup
	IngRouteServiceGroup    ingroute.ServiceGroup
	StatefulSetServiceGroup statefulset.ServiceGroup
	DeploymentServiceGroup  deployment.ServiceGroup
	DaemonSetServiceGroup   daemonset.ServiceGroup
	JobServiceGroup         job.ServiceGroup
	CronJobServiceGroup     cronjob.ServiceGroup
	RbacServiceGroup        rbac.RbacServiceGroup
	MetricsServiceGroup     metrics.MetricsServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
