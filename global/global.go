package global

import (
	"k8s.io/client-go/kubernetes"
	"kubemanager.com/config"
	"kubemanager.com/plugins/harbor"
)

var (
	CONF          config.Server
	KubeConfigSet *kubernetes.Clientset
	HarborClient  *harbor.Harbor
)
