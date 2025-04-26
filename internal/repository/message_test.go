package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"auto-messaging/internal/model"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	ctx := context.Background()

	// Start PostgreSQL container
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "test",
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections"),
			wait.ForListeningPort("5432/tcp"),
		),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start container: %v", err)
	}
	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("Failed to terminate container: %v", err)
		}
	})

	// Get container host and port
	mappedPort, err := container.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("Failed to get container port: %v", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get container host: %v", err)
	}

	// Wait a bit to ensure the database is ready
	time.Sleep(2 * time.Second)

	// Connect to the database
	dsn := fmt.Sprintf("host=%s port=%s user=test password=test dbname=test sslmode=disable",
		host, mappedPort.Port())
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate the Message model
	if err := db.AutoMigrate(&model.Message{}); err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func TestMessageRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMessageRepository(db)

	message := &model.Message{
		Content:     "Test message",
		To:          "test@example.com",
		Status:      model.MessageStatusPending,
		ScheduledAt: time.Now().Add(1 * time.Hour),
	}

	if err := repo.Create(context.Background(), message); err != nil {
		t.Errorf("Create() error = %v", err)
	}

	// Verify the message was created
	var count int64
	if err := db.Model(&model.Message{}).Count(&count).Error; err != nil {
		t.Errorf("Failed to count messages: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 message, got %d", count)
	}
}

func TestMessageRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMessageRepository(db)

	// Create a test message
	message := &model.Message{
		Content:     "Test message",
		To:          "test@example.com",
		Status:      model.MessageStatusPending,
		ScheduledAt: time.Now().Add(1 * time.Hour),
	}
	if err := repo.Create(context.Background(), message); err != nil {
		t.Fatalf("Failed to create test message: %v", err)
	}

	// Test finding the message
	found, err := repo.FindByID(context.Background(), message.ID)
	if err != nil {
		t.Errorf("FindByID() error = %v", err)
	}
	if found == nil {
		t.Error("Expected to find message, got nil")
	}
	if found.Content != message.Content {
		t.Errorf("Expected content %q, got %q", message.Content, found.Content)
	}
}

func TestMessageRepository_UpdateStatus(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMessageRepository(db)

	// Create a test message
	message := &model.Message{
		Content:     "Test message",
		To:          "test@example.com",
		Status:      model.MessageStatusPending,
		ScheduledAt: time.Now().Add(1 * time.Hour),
	}
	if err := repo.Create(context.Background(), message); err != nil {
		t.Fatalf("Failed to create test message: %v", err)
	}

	// Update the status
	newStatus := model.MessageStatusSent
	if err := repo.UpdateStatus(context.Background(), message.ID, newStatus); err != nil {
		t.Errorf("UpdateStatus() error = %v", err)
	}

	// Verify the update
	found, err := repo.FindByID(context.Background(), message.ID)
	if err != nil {
		t.Errorf("FindByID() error = %v", err)
	}
	if found.Status != newStatus {
		t.Errorf("Expected status %q, got %q", newStatus, found.Status)
	}
}

func TestMessageRepository_FindPendingBefore(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMessageRepository(db)

	now := time.Now()

	// Create test messages
	messages := []*model.Message{
		{
			Content:     "Past message",
			To:          "test@example.com",
			Status:      model.MessageStatusPending,
			ScheduledAt: now.Add(-1 * time.Hour),
		},
		{
			Content:     "Future message",
			To:          "test@example.com",
			Status:      model.MessageStatusPending,
			ScheduledAt: now.Add(1 * time.Hour),
		},
		{
			Content:     "Sent message",
			To:          "test@example.com",
			Status:      model.MessageStatusSent,
			ScheduledAt: now.Add(-1 * time.Hour),
		},
	}

	for _, msg := range messages {
		if err := repo.Create(context.Background(), msg); err != nil {
			t.Fatalf("Failed to create test message: %v", err)
		}
	}

	// Test finding pending messages before now
	found, err := repo.FindPendingBefore(context.Background(), now, 10)
	if err != nil {
		t.Errorf("FindPendingBefore() error = %v", err)
	}
	if len(found) != 1 {
		t.Errorf("Expected 1 pending message, got %d", len(found))
	}
	if found[0].Content != "Past message" {
		t.Errorf("Expected message content %q, got %q", "Past message", found[0].Content)
	}
}
