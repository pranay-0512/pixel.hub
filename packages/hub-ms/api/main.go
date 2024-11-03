package main

import (
	"hub-api/config"
	"hub-api/db"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()

	db.InitDB()

	router := gin.Default()

	router.Run(":8080")
}
