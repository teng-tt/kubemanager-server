package k8s

import (
	"github.com/gin-gonic/gin"
	pvc_req "kubmanager/model/pvc/request"
	"kubmanager/response"
)

type PVCApi struct {
}

func (p *PVCApi) CreatePVC(c *gin.Context) {
	var pvcReq pvc_req.PersistentVolumeClaim
	if err := c.ShouldBindUri(&pvcReq); err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	err := pvcService.CreatePVC(pvcReq)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.Success(c)
}

func (p *PVCApi) DeletePVC(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := pvcService.DeletePVC(namespace, name)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.Success(c)
}

func (p *PVCApi) GetPVCList(c *gin.Context) {
	namespace := c.Param("namespace")
	pvcResList, err := pvcService.GetPVCList(namespace, c.Query("keyword"))
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.SuccessWithDetailed(c, "获取数据成功！", pvcResList)
}
