package repository

import (
	"auto-messaging/config"
	"auto-messaging/internal/model"
	"context"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// MessageRepository defines the interface for message data access
type MessageRepository interface {
	Create(ctx context.Context, message *model.Message) error
	FindAll() ([]model.Message, error)
	FindByID(ctx context.Context, id uint) (*model.Message, error)
	UpdateStatus(ctx context.Context, id uint, status string) error
	FindPendingBefore(ctx context.Context, before time.Time, limit int) ([]*model.Message, error)
	FindByStatus(status string) ([]*model.Message, error)
	UpdateMessageID(ctx context.Context, id uint, messageID string) error
	UpdateSentAt(ctx context.Context, id uint, sentAt time.Time) error
}

// MessageRepositoryImpl implements the MessageRepository interface
type MessageRepositoryImpl struct {
	db *gorm.DB
}

// NewMessageRepository creates a new message repository
func NewMessageRepository(db *gorm.DB) *MessageRepositoryImpl {
	return &MessageRepositoryImpl{
		db: db,
	}
}

// InitDB initializes the database connection
func InitDB(cfg config.DB) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)

	log.Printf("Connecting to database with DSN: %s", dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate the Message model
	if err := db.AutoMigrate(&model.Message{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}

func (r *MessageRepositoryImpl) Create(ctx context.Context, message *model.Message) error {
	return r.db.WithContext(ctx).Create(message).Error
}

func (r *MessageRepositoryImpl) FindAll() ([]model.Message, error) {
	var messages []model.Message
	if err := r.db.Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *MessageRepositoryImpl) FindByID(ctx context.Context, id uint) (*model.Message, error) {
	var message model.Message
	err := r.db.WithContext(ctx).First(&message, id).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *MessageRepositoryImpl) UpdateStatus(ctx context.Context, id uint, status string) error {
	return r.db.WithContext(ctx).
		Model(&model.Message{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *MessageRepositoryImpl) FindPendingBefore(ctx context.Context, before time.Time, limit int) ([]*model.Message, error) {
	var messages []*model.Message
	err := r.db.WithContext(ctx).
		Where("status = ? AND scheduled_at <= ?", model.MessageStatusPending, before).
		Order("scheduled_at ASC").
		Limit(limit).
		Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *MessageRepositoryImpl) FindByStatus(status string) ([]*model.Message, error) {
	var messages []*model.Message
	if err := r.db.Where("status = ?", status).Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *MessageRepositoryImpl) UpdateMessageID(ctx context.Context, id uint, messageID string) error {
	return r.db.WithContext(ctx).
		Model(&model.Message{}).
		Where("id = ?", id).
		Update("message_id", messageID).Error
}

func (r *MessageRepositoryImpl) UpdateSentAt(ctx context.Context, id uint, sentAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&model.Message{}).
		Where("id = ?", id).
		Update("sent_at", sentAt).Error
}
