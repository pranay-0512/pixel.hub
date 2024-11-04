package routes

import (
	"hub-api/controllers"
	"hub-api/middlewares"

	"github.com/gin-gonic/gin"
)

func SetUserRoutes(router *gin.RouterGroup) {
	authRoute := router.Group("/")
	{
		authRoute.POST("/signup", controllers.SignUp)
		authRoute.POST("/login", controllers.Login)
	}
	roomRoute := router.Group("/room")
	{
		roomRoute.POST("/create", middlewares.AuthMiddleware(), controllers.CreateRoom)
		roomRoute.GET("/list", middlewares.AuthMiddleware(), controllers.GetRooms)
		roomRoute.GET("/all", middlewares.AuthMiddleware(), controllers.GetAllRooms)
		roomRoute.GET("/:roomId", middlewares.AuthMiddleware(), controllers.GetRoom)
		roomRoute.PUT("/update/:roomId", middlewares.AuthMiddleware(), controllers.UpdateRoom)
	}
	mapRoute := router.Group("/map")
	{
		mapRoute.GET("/all", middlewares.AuthMiddleware(), controllers.GetAllMaps)
		mapRoute.GET("/:mapId", middlewares.AuthMiddleware(), controllers.GetMap)
	}
	eleRoute := router.Group("/element")
	{
		eleRoute.GET("/list", middlewares.AuthMiddleware(), controllers.GetElements)
		eleRoute.GET("/:elementId", middlewares.AuthMiddleware(), controllers.GetElement)
	}
}
