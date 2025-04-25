package controller

import (
	"time"

	"auto-messaging/internal/model"
	"auto-messaging/internal/repository"
)

// MessageController handles message operations
type MessageController struct {
	repo repository.MessageRepository
}

// NewMessageController creates a new message controller instance
func NewMessageController(repo repository.MessageRepository) *MessageController {
	return &MessageController{repo: repo}
}

// CreateMessage adds a new message to the system
func (c *MessageController) CreateMessage(content string, scheduledTime time.Time) (*model.Message, error) {
	msg := &model.Message{
		Content:       content,
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
		ticker := time.NewTicker(2 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			msgs, err := c.repo.FindPendingBefore(time.Now())
			if err != nil {
				continue
			}

			for _, msg := range msgs {
				// TODO: Implement actual message sending
				if err := c.repo.UpdateStatus(msg.ID, "sent"); err != nil {
					continue
				}
			}
		}
	}()

	return nil
}

// Stop halts message processing
func (c *MessageController) Stop() error {
	// TODO: Implement proper shutdown
	return nil
}

// GetSentMessages retrieves all sent messages
func (c *MessageController) GetSentMessages() ([]*model.Message, error) {
	return c.repo.FindByStatus("sent")
}
