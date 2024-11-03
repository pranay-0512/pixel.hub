package db

import (
	"fmt"
	"hub-api/config"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(postgres.Open(config.AppConfig.DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	if err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}
	fmt.Println("Connected to database")
}
