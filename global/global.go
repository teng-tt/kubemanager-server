package global

import (
	"k8s.io/client-go/kubernetes"
	"kubmanager/config"
	"kubmanager/plugins/harbor"
)

var (
	CONF          config.Server
	KubeConfigSet *kubernetes.Clientset
	HarborClient  *harbor.Harbor
)
