package config

import (
	"log"
)

type Config struct {
	DatabaseURL string
	PORT        string
	JWTSecret   string
}

var AppConfig Config

func LoadConfig() {
	AppConfig = Config{
		DatabaseURL: "mongodb://localhost:27017",
		PORT:        "8080",
		JWTSecret:   "wt2i895a8cwkomyR8rBUgw9xsNAVvZeG",
	}
	log.Println("Configuration loaded.")
}
