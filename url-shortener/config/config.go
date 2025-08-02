package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	Port            string
	TrustedProxies  []string
	StaticFilesPath string
	IndexFilePath   string
	DatabaseMode    string // in_mem, redis, etc.
	RedisHost       string // default: localhost
	RedisPort       string // default: 6379
	RedisPassword   string // default: ""
	RedisDB         int    // default: 7
}

// LoadConfig loads configuration from environment variables with defaults
func LoadConfig() *Config {
	return &Config{
		Port:            getEnv("PORT", "8000"),
		TrustedProxies:  []string{"127.0.0.1"},
		StaticFilesPath: "./www",
		IndexFilePath:   "./www/index.html",
		DatabaseMode:    getEnv("DATABASE_MODE", "in_mem"),
		RedisHost:       getEnv("REDIS_HOST", "localhost"),
		RedisPort:       getEnv("REDIS_PORT", "6379"),
		RedisPassword:   getEnv("REDIS_PASSWORD", ""),
		RedisDB:         getEnvAsInt("REDIS_DB", 7),
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as int or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}
