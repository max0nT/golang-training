package main

import (
	"src/src/controller"
	"src/src/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Past GitHub Api key to perform request
func main() {
	server := gin.Default()
	v1 := server.Group("/api/v1/")

	docs.SwaggerInfo.Title = "Training API to interact with GitHub APi"
	docs.SwaggerInfo.Version = "1.0.0"

	apiController := controller.GetController()

	v1.GET("/me/", apiController.GetMe)
	v1.GET("/repo/detail/", apiController.GetDetailedRepoData)
	v1.GET("/repo/list/", apiController.GetRepoList)
	v1.POST("/repo/create/", apiController.CreateRepo)

	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	server.Run(":8001")
}
