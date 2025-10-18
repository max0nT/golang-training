package controller

import (
	"encoding/json"
	githubapi "src/src/github_api"

	"github.com/gin-gonic/gin"
)

// @Tags GitHub API
// @Body {object} githubapi.CreateRepoData
// @Router /api/v1/repo/create/ [post]
func (controller *Controller) CreateRepo(ctx *gin.Context) {
	token := controller.getToken(ctx)

	var jsonData map[string]any
	rawData, _ := ctx.GetRawData()
	parseErr := json.Unmarshal(rawData, &jsonData)

	if parseErr != nil {
		ctx.JSON(
			400,
			map[string]string{
				"attr": "invalid",
				"msg":  "Invalid format data",
			},
		)
	}

	var githubData *githubapi.CreateRepoData
	response, _ := githubapi.CreateRepo(
		token,
		githubData.FromMap(jsonData),
	)

	ctx.JSON(response.StatusCode, response.Body)

}
