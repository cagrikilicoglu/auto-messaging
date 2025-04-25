package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DB struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
	}
	Server struct {
		Port int
	}
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	cfg := &Config{}

	// Database configuration
	cfg.DB.Host = getEnv("DB_HOST", "localhost")
	cfg.DB.Port = getEnvAsInt("DB_PORT", 5432)
	cfg.DB.User = getEnv("DB_USER", "postgres")
	cfg.DB.Password = getEnv("DB_PASSWORD", "")
	cfg.DB.DBName = getEnv("DB_NAME", "auto_messaging")

	// Server configuration
	cfg.Server.Port = getEnvAsInt("API_PORT", 8080)

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
