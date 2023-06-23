package k8s

import (
	"github.com/gin-gonic/gin"
	req_secret "kubmanager/model/secret/request"
	"kubmanager/response"
)

type SecretApi struct {
}

func (s *SecretApi) CreateOrUpdateSecret(c *gin.Context) {
	var secretReq req_secret.Secret
	if err := c.ShouldBindUri(&secretReq); err != nil {
		response.FailWithMessage(c, "参数解析失败！")
		return
	}
	err := secretServicer.CreateOrUpdateSecret(secretReq)
	if err != nil {
		response.FailWithDetailed(c, "secret配置创建或更新失败！", err.Error())
		return
	}
	response.Success(c)
}

func (s *SecretApi) GetSecretListOrDetail(c *gin.Context) {
	name := c.Query("name")
	namespace := c.Param("namespace")
	keyword := c.Query("keyword")
	if name != "" {
		secretDetailRes, err := secretServicer.GetSecretDetail(namespace, namespace)
		if err != nil {
			response.SuccessWithDetailed(c, "查询secret失败！", err.Error())
			return
		}
		response.SuccessWithDetailed(c, "查询secret成功！", secretDetailRes)
		return
	}
	secretResList, err := secretServicer.GetSecretList(name, keyword)
	if err != nil {
		response.SuccessWithDetailed(c, "查询secret列表失败！", err.Error())
		return
	}
	response.SuccessWithDetailed(c, "查询secret列表成功！", secretResList)
	return
}

func (s *SecretApi) DeleteSecret(c *gin.Context) {
	err := secretServicer.DeleteSecret(c.Param("namespace"), c.Param("name"))
	if err != nil {
		response.FailWithDetailed(c, "secret配置删除失败！", err.Error())
		return
	}
	response.Success(c)
}
