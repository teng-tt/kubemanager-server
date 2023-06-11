package global

import (
	"k8s.io/client-go/kubernetes"
	"kubmanager/config"
)

var (
	CONF          config.Server
	KubeConfigSet *kubernetes.Clientset
)
