package main

import (
	"hub-api/config"
	"hub-api/db"
	"hub-api/routes/admin"
	"hub-api/routes/user"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()

	db.InitDB()

	router := gin.Default()

	admin.SetupAdminRoutes(router)
	user.SetupUserRoutes(router)

	router.Run(":8080")
}
