package main

import (
	"GinSwagger/docs"
	"GinSwagger/services"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

// @title Jewels API
// @version 1.0
// @description This is a sample practise for gin.
// @BasePath  /api/v1
func main() {
	r := gin.Default()

	r.POST("/issue/add", services.IssueAdd)
	r.GET("/issue/list", services.IssueList)
	//r.PUT("/issue/update", services.IssueUpdate)
	//r.DELETE("/issue/delete", services.IssueDelete)
	//r.GET("/issue/detail", services.IssueDetail)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	//swagger
	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "This is a sample server Petstore server."
	docs.SwaggerInfo.Version = "1.0"
	//docs.SwaggerInfo.Host = "petstore.swagger.io"
	//docs.SwaggerInfo.BasePath = ""  // xxx 覆盖上面的 BasePath
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run()
}
