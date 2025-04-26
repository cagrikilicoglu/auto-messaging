package controller

import (
	"context"
	"errors"
	"log"
	"os"
	"testing"
	"time"

	"auto-messaging/internal/model"
)

// MockWebhookClient implements the WebhookClient interface for testing
type mockWebhookClient struct {
	sendMessageFunc func(req *model.WebhookRequest) (*model.WebhookResponse, error)
}

func (m *mockWebhookClient) SendMessage(req *model.WebhookRequest) (*model.WebhookResponse, error) {
	return m.sendMessageFunc(req)
}

// MockMessageRepository implements the MessageRepository interface for testing
type mockMessageRepository struct {
	createFunc            func(ctx context.Context, message *model.Message) error
	findAllFunc           func() ([]model.Message, error)
	findByIDFunc          func(ctx context.Context, id uint) (*model.Message, error)
	updateStatusFunc      func(ctx context.Context, id uint, status string) error
	findPendingBeforeFunc func(ctx context.Context, before time.Time, limit int) ([]*model.Message, error)
	findByStatusFunc      func(status string) ([]*model.Message, error)
	updateMessageIDFunc   func(ctx context.Context, id uint, messageID string) error
	updateSentAtFunc      func(ctx context.Context, id uint, sentAt time.Time) error
	messages              map[uint]*model.Message
}

func (m *mockMessageRepository) Create(ctx context.Context, message *model.Message) error {
	return m.createFunc(ctx, message)
}

func (m *mockMessageRepository) FindAll() ([]model.Message, error) {
	return m.findAllFunc()
}

func (m *mockMessageRepository) FindByID(ctx context.Context, id uint) (*model.Message, error) {
	return m.findByIDFunc(ctx, id)
}

func (m *mockMessageRepository) UpdateStatus(ctx context.Context, id uint, status string) error {
	return m.updateStatusFunc(ctx, id, status)
}

func (m *mockMessageRepository) FindPendingBefore(ctx context.Context, before time.Time, limit int) ([]*model.Message, error) {
	if m.findPendingBeforeFunc != nil {
		return m.findPendingBeforeFunc(ctx, before, limit)
	}
	return []*model.Message{}, nil
}

func (m *mockMessageRepository) FindByStatus(status string) ([]*model.Message, error) {
	return m.findByStatusFunc(status)
}

func (m *mockMessageRepository) UpdateMessageID(ctx context.Context, id uint, messageID string) error {
	return m.updateMessageIDFunc(ctx, id, messageID)
}

func (m *mockMessageRepository) UpdateSentAt(ctx context.Context, id uint, sentAt time.Time) error {
	return m.updateSentAtFunc(ctx, id, sentAt)
}

// MockMessageCache implements the MessageCache interface for testing
type mockMessageCache struct {
	storeMessageIDFunc     func(ctx context.Context, messageID string, sentAt time.Time) error
	getMessageSentTimeFunc func(ctx context.Context, messageID string) (*time.Time, error)
}

func (m *mockMessageCache) StoreMessageID(ctx context.Context, messageID string, sentAt time.Time) error {
	return m.storeMessageIDFunc(ctx, messageID, sentAt)
}

func (m *mockMessageCache) GetMessageSentTime(ctx context.Context, messageID string) (*time.Time, error) {
	return m.getMessageSentTimeFunc(ctx, messageID)
}

func TestMessageController_CreateMessage(t *testing.T) {
	tests := []struct {
		name          string
		message       *model.Message
		createErr     error
		expectedError bool
	}{
		{
			name: "successful message creation",
			message: &model.Message{
				Content:     "Test message",
				To:          "test@example.com",
				ScheduledAt: time.Now().Add(1 * time.Hour),
				Status:      model.MessageStatusPending,
			},
			createErr:     nil,
			expectedError: false,
		},
		{
			name: "message creation fails",
			message: &model.Message{
				Content:     "Test message",
				To:          "test@example.com",
				ScheduledAt: time.Now().Add(1 * time.Hour),
				Status:      model.MessageStatusPending,
			},
			createErr:     errors.New("database error"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockMessageRepository{
				createFunc: func(ctx context.Context, message *model.Message) error {
					return tt.createErr
				},
			}

			controller := NewMessageController(
				repo,
				&mockWebhookClient{},
				&mockMessageCache{},
				nil,
			)

			err := controller.repo.Create(context.Background(), tt.message)
			if (err != nil) != tt.expectedError {
				t.Errorf("CreateMessage() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

func TestMessageController_ProcessMessage(t *testing.T) {
	tests := []struct {
		name          string
		message       *model.Message
		webhookResp   *model.WebhookResponse
		webhookErr    error
		expectedError bool
	}{
		{
			name: "successful message processing",
			message: &model.Message{
				ID:          1,
				Content:     "Test message",
				To:          "test@example.com",
				Status:      model.MessageStatusPending,
				ScheduledAt: time.Now().Add(-1 * time.Hour), // Past time
			},
			webhookResp: &model.WebhookResponse{
				MessageID: "test-message-id",
			},
			webhookErr:    nil,
			expectedError: false,
		},
		{
			name: "webhook error",
			message: &model.Message{
				ID:          1,
				Content:     "Test message",
				To:          "test@example.com",
				Status:      model.MessageStatusPending,
				ScheduledAt: time.Now().Add(-1 * time.Hour),
			},
			webhookResp:   nil,
			webhookErr:    errors.New("webhook error"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			webhookClient := &mockWebhookClient{
				sendMessageFunc: func(req *model.WebhookRequest) (*model.WebhookResponse, error) {
					return tt.webhookResp, tt.webhookErr
				},
			}

			repo := &mockMessageRepository{
				updateStatusFunc: func(ctx context.Context, id uint, status string) error {
					return nil
				},
				updateMessageIDFunc: func(ctx context.Context, id uint, messageID string) error {
					return nil
				},
				updateSentAtFunc: func(ctx context.Context, id uint, sentAt time.Time) error {
					return nil
				},
			}

			controller := NewMessageController(
				repo,
				webhookClient,
				&mockMessageCache{},
				nil,
			)

			err := controller.processMessage(tt.message)
			if (err != nil) != tt.expectedError {
				t.Errorf("ProcessMessage() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

func TestMessageController_StartStop(t *testing.T) {
	// Create mock repository
	mockRepo := &mockMessageRepository{
		messages: make(map[uint]*model.Message),
		findPendingBeforeFunc: func(ctx context.Context, before time.Time, limit int) ([]*model.Message, error) {
			return []*model.Message{}, nil
		},
	}

	// Create mock webhook client
	mockWebhook := &mockWebhookClient{}

	// Create mock cache
	mockCache := &mockMessageCache{}

	// Create logger
	logger := log.New(os.Stdout, "", log.LstdFlags)

	// Create controller
	controller := NewMessageController(mockRepo, mockWebhook, mockCache, logger)

	// Test Start
	err := controller.Start()
	if err != nil {
		t.Errorf("Start() error = %v", err)
	}

	// Wait a bit to ensure the goroutine has started
	time.Sleep(100 * time.Millisecond)

	// Test Stop
	err = controller.Stop()
	if err != nil {
		t.Errorf("Stop() error = %v", err)
	}
}
