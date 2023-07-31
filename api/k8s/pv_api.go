package k8s

import (
	"github.com/gin-gonic/gin"
	pv_req "kubmanager/model/pv/request"
	"kubmanager/response"
)

type PVApi struct {
}

func (p *PVApi) CreatePV(c *gin.Context) {
	var pvReq pv_req.PersistentVolume
	if err := c.ShouldBindUri(&pvReq); err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	err := pvService.CreatePV(pvReq)
	if err != nil {
		response.FailWithMessage(c, err.Error())
	} else {
		response.Success(c)
	}
}

func (p *PVApi) DeletePV(c *gin.Context) {
	err := pvService.DeletePV(c.Param("namespace"), c.Param("name"))
	if err != nil {
		response.FailWithMessage(c, err.Error())
	} else {
		response.Success(c)
	}
}

func (p *PVApi) GetPVList(c *gin.Context) {
	pvResList, err := pvService.GetPvList(c.Query("keyword"))
	if err != nil {
		response.FailWithMessage(c, err.Error())
	} else {
		response.SuccessWithDetailed(c, "获取数据成功", pvResList)
	}
}
