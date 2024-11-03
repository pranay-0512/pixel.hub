package user

import (
	"hub-api/controllers"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.Engine) {
	auth := router.Group("/user/auth")
	{
		auth.POST("/signup", controllers.Signup)
		auth.POST("/login", controllers.Login)
	}
}
