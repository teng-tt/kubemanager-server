package k8s

import (
	"github.com/gin-gonic/gin"
	conf_req "kubemanager.com/model/configmap/request"
	"kubemanager.com/response"
)

type ConfigMapApi struct {
}

// GetConfigMApDetailOrList 查询ConfigMap详情或列表
func (c *ConfigMapApi) GetConfigMApDetailOrList(ctx *gin.Context) {
	name := ctx.Query("name")
	namespace := ctx.Param("namespace")
	keyword := ctx.Query("keyword")

	// 名称为空，查询列表
	if name == "" {
		list, err := configMapService.GetConfigMapList(namespace, keyword)
		if err != nil {
			response.FailWithMessage(ctx, "查询ConfigMap列表失败！")
			return
		}
		response.SuccessWithDetailed(ctx, "查询ConfigMap列表成功", list)
		return
	} else {
		// 查询详情
		detail, err := configMapService.GetConfigMapDetail(namespace, keyword)
		if err != nil {
			response.FailWithMessage(ctx, "查询ConfigMap详情失败！")
			return
		}
		response.SuccessWithDetailed(ctx, "查询ConfigMap详情成功", detail)
	}
	return
}

// CreateOrUpdateConfigMap 创建或更新
func (c *ConfigMapApi) CreateOrUpdateConfigMap(ctx *gin.Context) {
	var configMapReq conf_req.ConfigMap
	err := ctx.ShouldBind(&configMapReq)
	if err != nil {
		response.FailWithMessage(ctx, "ConfigMap参数解析失败！")
		return
	}
	err = configMapService.CreateOrUpdateConfigMap(configMapReq)
	if err != nil {
		response.FailWithMessage(ctx, err.Error())
		return
	}
	response.Success(ctx)
}

// DeleteConfigMap 删除
func (c *ConfigMapApi) DeleteConfigMap(ctx *gin.Context) {
	err := configMapService.DeleteConfigMap(ctx.Param("namespace"), ctx.Param("name"))
	if err != nil {
		response.FailWithMessage(ctx, err.Error())
		return
	}
	response.Success(ctx)
}
