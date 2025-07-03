package api

import (
	"os"
)

// Config holds all configuration for the application
type Config struct {
	Port            string
	TrustedProxies  []string
	StaticFilesPath string
	IndexFilePath   string
	DatabaseMode    string // in_mem, redis, etc.
}

// LoadConfig loads configuration from environment variables with defaults
func LoadConfig() *Config {
	return &Config{
		Port:            getEnv("PORT", "8000"),
		TrustedProxies:  []string{"127.0.0.1"},
		StaticFilesPath: "./www",
		IndexFilePath:   "./www/index.html",
		DatabaseMode:    getEnv("DATABASE_MODE", "in_mem"),
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
