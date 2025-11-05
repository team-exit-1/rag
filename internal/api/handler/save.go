package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"repo-rag-server/internal/models"
	"repo-rag-server/internal/service"
)

// SaveConversationHandler handles conversation save requests
type SaveConversationHandler struct {
	conversationService *service.ConversationService
}

// NewSaveConversationHandler creates a new save conversation handler
func NewSaveConversationHandler(conversationService *service.ConversationService) *SaveConversationHandler {
	return &SaveConversationHandler{
		conversationService: conversationService,
	}
}

// Handle processes save conversation requests
// @Summary Save a conversation
// @Description Save a new conversation with question and answer
// @Tags conversations
// @Accept json
// @Produce json
// @Param request body models.ConversationSaveRequest true "Conversation save request"
// @Success 201 {object} models.ConversationResponse
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Server error"
// @Router /api/v1/conversations [post]
func (sch *SaveConversationHandler) Handle(c *gin.Context) {
	var req models.ConversationSaveRequest

	// Bind JSON request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Validate required fields
	if req.UserID == "" || req.Question == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_id and question are required",
		})
		return
	}

	// Save conversation
	response, err := sch.conversationService.SaveConversation(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to save conversation",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response)
}
