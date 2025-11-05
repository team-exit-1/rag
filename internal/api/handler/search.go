package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"repo-rag-server/internal/models"
	"repo-rag-server/internal/service"
)

// SearchConversationHandler handles conversation search requests
type SearchConversationHandler struct {
	conversationService *service.ConversationService
}

// NewSearchConversationHandler creates a new search conversation handler
func NewSearchConversationHandler(conversationService *service.ConversationService) *SearchConversationHandler {
	return &SearchConversationHandler{
		conversationService: conversationService,
	}
}

// Handle processes search conversation requests
func (sch *SearchConversationHandler) Handle(c *gin.Context) {
	var req models.ConversationSearchRequest

	// Bind JSON request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Validate required fields
	if req.Query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "query is required",
		})
		return
	}

	// Search conversations
	results, err := sch.conversationService.SearchConversations(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to search conversations",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"results": results,
		"count":   len(results),
	})
}
