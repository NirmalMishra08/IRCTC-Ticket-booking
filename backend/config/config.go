package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT                string
	POSTGRES_CONNECTION string
	REDIS_DB_URL        string
	REDIS_PASSWORD      string
	STRIPE_KEY          string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment")
	}

	return &Config{
		PORT:                getEnv("PORT", "8080"),
		POSTGRES_CONNECTION: getEnv("POSTGRES_CONNECTION", ""),
		REDIS_DB_URL:        getEnv("REDIS_DB_URL", ""),
		REDIS_PASSWORD:      getEnv("REDIS_PASSWORD", ""),
		STRIPE_KEY:          getEnv("STRIPE_KEY", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
