package router

import (
	"auto-messaging/internal/handler"

	"github.com/gin-gonic/gin"
)

// SetupRouter initializes the API routes
func SetupRouter(messageHandler *handler.MessageHandler) *gin.Engine {
	r := gin.Default()

	// API routes
	api := r.Group("/api/v1")
	{
		// Message management
		msgs := api.Group("/messages")
		{
			msgs.POST("", messageHandler.CreateMessage)
			msgs.GET("", messageHandler.GetMessages)
			msgs.GET("/:id", messageHandler.GetMessageByID)
			msgs.PUT("/:id/status", messageHandler.UpdateMessageStatus)
		}

		// Message processing control
		ctrl := api.Group("/messaging")
		{
			ctrl.POST("/start", messageHandler.StartMessaging)
			ctrl.POST("/stop", messageHandler.StopMessaging)
			ctrl.GET("/sent", messageHandler.GetSentMessages)
		}
	}

	return r
}
