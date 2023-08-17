package initiallize

import (
	"context"
	"fmt"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	kubeConfig := "E:/Goproject/src/kubemanager-server/.kube/config"
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
		K8S()
	}
}

func K8SWithToken() {
	cAData, err := ioutil.ReadFile("k8s_user/identity/ca.crt")
	if err != nil {
		panic(err)
	}
	config := &rest.Config{
		Host:            "https://192.168.2.11:6443",
		BearerTokenFile: "k8s_use/identity/token",
		TLSClientConfig: rest.TLSClientConfig{
			CAData: cAData,
		},
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	list, err := clientSet.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, item := range list.Items {
		fmt.Println(item.Namespace, item.Name)

	}
}
