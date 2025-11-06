package config

import (
	"os"
)

type Config struct {
	AppPort     string
	DatabaseURL string
	LogLevel    string
}

func Load() *Config {
	return &Config{
		AppPort:     getEnv("APP_PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://tasks:taskspass@db:5432/calendar_db?sslmode=disable"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return def
}
