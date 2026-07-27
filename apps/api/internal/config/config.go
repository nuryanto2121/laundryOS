package config

import (
	"os"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	Port       string
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	RedisHost  string
	RedisPort  string
	AppVersion string
}

// Load reads configuration from environment variables.
// Each variable falls back to a safe default if not set.
func Load() *Config {
	return &Config{
		Port:       getEnv("PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     getEnv("DB_NAME", "laundryos"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		RedisHost:  getEnv("REDIS_HOST", "localhost"),
		RedisPort:  getEnv("REDIS_PORT", "6379"),
		AppVersion: getEnv("APP_VERSION", "0.1.0"),
	}
}

// getEnv returns the value of the environment variable named by key,
// or fallback if the variable is not set or is empty.
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
