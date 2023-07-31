package k8s

import (
	"github.com/gin-gonic/gin"
	scReq "kubmanager/model/sc/request"
	"kubmanager/response"
)

type SCApi struct {
}

func (s *SCApi) CreateSc(c *gin.Context) {
	var scReqs scReq.StorageClass
	if err := c.ShouldBindUri(&scReqs); err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	err := scService.CreateSc(scReqs)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.Success(c)
}

func (s *SCApi) DeleteSc(c *gin.Context) {
	name := c.Param("name")
	err := scService.DeleteSc(name)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.Success(c)
}

func (s *SCApi) GetScList(c *gin.Context) {
	scRespList, err := scService.GetScList(c.Query("keyword"))
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.SuccessWithDetailed(c, "获取数据成功！", scRespList)
}
