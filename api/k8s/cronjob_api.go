package k8s

import (
	"github.com/gin-gonic/gin"
	cronjobReq "kubmanager/model/cronjob/request"
	"kubmanager/response"
)

type CronJobApi struct {
}

func (cr *CronJobApi) CreateOrUpdateCronjob(c *gin.Context) {
	var cronJobReq cronjobReq.CronJob
	if err := c.ShouldBindUri(&cronJobReq); err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	err := cronJobService.CreateOrUpdateCronJob(cronJobReq)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.Success(c)
}

func (cr *CronJobApi) DeleteCronjob(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	if err := cronJobService.DeleteCronJob(namespace, name); err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.Success(c)
}

func (cr *CronJobApi) GetCronjobDetailOrList(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Query("name")
	keyword := c.Param("keyword")
	if name == "" {
		jobResList, err := cronJobService.GetCronJobList(name, keyword)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		response.SuccessWithDetailed(c, "获取cronjob列表成功！", jobResList)
	} else {
		jobDetail, err := cronJobService.GetCronJobDetail(namespace, name)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		response.SuccessWithDetailed(c, "获取cronjob详情成功！", jobDetail)
	}
}
