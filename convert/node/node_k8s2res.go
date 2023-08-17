package node

import (
	corev1 "k8s.io/api/core/v1"
	"kubemanager.com/model/base"
	node_res "kubemanager.com/model/node/response"
)

type NodeK8s2Res struct {
}

func (n *NodeK8s2Res) getNodeStatus(nodeConditions []corev1.NodeCondition) string {
	nodeStatus := "NotReady"
	for _, condition := range nodeConditions {
		if condition.Type == "Ready" && condition.Status == "True" {
			nodeStatus = "Ready"
			break
		}
	}
	return nodeStatus
}

func (n *NodeK8s2Res) getNodeIp(addresses []corev1.NodeAddress, addressType corev1.NodeAddressType) string {
	for _, item := range addresses {
		if item.Type == addressType {
			return item.Address
		}
	}
	return "<none>"
}

func (n *NodeK8s2Res) mapToList(m map[string]string) []base.ListMapItem {
	res := make([]base.ListMapItem, 0)
	for k, v := range m {
		res = append(res, base.ListMapItem{
			Key:   k,
			Value: v,
		})
	}
	return res
}

func (n *NodeK8s2Res) GetNodeDetail(nodeK8s corev1.Node) node_res.Node {
	nodeRes := n.GetNodeResItem(nodeK8s)
	// 计算label 和taint
	nodeRes.Taints = nodeK8s.Spec.Taints
	nodeRes.Labels = n.mapToList(nodeK8s.Labels)
	return nodeRes
}

func (n *NodeK8s2Res) GetNodeResItem(nodeK8s corev1.Node) node_res.Node {
	nodeInfo := nodeK8s.Status.NodeInfo
	return node_res.Node{
		Name:             nodeK8s.Name,
		Status:           n.getNodeStatus(nodeK8s.Status.Conditions),
		Age:              nodeK8s.CreationTimestamp.Unix(),
		InternalIp:       n.getNodeIp(nodeK8s.Status.Addresses, corev1.NodeInternalIP),
		ExternalIp:       n.getNodeIp(nodeK8s.Status.Addresses, corev1.NodeExternalIP),
		OsImage:          nodeInfo.OSImage,
		Version:          nodeInfo.KubeletVersion,
		KernelVersion:    nodeInfo.KernelVersion,
		ContainerRuntime: nodeInfo.ContainerRuntimeVersion,
	}
}
