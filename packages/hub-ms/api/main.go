package main

import (
	"hub-api/config"
	"hub-api/db"
	"hub-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()

	db.InitDB()

	router := gin.Default()

	adminRoute := router.Group("/api/v1/admin")
	userRoute := router.Group("/api/v1/user")

	routes.SetAdminRoutes(adminRoute)
	routes.SetUserRoutes(userRoute)
	router.Run(":8080")
}
