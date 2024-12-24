package config

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ClerkSecretKey string
	DatabaseURL    string
	Port           string
}

func Load() *Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	config := &Config{
		ClerkSecretKey: getEnvOrDefault("CLERK_SECRET_KEY", ""),
		DatabaseURL: fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
			url.QueryEscape(getEnvOrDefault("DB_USER", "postgres")),
			url.QueryEscape(getEnvOrDefault("DB_PASSWORD", "")),
			getEnvOrDefault("DB_HOST", "localhost"),
			getEnvOrDefault("DB_PORT", "5432"),
			getEnvOrDefault("DB_NAME", "carpool")),
		//Port: getEnvOrDefault("PORT", "8080"),
	}

	// Validate required fields
	if config.ClerkSecretKey == "" {
		log.Fatal("CLERK_SECRET_KEY is required")
	}

	log.Printf("Attempting to connect to: postgresql://%s:****@%s:%s/%s",
		getEnvOrDefault("DB_USER", "postgres"),
		getEnvOrDefault("DB_HOST", "localhost"),
		getEnvOrDefault("DB_PORT", "5432"),
		getEnvOrDefault("DB_NAME", "carpool"))

	return config
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
