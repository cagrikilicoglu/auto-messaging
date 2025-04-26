package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

// DB holds database connection settings
type DB struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

// Redis holds Redis connection settings
type Redis struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// Server holds server settings
type Server struct {
	Port int
}

// Webhook holds webhook settings
type Webhook struct {
	URL     string
	AuthKey string
}

// Config holds all configuration settings
type Config struct {
	DB      DB
	Server  Server
	Webhook Webhook
	Redis   Redis
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("$HOME/go/src/auto-messaging/config")

	// Set environment variable prefix
	viper.SetEnvPrefix("")
	viper.AutomaticEnv()

	// Bind environment variables with proper structure
	viper.BindEnv("DB.Host", "DB_HOST")
	viper.BindEnv("DB.Port", "DB_PORT")
	viper.BindEnv("DB.User", "DB_USER")
	viper.BindEnv("DB.Password", "DB_PASSWORD")
	viper.BindEnv("DB.Name", "DB_NAME")

	viper.BindEnv("Redis.Host", "REDIS_HOST")
	viper.BindEnv("Redis.Port", "REDIS_PORT")
	viper.BindEnv("Redis.Password", "REDIS_PASSWORD")
	viper.BindEnv("Redis.DB", "REDIS_DB")

	viper.BindEnv("Server.Port", "SERVER_PORT")
	viper.BindEnv("Webhook.URL", "WEBHOOK_URL")
	viper.BindEnv("Webhook.AuthKey", "WEBHOOK_AUTH_KEY")

	// Set defaults
	viper.SetDefault("DB.Host", "localhost")
	viper.SetDefault("DB.Port", 5432)
	viper.SetDefault("DB.User", "postgres")
	viper.SetDefault("DB.Name", "auto_messaging")

	viper.SetDefault("Redis.Host", "localhost")
	viper.SetDefault("Redis.Port", 6379)
	viper.SetDefault("Redis.Password", "")
	viper.SetDefault("Redis.DB", 0)

	viper.SetDefault("Server.Port", 8080)

	// Try to read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		fmt.Printf("No config file found, using defaults and environment variables\n")
	}

	// Debug log the loaded values
	fmt.Printf("Config values from file:\n")
	fmt.Printf("DB Host: %s\n", viper.GetString("DB.Host"))
	fmt.Printf("DB Port: %d\n", viper.GetInt("DB.Port"))
	fmt.Printf("DB User: %s\n", viper.GetString("DB.User"))
	fmt.Printf("DB Password: %s\n", viper.GetString("DB.Password"))
	fmt.Printf("DB Name: %s\n", viper.GetString("DB.Name"))

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
