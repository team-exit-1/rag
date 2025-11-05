package api

import (
	"github.com/gin-gonic/gin"

	"repo-rag-server/internal/api/handler"
	"repo-rag-server/internal/service"
)

// Router configures all API routes
func Router(conversationService *service.ConversationService) *gin.Engine {
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", handler.HealthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Save conversation endpoint
		saveHandler := handler.NewSaveConversationHandler(conversationService)
		v1.POST("/conversations", saveHandler.Handle)

		// Search conversations endpoint
		searchHandler := handler.NewSearchConversationHandler(conversationService)
		v1.POST("/conversations/search", searchHandler.Handle)
	}

	return router
}