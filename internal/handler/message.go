package handler

import (
	"net/http"
	"strconv"
	"time"

	"auto-messaging/internal/service"

	"github.com/gin-gonic/gin"
)

// MessageHandler handles HTTP requests for messages
type MessageHandler struct {
	messageService *service.MessageService
}

// NewMessageHandler creates a new message handler
func NewMessageHandler(messageService *service.MessageService) *MessageHandler {
	return &MessageHandler{messageService: messageService}
}

// CreateMessageRequest represents the request body for creating a message
type CreateMessageRequest struct {
	Content       string    `json:"content" binding:"required"`
	ScheduledTime time.Time `json:"scheduled_time" binding:"required"`
}

// CreateMessage handles the creation of a new message
func (h *MessageHandler) CreateMessage(c *gin.Context) {
	var req CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	msg, err := h.messageService.CreateMessage(req.Content, req.ScheduledTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create message"})
		return
	}

	c.JSON(http.StatusCreated, msg)
}

// GetMessages handles retrieving all messages
func (h *MessageHandler) GetMessages(c *gin.Context) {
	msgs, err := h.messageService.GetMessages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	c.JSON(http.StatusOK, msgs)
}

// GetMessageByID handles retrieving a message by its ID
func (h *MessageHandler) GetMessageByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	msg, err := h.messageService.GetMessageByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	c.JSON(http.StatusOK, msg)
}

// UpdateMessageStatus handles updating the status of a message
func (h *MessageHandler) UpdateMessageStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status format"})
		return
	}

	if err := h.messageService.UpdateMessageStatus(uint(id), req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status updated"})
}

// @Summary Start automatic message sending
// @Description Start the automatic message sending process
// @Tags messages
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/v1/messaging/start [post]
func (h *MessageHandler) StartMessaging(c *gin.Context) {
	if err := h.messageService.Start(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start messaging"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Messaging started"})
}

// @Summary Stop automatic message sending
// @Description Stop the automatic message sending process
// @Tags messages
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/v1/messaging/stop [post]
func (h *MessageHandler) StopMessaging(c *gin.Context) {
	if err := h.messageService.Stop(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stop messaging"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Messaging stopped"})
}

// @Summary Get sent messages
// @Description Get a list of all sent messages
// @Tags messages
// @Accept json
// @Produce json
// @Success 200 {array} model.Message
// @Router /api/v1/messaging/sent [get]
func (h *MessageHandler) GetSentMessages(c *gin.Context) {
	msgs, err := h.messageService.GetSentMessages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sent messages"})
		return
	}
	c.JSON(http.StatusOK, msgs)
}
