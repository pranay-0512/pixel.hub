package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL  string
	PORT         string
	JWTSecret    string
	DatabaseName string
}

var AppConfig Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("error loading from .env")
	}
	AppConfig = Config{
		DatabaseURL:  os.Getenv("DATA_BASE_URL"),
		PORT:         os.Getenv("PORT"),
		JWTSecret:    os.Getenv("JWT_SECRET_KEY"),
		DatabaseName: os.Getenv("DATA_BASE_NAME"),
	}
	log.Println("Configuration loaded.")
}
