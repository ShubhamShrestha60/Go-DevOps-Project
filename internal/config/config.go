package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	Env      string
	LogLevel string

	DB struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
		SSLMode  string
	}

	Auth struct {
		JWTSecret     string
		JWTExpiryH    int
		AdminPassword string
	}

	MigrationsPath string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := &Config{}

	cfg.Port = getEnv("PORT", "8080")
	cfg.Env = getEnv("ENV", "development")
	cfg.LogLevel = getEnv("LOG_LEVEL", "info")

	cfg.DB.Host = getEnv("DB_HOST", "localhost")
	cfg.DB.Port = getEnv("DB_PORT", "5432")
	cfg.DB.User = getEnv("DB_USER", "postgres")
	cfg.DB.Password = getEnv("DB_PASSWORD", "postgres")
	cfg.DB.Name = getEnv("DB_NAME", "devpulse")
	cfg.DB.SSLMode = getEnv("DB_SSLMODE", "disable")

	cfg.Auth.JWTSecret = getEnv("JWT_SECRET", "change-me")
	cfg.Auth.JWTExpiryH, _ = strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	cfg.Auth.AdminPassword = getEnv("ADMIN_PASSWORD", "")

	cfg.MigrationsPath = getEnv("MIGRATIONS_PATH", "migrations")

	return cfg
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
