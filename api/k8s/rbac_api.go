package k8s

import (
	"github.com/gin-gonic/gin"
	rbacReqs "kubmanager/model/rbac/request"
	"kubmanager/response"
)

type RbacApi struct {
}

func (r *RbacApi) CreateServiceAccount(c *gin.Context) {
	var saReq rbacReqs.ServiceAccount
	if err := c.ShouldBind(&saReq); err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	err := rbacService.CreateServiceAccount(saReq)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.Success(c)
}

func (r *RbacApi) DeleteServiceAccount(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	err := rbacService.DeleteServiceAccount(name, namespace)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.Success(c)
}

func (r *RbacApi) GetServiceAccountList(c *gin.Context) {
	namespace := c.Param("namespace")
	keyword := c.Query("keyword")
	saResList, err := rbacService.GetServiceAccountList(namespace, keyword)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.SuccessWithDetailed(c, "获取ServiceAccount列表成功!", saResList)
}

func (r *RbacApi) CreateOrUpdateRole(c *gin.Context) {
	var roleReq rbacReqs.RoleReq
	if err := c.ShouldBindUri(&roleReq); err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	err := rbacService.CreateOrUpdateRole(roleReq)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.Success(c)
}

func (r *RbacApi) GetRoleDetailOrList(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	keyword := c.Query("keyword")
	if name != "" {
		// 查看详情
		roleResInfo, err := rbacService.GetRoleDetail(namespace, name)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		response.SuccessWithDetailed(c, "查询Role/ClusterRole成功!", roleResInfo)
	} else {
		// 查看列表
		resRoleList, err := rbacService.GetRoleList(namespace, keyword)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		response.SuccessWithDetailed(c, "获取Role/ClusterRole列表成功!", resRoleList)
		return
	}
}

func (r *RbacApi) DeleteRole(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	err := rbacService.DeleteRole(namespace, name)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.Success(c)
}

func (r *RbacApi) CreateOrUpdateRb(c *gin.Context) {
	var rbReq rbacReqs.RoleBinding
	if err := c.ShouldBind(&rbReq); err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	err := rbacService.CreateOrUpdateRb(rbReq)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.Success(c)
}

func (r *RbacApi) DeleteRb(c *gin.Context) {
	err := rbacService.DeleteRb(c.Query("namespace"), c.Query("name"))
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.Success(c)
}

func (r *RbacApi) GetRbDetailOrList(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	keyword := c.Query("keyword")
	if name != "" {
		detail, err := rbacService.GetRbDetail(namespace, name)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		response.SuccessWithDetailed(c, "获取RoleBinding成功！", detail)
	} else {
		list, err := rbacService.GetRbList(namespace, keyword)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		response.SuccessWithDetailed(c, "获取RoleBinding列表成功！", list)
	}
}
