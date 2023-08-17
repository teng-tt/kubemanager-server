package response

import "kubemanager.com/model/base"

type ConfigMap struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	DataNum   int    `json:"dataNum"`
	Age       int64  `json:"age"`
	// 查询configMap详情信息
	Data   []base.ListMapItem `json:"data"`
	Labels []base.ListMapItem `json:"labels"`
}
