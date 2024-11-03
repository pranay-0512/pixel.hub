package config

import (
	"log"
)

type Config struct {
	DatabaseURL string
	PORT        string
	DSN         string
}

var AppConfig Config

func LoadConfig() {
	AppConfig = Config{
		DatabaseURL: "postgres://postgres:qwerty@localhost:5432/pixelhub?sslmode=disable",
		PORT:        "8080",
		DSN:         "host=localhost user=postgres password=qwerty dbname=pixelhub port=5432 sslmode=disable",
	}
	log.Println("Configuration loaded.")
}
