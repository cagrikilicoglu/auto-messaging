package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type MessageCache interface {
	StoreMessageID(ctx context.Context, messageID string, sentAt time.Time) error
	GetMessageSentTime(ctx context.Context, messageID string) (*time.Time, error)
}

type redisCache struct {
	client *redis.Client
}

func NewRedisCache(host string, port int, password string, db int) MessageCache {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
	})

	return &redisCache{
		client: client,
	}
}

func (c *redisCache) StoreMessageID(ctx context.Context, messageID string, sentAt time.Time) error {
	key := fmt.Sprintf("message:%s", messageID)
	return c.client.Set(ctx, key, sentAt.Unix(), 24*time.Hour).Err()
}

func (c *redisCache) GetMessageSentTime(ctx context.Context, messageID string) (*time.Time, error) {
	key := fmt.Sprintf("message:%s", messageID)
	val, err := c.client.Get(ctx, key).Int64()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	t := time.Unix(val, 0)
	return &t, nil
}
