package initiallize

import (
	"io/ioutil"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"kubemanager.com/global"
	"os"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

func K8S() {
	kubeConfig := global.CONF.System.K8sConfig.KubeConfig
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	global.KubeConfigSet = clientSet
}

// 判断集群环境
func isInCluster() (isInCluster bool) {
	tokenFile := "/var/run/secrets/kubernetes.io/serviceaccount/token"
	_, err := os.Stat(tokenFile)
	if err == nil {
		isInCluster = true
	}
	return
}

// K8SWithDiscovery 根据集群环境判断使用集群内认证还是集群外认证
func K8SWithDiscovery() {
	if isInCluster() {
		// 集群内环境api方法已经预置好了证书，token
		config, err := rest.InClusterConfig()
		if err != nil {
			panic(err)
		}
		clientSet, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err)
		}
		global.KubeConfigSet = clientSet
	} else {
		kubeConfig := global.CONF.System.K8sConfig.KubeConfig
		if len(kubeConfig) > 0 && kubeConfig != "" {
			K8S()
		} else {
			K8SWithToken()
		}
	}
}

func K8SWithToken() {
	cAData, err := ioutil.ReadFile(global.CONF.System.K8sConfig.CacertPath)
	if err != nil {
		panic(err)
	}
	config := &rest.Config{
		Host:            global.CONF.System.K8sConfig.Host,
		BearerTokenFile: global.CONF.System.K8sConfig.TokenFile,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: cAData,
		},
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	global.KubeConfigSet = clientSet
}
