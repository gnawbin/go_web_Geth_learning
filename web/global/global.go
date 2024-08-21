package global

import (
	"k8s-web/config"
	"k8s-web/plugins/harbor"
	"k8s.io/client-go/kubernetes"
)

var (
	CONF          config.Server
	KubeConfigSet *kubernetes.Clientset
	HarborClient  *harbor.Harbor
)
