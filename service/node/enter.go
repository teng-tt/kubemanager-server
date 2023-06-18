package node

import "kubmanager/convert"

type NodeServiceGroup struct {
	NodeService NodeService
}

var nodeConvert = convert.ConvertGroupApp.NodeConver
