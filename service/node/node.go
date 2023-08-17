package node

import (
	"context"
	"encoding/json"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"kubemanager.com/global"
	node_req "kubemanager.com/model/node/request"
	node_res "kubemanager.com/model/node/response"
	"strings"
)

type NodeService struct {
}

func (n *NodeService) UpdateNodeTaint(updateTaint node_req.UpdatedTaint) error {
	patchData := map[string]any{
		"spec": map[string]any{
			"taints": updateTaint.Taints,
		},
	}
	patchDataBytes, _ := json.Marshal(&patchData)
	_, err := global.KubeConfigSet.CoreV1().Nodes().Patch(
		context.TODO(),
		updateTaint.Name,
		types.StrategicMergePatchType,
		patchDataBytes,
		metav1.PatchOptions{},
	)
	return err
}

func (n *NodeService) UpdateNodeLabel(updateLabel node_req.UpdatedLabel) error {
	labelsMap := make(map[string]string)
	for _, label := range updateLabel.Labels {
		labelsMap[label.Key] = label.Value
	}
	labelsMap["$patch"] = "replace"
	patchData := map[string]any{
		"metadata": map[string]any{
			"labels": labelsMap,
		},
	}
	patchDataBytes, _ := json.Marshal(&patchData)
	_, err := global.KubeConfigSet.CoreV1().Nodes().Patch(
		context.TODO(),
		updateLabel.Name,
		types.StrategicMergePatchType,
		patchDataBytes,
		metav1.PatchOptions{},
	)
	return err
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
