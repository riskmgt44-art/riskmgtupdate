// config/config.go
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	MongoURI   string
	DBName     string
	JWTSecret  string
	Port       string
)

func LoadConfig() {
	// Load .env file if present (development)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	MongoURI = getEnv("MONGO_URI", "")
	if MongoURI == "" {
		log.Fatal("MONGO_URI is required")
	}

	DBName = getEnv("DB_NAME", "riskmgt")

	JWTSecret = getEnv("JWT_SECRET", "")
	if JWTSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	Port = getEnv("PORT", "8080")
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}