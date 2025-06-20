package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	Environment    string
	DefaultCurrency string
}

func Load() *Config {
	// Load .env file if exists
	godotenv.Load()

	return &Config{
		Port:           getEnvOrDefault("PORT", "8080"),
		Environment:    getEnvOrDefault("ENVIRONMENT", "development"),
		DefaultCurrency: getEnvOrDefault("DEFAULT_CURRENCY", "ARS"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}