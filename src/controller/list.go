package controller

import (
	githubapi "src/src/github_api"

	"github.com/gin-gonic/gin"
)

// @Tags GitHub API
// @Param  kind  path string true "User kind"
// @Param  name  path string true "User name"
// @Router /api/v1/repo/list/	[get]
func (controller *Controller) GetRepoList(ctx *gin.Context) {
	token := controller.getToken(ctx)

	kind := ctx.Query("kind")
	name := ctx.Query("name")

	response, _ := githubapi.GetRepoList(
		token,
		kind,
		name,
	)

	ctx.JSON(response.StatusCode, response.Body)
}
