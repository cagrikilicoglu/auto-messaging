package controller

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"auto-messaging/internal/client"
	"auto-messaging/internal/model"
	"auto-messaging/internal/repository"
	"auto-messaging/pkg/cache"
)

var (
	ErrContentTooLong = errors.New("message content exceeds maximum length")
)

const (
	maxContentLength = 500
	batchSize        = 2
	processInterval  = 2 * time.Minute
)

// MessageController handles message operations
type MessageController struct {
	repo    repository.MessageRepository
	webhook client.WebhookClient
	cache   cache.MessageCache
	stopCh  chan struct{}
	logger  *log.Logger
}

// NewMessageController creates a new message controller instance
func NewMessageController(repo repository.MessageRepository, webhook client.WebhookClient, cache cache.MessageCache, logger *log.Logger) *MessageController {
	return &MessageController{
		repo:    repo,
		webhook: webhook,
		cache:   cache,
		stopCh:  make(chan struct{}),
		logger:  logger,
	}
}

// CreateMessage adds a new message to the system
func (c *MessageController) CreateMessage(content string, to string, scheduledTime time.Time) (*model.Message, error) {
	if len(content) > maxContentLength {
		return nil, ErrContentTooLong
	}

	msg := &model.Message{
		Content:       content,
		To:            to,
		ScheduledTime: scheduledTime,
		Status:        "pending",
	}

	if err := c.repo.Create(msg); err != nil {
		return nil, err
	}
	return msg, nil
}

// GetMessages fetches all messages
func (c *MessageController) GetMessages() ([]model.Message, error) {
	return c.repo.FindAll()
}

// GetMessageByID finds a message by its ID
func (c *MessageController) GetMessageByID(id uint) (*model.Message, error) {
	return c.repo.FindByID(id)
}

// UpdateMessageStatus changes the status of a message
func (c *MessageController) UpdateMessageStatus(id uint, status string) error {
	return c.repo.UpdateStatus(id, status)
}

// Start begins processing scheduled messages
func (c *MessageController) Start() error {
	go func() {
		// Process messages immediately when started
		c.processMessages()

		ticker := time.NewTicker(processInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				c.processMessages()
			case <-c.stopCh:
				return
			}
		}
	}()

	return nil
}

// processMessages handles the message processing logic
func (c *MessageController) processMessages() {
	// Get pending messages
	msgs, err := c.repo.FindPendingBefore(time.Now(), batchSize)
	if err != nil {
		c.logger.Printf("Failed to fetch pending messages: %v", err)
		return
	}

	for _, msg := range msgs {
		if err := c.processMessage(msg); err != nil {
			c.logger.Printf("Failed to process message %d: %v", msg.ID, err)
		}
	}
}

// processMessage handles the message processing logic for a single message
func (c *MessageController) processMessage(msg model.Message) error {
	// Prepare webhook request
	req := &model.WebhookRequest{
		To:      msg.To,
		Content: msg.Content,
	}

	// Send message
	resp, err := c.webhook.SendMessage(req)
	if err != nil {
		msg.Status = model.StatusFailed
		_ = c.repo.UpdateStatus(msg.ID, msg.Status)
		return fmt.Errorf("webhook send failed: %w", err)
	}

	// Update message status
	now := time.Now()
	msg.MessageID = resp.MessageID
	msg.Status = model.StatusSent
	msg.SentAt = &now

	// Cache message ID (non-critical operation)
	if err := c.cache.StoreMessageID(context.Background(), resp.MessageID, now); err != nil {
		c.logger.Printf("Warning: Failed to cache message ID %s: %v", resp.MessageID, err)
	}

	// Update database
	if err := c.repo.UpdateStatus(msg.ID, msg.Status); err != nil {
		return fmt.Errorf("failed to update message status: %w", err)
	}

	c.logger.Printf("Successfully processed message ID %d", msg.ID)
	return nil
}

// Stop halts message processing
func (c *MessageController) Stop() error {
	c.stopCh <- struct{}{}
	return nil
}

// GetSentMessages retrieves all sent messages
func (c *MessageController) GetSentMessages() ([]*model.Message, error) {
	return c.repo.FindByStatus("sent")
}
