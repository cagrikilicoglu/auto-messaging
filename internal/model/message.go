package model

import (
	"time"

	"gorm.io/gorm"
)

// Message status constants
const (
	MessageStatusPending   = "pending"
	MessageStatusSent      = "sent"
	MessageStatusFailed    = "failed"
	MessageStatusCancelled = "cancelled"
)

// Message represents a message in the system
type Message struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Content     string    `json:"content"`
	To          string    `json:"to"`
	Status      string    `json:"status"`
	MessageID   string    `json:"message_id"`
	SentAt      time.Time `json:"sent_at"`
	ScheduledAt time.Time `json:"scheduled_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (m *Message) BeforeCreate(tx *gorm.DB) error {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return nil
}

func (m *Message) BeforeUpdate(tx *gorm.DB) error {
	m.UpdatedAt = time.Now()
	return nil
}
