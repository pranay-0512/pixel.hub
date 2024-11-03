package config

import (
	"log"
)

type Config struct {
	DatabaseURL string
	PORT        string
}

var AppConfig Config

func LoadConfig() {
	AppConfig = Config{
		DatabaseURL: "mongodb://localhost:27017",
		PORT:        "8080",
	}
	log.Println("Configuration loaded.")
}
