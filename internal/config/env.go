package config

import (
	"log"

	"github.com/joho/godotenv"
)

// LoadEnv loads variables from .env when present; it falls back to OS env vars.
func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found; relying on environment variables")
	}
}
