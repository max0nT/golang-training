package controller

import (
	"fmt"
	"net/http"
	githubapi "src/src/github_api"
	"src/src/rabbitmq"

	"github.com/gin-gonic/gin"
)

// @Tags GitHub API
// @Summary Get User Info from GitHub API
// @Router /api/v1/me/	[get]
func (controller *Controller) GetMe(ctx *gin.Context) {
	token := controller.getToken(ctx)
	response, _ := githubapi.GetMeInfo(token)
	if response.StatusCode == http.StatusOK {
		go rabbitmq.SendSuccessfulMessage(
			fmt.Sprintf("Data from token %s is received", token),
		)
	}
	ctx.JSON(response.StatusCode, controller.processResponse(response))
}
