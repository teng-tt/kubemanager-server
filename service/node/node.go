package node

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubmanager/global"
	node_res "kubmanager/model/node/response"
	"strings"
)

type NodeService struct {
}

func (n *NodeService) GteNodeDetail(nodeName string) (*node_res.Node, error) {
	nodeK8s, err := global.KubeConfigSet.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	detail := nodeConvert.K8s2ResConver.GetNodeDetail(*nodeK8s)
	return &detail, err
}

func (n *NodeService) GetNodeList(keyword string) ([]node_res.Node, error) {
	ctx := context.TODO()
	list, err := global.KubeConfigSet.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	nodeResList := make([]node_res.Node, 0)
	for _, item := range list.Items {
		if strings.Contains(item.Name, keyword) {
			nodeRes := nodeConvert.K8s2ResConver.GetNodeResItem(item)
			nodeResList = append(nodeResList, nodeRes)
		}
	}
	return nodeResList, err
}
