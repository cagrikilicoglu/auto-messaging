package main

import (
	"auto-messaging/config"
	"auto-messaging/internal/client"
	"auto-messaging/internal/controller"
	"auto-messaging/internal/handler"
	"auto-messaging/internal/repository"
	"auto-messaging/internal/router"
	"auto-messaging/pkg/cache"
	"log"
	"os"
	"strconv"
)

func main() {
	// Initialize logger
	logger := log.New(os.Stdout, "AUTO-MSG: ", log.LstdFlags)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := repository.InitDB(cfg.DB)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize Redis cache
	redisCache := cache.NewRedisCache(
		cfg.Redis.Host,
		cfg.Redis.Port,
		cfg.Redis.Password,
		cfg.Redis.DB,
	)

	// Initialize repository
	messageRepo := repository.NewMessageRepository(db)

	// Initialize webhook client
	webhookClient := client.NewWebhookClient(cfg.Webhook.URL, cfg.Webhook.AuthKey)

	// Initialize controller with logger
	messageController := controller.NewMessageController(messageRepo, webhookClient, redisCache, logger)

	// Initialize handlers
	messageHandler := handler.NewMessageHandler(messageController)

	// Setup router
	r := router.SetupRouter(messageHandler)

	// Start server
	logger.Printf("Server starting on port %d", cfg.Server.Port)
	if err := r.Run(":" + strconv.Itoa(cfg.Server.Port)); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}
