package k8s

import (
	"github.com/gin-gonic/gin"
	"kubmanager/response"
)

type NodeApi struct {
}

func (n *NodeApi) GetNodeDetailOrList(ctx *gin.Context) {
	keyword := ctx.Query("keyword")
	nodeName := ctx.Query("nodeName")
	if nodeName != "" {
		detail, err := nodeService.GteNodeDetail(nodeName)
		if err != nil {
			response.FailWithMessage(ctx, "查询node详情失败！")
		} else {
			response.FailWithDetailed(ctx, "查询node详情成功！", detail)
		}
	} else {
		list, err := nodeService.GetNodeList(keyword)
		if err != nil {
			response.FailWithMessage(ctx, "查询Node列表失败！")
		} else {
			response.SuccessWithDetailed(ctx, "查询Node列表成功！", list)
		}
	}
}
