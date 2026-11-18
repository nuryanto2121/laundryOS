package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	Port       string `env:"PORT"         env-default:"8080"`
	DBHost     string `env:"DB_HOST"      env-default:"localhost"`
	DBPort     string `env:"DB_PORT"      env-default:"5432"`
	DBName     string `env:"DB_NAME"      env-default:"laundryos"`
	DBUser     string `env:"DB_USER"      env-default:"postgres"`
	DBPassword string `env:"DB_PASSWORD"  env-default:""`
	RedisHost  string `env:"REDIS_HOST"   env-default:"localhost"`
	RedisPort  string `env:"REDIS_PORT"   env-default:"6379"`
	GatewayURL string `env:"GATEWAY_URL"  env-default:"http://localhost:3000"`
	StaffPhone string `env:"STAFF_PHONE"  env-default:""`
	AppVersion string `env:"APP_VERSION"  env-default:"0.1.0"`
}

// Load reads configuration from a .env file (if present in the current directory)
// and from environment variables. Environment variables always take precedence.
func Load() *Config {
	var cfg Config

	const envFile = ".env"

	if _, err := os.Stat(envFile); err == nil {
		// .env file found — read file first, then override with real env vars
		if err := cleanenv.ReadConfig(envFile, &cfg); err != nil {
			log.Fatalf("Config error: %s", err)
		}
	} else {
		// No .env file — read from environment variables only (production / Docker)
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			log.Fatalf("Config error: %s", err)
		}
	}

	return &cfg
}
