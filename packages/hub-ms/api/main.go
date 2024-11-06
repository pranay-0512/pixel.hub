package main

import (
	"context"
	"hub-api/config"
	"hub-api/db"
	"hub-api/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()

	ctx := context.Background()
	_, err := db.InitDB(ctx)
	if err != nil {
		log.Fatalln("error connecting to db")
		return
	}

	router := gin.Default()

	adminRoute := router.Group("/api/v1/admin")
	userRoute := router.Group("/api/v1/user")

	routes.SetAdminRoutes(adminRoute)
	routes.SetUserRoutes(userRoute)
	router.Run(":8080")
}
