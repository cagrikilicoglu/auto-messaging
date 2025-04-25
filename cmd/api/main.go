package main

import (
	"auto-messaging/config"
	"auto-messaging/internal/handler"
	"auto-messaging/internal/router"
	"auto-messaging/internal/service"
	"auto-messaging/pkg/database"
	"log"
	"strconv"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	db, err := database.NewPostgresDB(
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.DBName,
	)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize services
	messageService := service.NewMessageService(db)

	// Initialize handlers
	messageHandler := handler.NewMessageHandler(messageService)

	// Setup router
	r := router.SetupRouter(messageHandler)

	// Start server
	log.Printf("Server starting on port %d", cfg.Server.Port)
	if err := r.Run(":" + strconv.Itoa(cfg.Server.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
