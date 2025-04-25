package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

// DBConfig holds database connection settings
type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

// RedisConfig holds Redis connection settings
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// ServerConfig holds server settings
type ServerConfig struct {
	Port int
}

// WebhookConfig holds webhook settings
type WebhookConfig struct {
	URL     string
	AuthKey string
}

// Config holds all configuration settings
type Config struct {
	DB      DBConfig
	Server  ServerConfig
	Webhook WebhookConfig
	Redis   RedisConfig
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// Set environment variable prefix
	viper.SetEnvPrefix("")
	viper.AutomaticEnv()

	// Bind environment variables
	viper.BindEnv("DB.Host", "DB_HOST")
	viper.BindEnv("DB.Port", "DB_PORT")
	viper.BindEnv("DB.User", "DB_USER")
	viper.BindEnv("DB.Password", "DB_PASSWORD")
	viper.BindEnv("DB.DBName", "DB_NAME")

	viper.BindEnv("Redis.Host", "REDIS_HOST")
	viper.BindEnv("Redis.Port", "REDIS_PORT")
	viper.BindEnv("Redis.Password", "REDIS_PASSWORD")
	viper.BindEnv("Redis.DB", "REDIS_DB")

	viper.BindEnv("Server.Port", "API_PORT")
	viper.BindEnv("Webhook.URL", "WEBHOOK_URL")
	viper.BindEnv("Webhook.AuthKey", "WEBHOOK_AUTH_KEY")

	// Set defaults
	viper.SetDefault("DB.Host", "localhost")
	viper.SetDefault("DB.Port", 5432)
	viper.SetDefault("DB.User", "postgres")
	viper.SetDefault("DB.Password", "postgres")
	viper.SetDefault("DB.DBName", "auto_messaging")

	viper.SetDefault("Redis.Host", "localhost")
	viper.SetDefault("Redis.Port", 6379)
	viper.SetDefault("Redis.Password", "")
	viper.SetDefault("Redis.DB", 0)

	viper.SetDefault("Server.Port", 8080)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
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
