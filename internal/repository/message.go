package repository

import (
	"time"

	"auto-messaging/internal/model"

	"gorm.io/gorm"
)

// MessageRepository defines the interface for message data access
type MessageRepository interface {
	Create(message *model.Message) error
	FindAll() ([]model.Message, error)
	FindByID(id uint) (*model.Message, error)
	UpdateStatus(id uint, status string) error
	FindPendingBefore(time time.Time) ([]model.Message, error)
	FindByStatus(status string) ([]*model.Message, error)
}

// messageRepository implements MessageRepository
type messageRepository struct {
	db *gorm.DB
}

// NewMessageRepository creates a new message repository
func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(message *model.Message) error {
	return r.db.Create(message).Error
}

func (r *messageRepository) FindAll() ([]model.Message, error) {
	var messages []model.Message
	if err := r.db.Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *messageRepository) FindByID(id uint) (*model.Message, error) {
	var message model.Message
	if err := r.db.First(&message, id).Error; err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *messageRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&model.Message{}).Where("id = ?", id).Update("status", status).Error
}

func (r *messageRepository) FindPendingBefore(time time.Time) ([]model.Message, error) {
	var messages []model.Message
	if err := r.db.Where("status = ? AND scheduled_time <= ?", "pending", time).Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *messageRepository) FindByStatus(status string) ([]*model.Message, error) {
	var messages []*model.Message
	if err := r.db.Where("status = ?", status).Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}
