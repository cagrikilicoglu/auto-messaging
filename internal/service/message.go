package service

import (
	"time"

	"gorm.io/gorm"
)

// Message represents a message in the system
type Message struct {
	gorm.Model
	Content       string    `json:"content"`
	ScheduledTime time.Time `json:"scheduled_time"`
	Status        string    `json:"status"` // pending, sent, failed
}

// MessageService handles message operations
type MessageService struct {
	db *gorm.DB
}

// NewMessageService creates a new message service instance
func NewMessageService(db *gorm.DB) *MessageService {
	return &MessageService{db: db}
}

// CreateMessage adds a new message to the system
func (s *MessageService) CreateMessage(content string, scheduledTime time.Time) (*Message, error) {
	msg := &Message{
		Content:       content,
		ScheduledTime: scheduledTime,
		Status:        "pending",
	}

	if err := s.db.Create(msg).Error; err != nil {
		return nil, err
	}
	return msg, nil
}

// GetMessages fetches all messages
func (s *MessageService) GetMessages() ([]Message, error) {
	var msgs []Message
	if err := s.db.Find(&msgs).Error; err != nil {
		return nil, err
	}
	return msgs, nil
}

// GetMessageByID finds a message by its ID
func (s *MessageService) GetMessageByID(id uint) (*Message, error) {
	var msg Message
	if err := s.db.First(&msg, id).Error; err != nil {
		return nil, err
	}
	return &msg, nil
}

// UpdateMessageStatus changes the status of a message
func (s *MessageService) UpdateMessageStatus(id uint, status string) error {
	return s.db.Model(&Message{}).Where("id = ?", id).Update("status", status).Error
}

// Start begins processing scheduled messages
func (s *MessageService) Start() error {
	go func() {
		ticker := time.NewTicker(2 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			var msgs []Message
			if err := s.db.Where("status = ? AND scheduled_time <= ?", "pending", time.Now()).Find(&msgs).Error; err != nil {
				continue
			}

			for _, msg := range msgs {
				// TODO: Implement actual message sending
				if err := s.UpdateMessageStatus(msg.ID, "sent"); err != nil {
					continue
				}
			}
		}
	}()
	return nil
}

// Stop halts message processing
func (s *MessageService) Stop() error {
	// TODO: Implement proper shutdown
	return nil
}

// GetSentMessages retrieves all sent messages
func (s *MessageService) GetSentMessages() ([]*Message, error) {
	var msgs []*Message
	if err := s.db.Where("status = ?", "sent").Find(&msgs).Error; err != nil {
		return nil, err
	}
	return msgs, nil
}
