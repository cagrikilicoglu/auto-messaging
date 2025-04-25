package controller

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"auto-messaging/internal/client"
	"auto-messaging/internal/model"
	"auto-messaging/internal/repository"
	"auto-messaging/pkg/cache"

	"github.com/gin-gonic/gin"
)

var (
	ErrContentTooLong = errors.New("message content exceeds maximum length")
)

const (
	maxContentLength = 500
	batchSize        = 2
	processInterval  = 2 * time.Minute
)

// MessageController handles HTTP requests for messages
type MessageController struct {
	repo    repository.MessageRepository
	webhook client.WebhookClient
	cache   cache.MessageCache
	stopCh  chan struct{}
	logger  *log.Logger
}

// NewMessageController creates a new MessageController
func NewMessageController(repo repository.MessageRepository, webhook client.WebhookClient, cache cache.MessageCache, logger *log.Logger) *MessageController {
	return &MessageController{
		repo:    repo,
		webhook: webhook,
		cache:   cache,
		stopCh:  make(chan struct{}),
		logger:  logger,
	}
}

// CreateMessageRequest represents the request body for creating a message
type CreateMessageRequest struct {
	Content     string    `json:"content" binding:"required"`
	To          string    `json:"to" binding:"required,email"`
	ScheduledAt time.Time `json:"scheduled_at" binding:"required"`
}

// @Summary Create a new message
// @Description Create a new message with the provided details
// @Tags messages
// @Accept json
// @Produce json
// @Param message body CreateMessageRequest true "Message details"
// @Success 201 {object} model.Message
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /messages [post]
func (c *MessageController) CreateMessage(ctx *gin.Context) {
	var req CreateMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if len(req.Content) > maxContentLength {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: ErrContentTooLong.Error()})
		return
	}

	message := &model.Message{
		Content:     req.Content,
		To:          req.To,
		ScheduledAt: req.ScheduledAt,
		Status:      model.MessageStatusPending,
	}

	if err := c.repo.Create(context.Background(), message); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create message"})
		return
	}

	ctx.JSON(http.StatusCreated, message)
}

// @Summary Get all messages
// @Description Get a list of all messages
// @Tags messages
// @Produce json
// @Success 200 {array} model.Message
// @Failure 500 {object} ErrorResponse
// @Router /messages [get]
func (c *MessageController) GetMessages(ctx *gin.Context) {
	messages, err := c.repo.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get messages"})
		return
	}

	ctx.JSON(http.StatusOK, messages)
}

// @Summary Get a message by ID
// @Description Get a message by its ID
// @Tags messages
// @Produce json
// @Param id path int true "Message ID"
// @Success 200 {object} model.Message
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /messages/{id} [get]
func (c *MessageController) GetMessage(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid message ID"})
		return
	}

	message, err := c.repo.FindByID(context.Background(), uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, ErrorResponse{Error: "Message not found"})
		return
	}

	ctx.JSON(http.StatusOK, message)
}

// @Summary Update a message
// @Description Update an existing message
// @Tags messages
// @Accept json
// @Produce json
// @Param id path int true "Message ID"
// @Param message body CreateMessageRequest true "Updated message details"
// @Success 200 {object} model.Message
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /messages/{id} [put]
func (c *MessageController) UpdateMessage(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid message ID"})
		return
	}

	var req CreateMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	message, err := c.repo.FindByID(context.Background(), uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, ErrorResponse{Error: "Message not found"})
		return
	}

	message.Content = req.Content
	message.To = req.To
	message.ScheduledAt = req.ScheduledAt

	if err := c.repo.UpdateStatus(context.Background(), uint(id), message.Status); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to update message"})
		return
	}

	ctx.JSON(http.StatusOK, message)
}

// @Summary Delete a message
// @Description Delete a message by its ID
// @Tags messages
// @Produce json
// @Param id path int true "Message ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /messages/{id} [delete]
func (c *MessageController) DeleteMessage(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid message ID"})
		return
	}

	if err := c.repo.UpdateStatus(context.Background(), uint(id), model.MessageStatusCancelled); err != nil {
		ctx.JSON(http.StatusNotFound, ErrorResponse{Error: "Message not found"})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// @Summary Start message processing
// @Description Start processing pending messages
// @Tags messaging
// @Produce json
// @Success 200 {object} MessageResponse
// @Failure 500 {object} ErrorResponse
// @Router /messaging/start [post]
func (c *MessageController) StartMessaging(ctx *gin.Context) {
	if err := c.Start(); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to start messaging"})
		return
	}
	ctx.JSON(http.StatusOK, MessageResponse{Message: "Messaging started"})
}

// @Summary Stop message processing
// @Description Stop processing messages
// @Tags messaging
// @Produce json
// @Success 200 {object} MessageResponse
// @Failure 500 {object} ErrorResponse
// @Router /messaging/stop [post]
func (c *MessageController) StopMessaging(ctx *gin.Context) {
	if err := c.Stop(); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to stop messaging"})
		return
	}
	ctx.JSON(http.StatusOK, MessageResponse{Message: "Messaging stopped"})
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
	messages, err := c.repo.FindPendingBefore(context.Background(), time.Now(), batchSize)
	if err != nil {
		c.logger.Printf("Error finding pending messages: %v", err)
		return
	}

	for _, msg := range messages {
		if err := c.processMessage(msg); err != nil {
			c.logger.Printf("Failed to process message %d: %v", msg.ID, err)
		}
	}
}

// processMessage handles the message processing logic for a single message
func (c *MessageController) processMessage(msg *model.Message) error {
	// Check if message is already processed
	if msg.Status != model.MessageStatusPending {
		return nil
	}

	// Send message via webhook
	req := &model.WebhookRequest{
		Content: msg.Content,
		To:      msg.To,
	}
	resp, err := c.webhook.SendMessage(req)
	if err != nil {
		c.logger.Printf("Failed to send message %d: %v", msg.ID, err)
		return err
	}

	// Update message ID
	if err := c.repo.UpdateMessageID(context.Background(), msg.ID, resp.MessageID); err != nil {
		c.logger.Printf("Failed to update message %d ID: %v", msg.ID, err)
		return err
	}

	// Update message status
	now := time.Now()
	if err := c.repo.UpdateStatus(context.Background(), msg.ID, model.MessageStatusSent); err != nil {
		c.logger.Printf("Failed to update message %d status: %v", msg.ID, err)
		return err
	}

	// Update sent time
	if err := c.repo.UpdateSentAt(context.Background(), msg.ID, now); err != nil {
		c.logger.Printf("Failed to update message %d sent time: %v", msg.ID, err)
		return err
	}

	return nil
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
