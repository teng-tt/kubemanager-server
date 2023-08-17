package k8s

import (
	"github.com/gin-gonic/gin"
	node_req "kubemanager.com/model/node/request"
	"kubemanager.com/response"
)

type NodeApi struct {
}

func (n *NodeApi) UpdateTaints(ctx *gin.Context) {
	var updateTaint node_req.UpdatedTaint
	err := ctx.ShouldBind(&updateTaint)
	if err != nil {
		response.FailWithMessage(ctx, "参数解析报错！")
		return
	}
	err = nodeService.UpdateNodeTaint(updateTaint)
	if err != nil {
		response.FailWithMessage(ctx, "更新节点污点(Taint)报错, detail:"+err.Error())
	} else {
		response.Success(ctx)
	}
}

func (n *NodeApi) UpdateNode(ctx *gin.Context) {
	var updateLabel node_req.UpdatedLabel
	err := ctx.ShouldBind(&updateLabel)
	if err != nil {
		response.FailWithMessage(ctx, "参数解析报错！")
		return
	}
	err = nodeService.UpdateNodeLabel(updateLabel)
	if err != nil {
		response.FailWithMessage(ctx, "更新节点标签报错")
	} else {
		response.Success(ctx)
	}
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
