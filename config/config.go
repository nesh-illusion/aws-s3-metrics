package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from .env file
func LoadEnv() {
	wd, _ := os.Getwd()
	log.Println("Current Working Directory:", wd)
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("No .env file found, using system envs")
	}
}

// GetEnv retrieves environment variables with fallback
func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
