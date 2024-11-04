package routes

import (
	"hub-api/controllers"
	"hub-api/middlewares"

	"github.com/gin-gonic/gin"
)

func SetAdminRoutes(router *gin.RouterGroup) {
	authRoute := router.Group("/")
	{
		authRoute.POST("/signup", controllers.SignUp)
		authRoute.POST("/login", controllers.Login)
	}
	mapRoute := router.Group("/map")
	{
		mapRoute.POST("/create", middlewares.AuthMiddleware(), controllers.CreateMap)
		mapRoute.GET("/list", middlewares.AuthMiddleware(), controllers.GetMaps)
		mapRoute.PUT("/update/:mapId", middlewares.AuthMiddleware(), controllers.GetMaps)
		mapRoute.GET("/:mapId", middlewares.AuthMiddleware(), controllers.GetMap)
	}
	eleRoute := router.Group("/element")
	{
		eleRoute.POST("/create", middlewares.AuthMiddleware(), controllers.CreateElement)
		eleRoute.GET("/list", middlewares.AuthMiddleware(), controllers.GetElements)
		eleRoute.PUT("/update/:elementId", middlewares.AuthMiddleware(), controllers.UpdateElement)
		eleRoute.GET("/:elementId", middlewares.AuthMiddleware(), controllers.GetElement)
	}
}
