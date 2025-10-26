package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Controller struct{}

func (controller *Controller) getToken(ctx *gin.Context) string {
	token := ctx.GetHeader("Authorization")
	return token
}

func (controller *Controller) processResponse(response *http.Response) any {
	formatData := response.Header.Get("Content-Type")

	if !strings.Contains(formatData, "application/json") {
		errorMessage := fmt.Sprintf(
			"Response content is not json pls consider add processing this type: %s",
			formatData,
		)
		return map[string]any{
			"msg": errorMessage,
		}
	}

	bodyBytes, _ := io.ReadAll(response.Body)
	var jsonData any
	json.Unmarshal(bodyBytes, &jsonData)

	return jsonData

}

func GetController() *Controller {
	return &Controller{}
}
