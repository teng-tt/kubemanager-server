package harbor

import (
	"github.com/gin-gonic/gin"
	"kubmanager/api"
)

func initHarborRouter(group *gin.RouterGroup) {
	harborApiGroup := api.ApiGroupApp.HarborApiGroup.HarborApi
	group.GET("match", harborApiGroup.MatchImage)
	group.GET("/projects", harborApiGroup.GetProjects)
	group.GET("/projects/:projectName", harborApiGroup.GetRepositories)
	group.GET("/projects/:projectName/repositories/:repositoryName", harborApiGroup.GetArtifacts)
}
