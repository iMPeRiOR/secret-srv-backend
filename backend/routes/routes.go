package routes

import (
	"backend/controllers"
	"backend/docs"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup(router *gin.Engine) {
	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "This is a sample server Petstore server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "secret-server-backend.herokuapp.com"
	docs.SwaggerInfo.Schemes = []string{"https"}

	router.Use(static.Serve("/", static.LocalFile("./website/dist", true)))
	router.POST("/generate", controllers.GenerateToken)
	router.POST("/get/:token", controllers.GetToken)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
