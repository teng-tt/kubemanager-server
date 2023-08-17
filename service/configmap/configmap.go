package configmap

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubemanager.com/global"
	conf_req "kubemanager.com/model/configmap/request"
	conf_res "kubemanager.com/model/configmap/response"
	"strings"
)

type ConfigMapService struct {
}

// GetConfigMapDetail 查询  configMap detail
func (c *ConfigMapService) GetConfigMapDetail(namespace, name string) (cm conf_res.ConfigMap, err error) {
	configMapK8s, err := global.KubeConfigSet.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return
	}
	cm = configMapConvert.K8s2ResConvert.GetCmReqDetail(*configMapK8s)
	return
}

// GetConfigMapList 查询 configMap List
func (c *ConfigMapService) GetConfigMapList(namespace, keyword string) (cmList []conf_res.ConfigMap, err error) {
	// 1. 从 k8s 查询
	list, errGetK8s := global.KubeConfigSet.CoreV1().ConfigMaps(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		err = errGetK8s
		return
	}
	configMapList := make([]conf_res.ConfigMap, 0)
	// 2. 转换为res(filter),按keyword过滤
	for _, item := range list.Items {
		if !strings.Contains(item.Name, keyword) {
			continue
		}
		resItem := configMapConvert.K8s2ResConvert.GetCmReqItem(item)
		configMapList = append(configMapList, resItem)
	}
	return configMapList, nil
}

func (c *ConfigMapService) CreateOrUpdateConfigMap(configReq conf_req.ConfigMap) error {
	// 将 request 转换为 k8s 结构
	configMapK8s := configMapConvert.Req2K8sConvert.CmReq2K8sConvert(configReq)
	ctx := context.TODO()
	// 查询是否存在
	_, errSearch := global.KubeConfigSet.CoreV1().ConfigMaps(configMapK8s.Namespace).Get(ctx, configMapK8s.Name, metav1.GetOptions{})
	// 存在更新
	if errSearch == nil {
		_, err := global.KubeConfigSet.CoreV1().ConfigMaps(configMapK8s.Namespace).Update(ctx, configMapK8s, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	} else {
		// 不存在创建
		_, err := global.KubeConfigSet.CoreV1().ConfigMaps(configMapK8s.Namespace).Create(ctx, configMapK8s, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteConfigMap 删除 configMap
func (c *ConfigMapService) DeleteConfigMap(ns, name string) error {
	return global.KubeConfigSet.CoreV1().ConfigMaps(ns).Delete(context.TODO(), name, metav1.DeleteOptions{})
}
