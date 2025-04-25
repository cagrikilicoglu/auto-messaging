package handler

import (
	"net/http"
	"time"

	"auto-messaging/internal/controller"

	"github.com/gin-gonic/gin"
)

// MessageHandler handles HTTP requests for messages
type MessageHandler struct {
	controller *controller.MessageController
}

// NewMessageHandler creates a new message handler
func NewMessageHandler(controller *controller.MessageController) *MessageHandler {
	return &MessageHandler{controller: controller}
}

// CreateMessageRequest represents the request body for creating a message
type CreateMessageRequest struct {
	Content       string    `json:"content" binding:"required"`
	To            string    `json:"to" binding:"required"`
	ScheduledTime time.Time `json:"scheduled_time" binding:"required"`
}

// CreateMessage handles the creation of a new message
func (h *MessageHandler) CreateMessage(c *gin.Context) {
	h.controller.CreateMessage(c)
}

// GetMessages handles retrieving all messages
func (h *MessageHandler) GetMessages(c *gin.Context) {
	h.controller.GetMessages(c)
}

// GetMessageByID handles retrieving a message by its ID
func (h *MessageHandler) GetMessageByID(c *gin.Context) {
	h.controller.GetMessage(c)
}

// UpdateMessageStatus handles updating the status of a message
func (h *MessageHandler) UpdateMessageStatus(c *gin.Context) {
	h.controller.UpdateMessage(c)
}

// @Summary Start automatic message sending
// @Description Start the automatic message sending process
// @Tags messages
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/v1/messaging/start [post]
func (h *MessageHandler) StartMessaging(c *gin.Context) {
	h.controller.StartMessaging(c)
}

// @Summary Stop automatic message sending
// @Description Stop the automatic message sending process
// @Tags messages
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/v1/messaging/stop [post]
func (h *MessageHandler) StopMessaging(c *gin.Context) {
	h.controller.StopMessaging(c)
}

// @Summary Get sent messages
// @Description Get a list of all sent messages
// @Tags messages
// @Accept json
// @Produce json
// @Success 200 {array} model.Message
// @Router /api/v1/messaging/sent [get]
func (h *MessageHandler) GetSentMessages(c *gin.Context) {
	msgs, err := h.controller.GetSentMessages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sent messages"})
		return
	}

	c.JSON(http.StatusOK, msgs)
}
