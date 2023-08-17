package request

import (
	"kubemanager.com/model/base"
	"kubemanager.com/model/ingroute"
)

type IngressRouteReq struct {
	Name      string             `json:"name"`
	Namespace string             `json:"namespace"`
	Labels    []base.ListMapItem `json:"labels"`
	ingroute.IngressRouteSpec
}
