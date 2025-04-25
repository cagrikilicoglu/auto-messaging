package controller

import (
	"errors"
	"log"
	"time"

	"auto-messaging/internal/client"
	"auto-messaging/internal/model"
	"auto-messaging/internal/repository"
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
	webhook *client.WebhookClient
	stopCh  chan struct{}
}

// NewMessageController creates a new message controller instance
func NewMessageController(repo repository.MessageRepository, webhook *client.WebhookClient) *MessageController {
	return &MessageController{
		repo:    repo,
		webhook: webhook,
		stopCh:  make(chan struct{}),
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
	msgs, err := c.repo.FindPendingBefore(time.Now())
	if err != nil {
		log.Printf("Error finding pending messages: %v", err)
		return
	}

	log.Printf("Found %d pending messages", len(msgs))

	// Process only up to batchSize messages
	if len(msgs) > batchSize {
		msgs = msgs[:batchSize]
	}

	for _, msg := range msgs {
		req := &model.WebhookRequest{
			To:      msg.To,
			Content: msg.Content,
		}

		log.Printf("Sending message ID %d to webhook", msg.ID)
		resp, err := c.webhook.SendMessage(req)
		if err != nil {
			log.Printf("Error sending message ID %d: %v", msg.ID, err)
			_ = c.repo.UpdateStatus(msg.ID, "failed")
			continue
		}
		log.Printf("Webhook response for message ID %d: %+v", msg.ID, resp)

		now := time.Now()
		msg.MessageID = resp.MessageID
		msg.Status = "sent"
		msg.SentAt = &now

		if err := c.repo.UpdateStatus(msg.ID, msg.Status); err != nil {
			log.Printf("Error updating message ID %d status: %v", msg.ID, err)
			continue
		}
		log.Printf("Successfully processed message ID %d", msg.ID)
	}
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
