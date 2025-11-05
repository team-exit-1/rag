package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

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
// @Param top_k query int false "Result limit (default: 10, max: 100)"
// @Success 200 {object} models.APIResponse "Search results with metadata"
// @Failure 400 {object} models.APIResponse "Invalid request"
// @Failure 500 {object} models.APIResponse "Server error"
// @Router /api/rag/conversation/search [get]
func (sch *SearchConversationHandler) Handle(c *gin.Context) {
	startTime := time.Now()

	// Get query parameters
	query := c.Query("query")
	topKStr := c.DefaultQuery("top_k", "5")

	// Parse top_k
	topK := 5
	if topKStr != "" {
		if k, err := strconv.Atoi(topKStr); err == nil && k > 0 && k <= 100 {
			topK = k
		}
	}

	// Validate required fields
	if query == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.ErrorInfo{
				Code:    "INVALID_REQUEST",
				Message: "Query text cannot be empty",
				Details: map[string]interface{}{
					"field":  "query",
					"reason": "required field missing",
				},
			},
			Metadata: models.Metadata{},
		})
		return
	}

	query = strings.TrimSpace(query)

	req := models.ConversationSearchRequest{
		Query: query,
		Limit: topK,
	}

	// Search conversations
	results, err := sch.conversationService.SearchConversations(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: &models.ErrorInfo{
				Code:    "INTERNAL_ERROR",
				Message: "failed to search conversations",
				Details: map[string]string{
					"error": err.Error(),
				},
			},
			Metadata: models.Metadata{},
		})
		return
	}

	searchTimeMs := time.Since(startTime).Milliseconds()

	// Build search response
	searchResp := models.SearchResponse{
		Query:        query,
		Results:      results,
		TotalResults: len(results),
		SearchMetadata: models.SearchMetadata{
			EmbeddingModel: "text-embedding-3-large",
			VectorDB:       "qdrant",
			SearchTimeMs:   searchTimeMs,
		},
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success:  true,
		Data:     searchResp,
		Metadata: models.Metadata{},
	})
}
