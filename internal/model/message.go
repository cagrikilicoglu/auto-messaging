package model

import "time"

const (
	StatusPending = "pending"
	StatusSent    = "sent"
	StatusFailed  = "failed"
)

type Message struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	Content       string     `json:"content" gorm:"size:500"`
	To            string     `json:"to"`
	Status        string     `json:"status" gorm:"default:'pending'"`
	MessageID     string     `json:"message_id"`
	ScheduledTime time.Time  `json:"scheduled_time"`
	SentAt        *time.Time `json:"sent_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
