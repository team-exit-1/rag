package handler

import (
	"net/http"
	"strconv"

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
// @Summary Search conversations
// @Description Search for conversations by semantic similarity
// @Tags conversations
// @Produce json
// @Param query query string true "Search query"
// @Param user_id query string false "User ID"
// @Param limit query int false "Result limit (default: 10, max: 100)"
// @Success 200 {object} map[string]interface{} "Search results with count"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Server error"
// @Router /api/v1/conversations/search [get]
func (sch *SearchConversationHandler) Handle(c *gin.Context) {
	// Get query parameters
	query := c.Query("query")
	userID := c.Query("user_id")
	limitStr := c.DefaultQuery("limit", "10")

	// Parse limit
	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	req := models.ConversationSearchRequest{
		Query:  query,
		UserID: userID,
		Limit:  limit,
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
