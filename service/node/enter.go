package node

import "kubemanager.com/convert"

type NodeServiceGroup struct {
	NodeService NodeService
}

var nodeConvert = convert.ConvertGroupApp.NodeConver
