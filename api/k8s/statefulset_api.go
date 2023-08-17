package k8s

import (
	"github.com/gin-gonic/gin"
	statefulsetReqs "kubemanager.com/model/statefulset/request"
	"kubemanager.com/response"
)

type StatefulSetApi struct {
}

func (s *StatefulSetApi) CreateOrUpdateStatefulSet(c *gin.Context) {
	var statefulsetReq statefulsetReqs.StatefulSte
	if err := c.ShouldBindUri(&statefulsetReq); err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	err := statefulsetService.CreateOrUpdateStatefulSet(statefulsetReq)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.Success(c)

}

func (s *StatefulSetApi) DeleteStatefulSet(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	err := statefulsetService.DeleteStatefulSet(name, namespace)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.Success(c)
}

func (s *StatefulSetApi) GetStatefulSetDetailOrList(c *gin.Context) {
	namespace := c.Param("namespace")
	keyword := c.Query("keyword")
	name := c.Query("name")
	if name == "" || len(name) == 0 {
		// 查询列表
		statefulRespList, err := statefulsetService.GetStatefulSetList(name, keyword)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		response.SuccessWithDetailed(c, "查询statefulSet列表成功!", statefulRespList)
	} else {
		// 查询指定statefulSet
		statefulSetRespInfo, err := statefulsetService.GetStatefulSetDetail(namespace, name)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		response.SuccessWithDetailed(c, "查询statefulSet详情成功!", statefulSetRespInfo)
	}
}
