package controller

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

type Controller struct{}

func (controller *Controller) getToken(ctx *gin.Context) string {
	rawToken := ctx.GetHeader("Authorization")
	fmt.Println(ctx.Request.Header)
	token, isFound := strings.CutSuffix("Bearer ", rawToken)
	if !isFound {
		token = rawToken
	}
	return token
}

func GetController() *Controller {
	return &Controller{}
}
