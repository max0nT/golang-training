package controller

import (
	githubapi "src/src/github_api"

	"github.com/gin-gonic/gin"
)

// @Tags GitHub API
// @Summary Get User Info from GitHub API
// @Router /api/v1/me/	[get]
func (controller *Controller) GetMe(ctx *gin.Context) {
	token := controller.getToken(ctx)
	response, _ := githubapi.GetMeInfo(token)
	ctx.JSON(response.StatusCode, controller.processResponse(response))
}
