package handler

import (
	"net/http"
	"time"

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
// @Description Save a new conversation with messages and metadata
// @Tags conversations
// @Accept json
// @Produce json
// @Param request body models.ConversationSaveRequest true "Conversation save request"
// @Success 201 {object} models.APIResponse "Conversation saved successfully"
// @Failure 400 {object} models.APIResponse "Invalid request"
// @Failure 500 {object} models.APIResponse "Server error"
// @Router /api/rag/conversation/store [post]
func (sch *SaveConversationHandler) Handle(c *gin.Context) {
	startTime := time.Now()

	var req models.ConversationSaveRequest

	// Bind JSON request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.ErrorInfo{
				Code:    "INVALID_REQUEST",
				Message: "Invalid request body",
				Details: map[string]interface{}{
					"error": err.Error(),
				},
			},
			Metadata: models.Metadata{},
		})
		return
	}

	// Validate required fields
	if req.ConversationID == "" || len(req.Messages) == 0 {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.ErrorInfo{
				Code:    "INVALID_REQUEST",
				Message: "conversation_id and messages are required",
				Details: map[string]interface{}{
					"fields": []string{"conversation_id", "messages"},
					"reason": "required fields missing",
				},
			},
			Metadata: models.Metadata{},
		})
		return
	}

	// Validate messages have content
	for i, msg := range req.Messages {
		if msg.Content == "" {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Error: &models.ErrorInfo{
					Code:    "INVALID_REQUEST",
					Message: "message content cannot be empty",
					Details: map[string]interface{}{
						"message_index": i,
						"reason":        "message content is required",
					},
				},
				Metadata: models.Metadata{},
			})
			return
		}
		if msg.Role != "user" && msg.Role != "assistant" {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Error: &models.ErrorInfo{
					Code:    "INVALID_REQUEST",
					Message: "invalid message role",
					Details: map[string]interface{}{
						"message_index": i,
						"valid_roles":   []string{"user", "assistant"},
						"provided_role": msg.Role,
					},
				},
				Metadata: models.Metadata{},
			})
			return
		}
	}

	// Save conversation
	_, err := sch.conversationService.SaveConversation(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: &models.ErrorInfo{
				Code:    "INTERNAL_ERROR",
				Message: "failed to save conversation",
				Details: map[string]interface{}{
					"error": err.Error(),
				},
			},
			Metadata: models.Metadata{},
		})
		return
	}

	processingTimeMs := time.Since(startTime).Milliseconds()

	// Build save response
	saveResp := models.SaveResponse{
		ConversationID:   req.ConversationID,
		VectorsCreated:   1,
		MessagesStored:   len(req.Messages),
		StoredAt:         time.Now().UTC().Format(time.RFC3339),
		ProcessingTimeMs: processingTimeMs,
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success:  true,
		Data:     saveResp,
		Metadata: models.Metadata{},
	})
}
