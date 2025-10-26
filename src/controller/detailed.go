package controller

import (
	githubapi "src/src/github_api"

	"github.com/gin-gonic/gin"
)

// @Tags GitHub API
// @Router /api/v1/repo/detail/	[get]
func (controller *Controller) GetDetailedRepoData(ctx *gin.Context) {
	token := controller.getToken(ctx)

	owner := ctx.Query("owner")
	name := ctx.Query("name")

	response, _ := githubapi.GetRepoDetailed(
		token,
		owner,
		name,
	)

	ctx.JSON(
		response.StatusCode,
		controller.processResponse(response),
	)
}
