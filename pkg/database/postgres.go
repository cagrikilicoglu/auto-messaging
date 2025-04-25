package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB(host string, port int, user, password, dbname string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

type Message struct {
	ID        uint   `gorm:"primaryKey"`
	Content   string `gorm:"size:500"`
	To        string
	Status    string `gorm:"default:'pending'"`
	MessageID string
	SentAt    *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
