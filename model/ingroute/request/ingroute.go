package request

import (
	"kubmanager/model/base"
	"kubmanager/model/ingroute"
)

type IngressRouteReq struct {
	Name      string             `json:"name"`
	Namespace string             `json:"namespace"`
	Labels    []base.ListMapItem `json:"labels"`
	ingroute.IngressRouteSpec
}
