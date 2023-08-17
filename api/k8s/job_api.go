package k8s

import (
	"github.com/gin-gonic/gin"
	jobReqs "kubemanager.com/model/job/request"
	"kubemanager.com/response"
)

type JoApi struct {
}

func (j *JoApi) CreateOrUpdateJob(c *gin.Context) {
	var jobReq jobReqs.Job
	if err := c.ShouldBind(&jobReq); err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	err := jobService.CreateOrUpdateJob(jobReq)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.Success(c)
}

func (j *JoApi) DeleteJob(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	err := jobService.DeleteJob(namespace, name)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.Success(c)
}

func (j *JoApi) GetJobDetailOrList(c *gin.Context) {
	namespace := c.Param("namespace")
	keyword := c.Query("keyword")
	name := c.Query("name")
	if name == "" || len(name) == 0 {
		// 查询列表
		jobRespList, err := jobService.GetJobList(namespace, keyword)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		response.SuccessWithDetailed(c, "查询job列表成功!", jobRespList)
	} else {
		// 查询指定deployment
		jobRespInfo, err := jobService.GetJobDetail(name, namespace)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		response.SuccessWithDetailed(c, "查询job成功!", jobRespInfo)
	}
}
