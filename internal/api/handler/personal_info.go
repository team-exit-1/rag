package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"refo-rag-server/internal/models"
	"refo-rag-server/internal/storage"
)

// PersonalInfoHandler handles personal information requests from guardians
type PersonalInfoHandler struct {
	personalInfoStore storage.PersonalInfoStore
}

// NewPersonalInfoHandler creates a new personal info handler
func NewPersonalInfoHandler(personalInfoStore storage.PersonalInfoStore) *PersonalInfoHandler {
	return &PersonalInfoHandler{
		personalInfoStore: personalInfoStore,
	}
}

// CreatePersonalInfo creates a new personal information entry
// @Summary Create personal information
// @Description Save personal information provided by guardians (medical, contact, emergency, etc.)
// @Tags personal-info
// @Accept json
// @Produce json
// @Param request body models.PersonalInfoCreateRequest true "Personal info creation request"
// @Success 201 {object} models.APIResponse "Personal info created successfully"
// @Failure 400 {object} models.APIResponse "Invalid request"
// @Failure 500 {object} models.APIResponse "Server error"
// @Router /api/rag/personal-info [post]
func (pih *PersonalInfoHandler) CreatePersonalInfo(c *gin.Context) {
	startTime := time.Now()

	var req models.PersonalInfoCreateRequest

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

	// Create new personal info entry
	now := time.Now()
	personalInfo := &models.PersonalInfo{
		ID:         uuid.New().String(),
		UserID:     req.UserID,
		Content:    req.Content,
		Category:   req.Category,
		Importance: req.Importance,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// Save personal info
	if err := pih.personalInfoStore.SavePersonalInfo(context.Background(), personalInfo); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: &models.ErrorInfo{
				Code:    "INTERNAL_ERROR",
				Message: "failed to create personal information",
				Details: map[string]interface{}{
					"error": err.Error(),
				},
			},
			Metadata: models.Metadata{},
		})
		return
	}

	processingTimeMs := time.Since(startTime).Milliseconds()

	// Build response
	infoResp := models.PersonalInfoResponse{
		ID:         personalInfo.ID,
		UserID:     personalInfo.UserID,
		Content:    personalInfo.Content,
		Category:   personalInfo.Category,
		Importance: personalInfo.Importance,
		CreatedAt:  personalInfo.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:  personalInfo.UpdatedAt.UTC().Format(time.RFC3339),
	}

	response := map[string]interface{}{
		"personal_info":        infoResp,
		"processing_time_ms":   processingTimeMs,
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success:  true,
		Data:     response,
		Metadata: models.Metadata{},
	})
}

// GetPersonalInfo retrieves a personal information entry by ID
// @Summary Get personal information
// @Description Retrieve a personal information entry by ID
// @Tags personal-info
// @Produce json
// @Param info_id path string true "Personal info ID"
// @Success 200 {object} models.APIResponse "Personal info retrieved successfully"
// @Failure 404 {object} models.APIResponse "Personal info not found"
// @Failure 500 {object} models.APIResponse "Server error"
// @Router /api/rag/personal-info/{info_id} [get]
func (pih *PersonalInfoHandler) GetPersonalInfo(c *gin.Context) {
	startTime := time.Now()
	infoID := c.Param("info_id")

	if infoID == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.ErrorInfo{
				Code:    "INVALID_REQUEST",
				Message: "info_id is required",
				Details: map[string]interface{}{
					"field":  "info_id",
					"reason": "required parameter missing",
				},
			},
			Metadata: models.Metadata{},
		})
		return
	}

	// Get personal info
	personalInfo, err := pih.personalInfoStore.GetPersonalInfo(context.Background(), infoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: &models.ErrorInfo{
				Code:    "INTERNAL_ERROR",
				Message: "failed to get personal information",
				Details: map[string]interface{}{
					"error": err.Error(),
				},
			},
			Metadata: models.Metadata{},
		})
		return
	}

	if personalInfo == nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error: &models.ErrorInfo{
				Code:    "PERSONAL_INFO_NOT_FOUND",
				Message: "personal information not found",
				Details: map[string]interface{}{
					"info_id": infoID,
				},
			},
			Metadata: models.Metadata{},
		})
		return
	}

	processingTimeMs := time.Since(startTime).Milliseconds()

	// Build response
	infoResp := models.PersonalInfoResponse{
		ID:         personalInfo.ID,
		UserID:     personalInfo.UserID,
		Content:    personalInfo.Content,
		Category:   personalInfo.Category,
		Importance: personalInfo.Importance,
		CreatedAt:  personalInfo.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:  personalInfo.UpdatedAt.UTC().Format(time.RFC3339),
	}

	response := map[string]interface{}{
		"personal_info":        infoResp,
		"processing_time_ms":   processingTimeMs,
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success:  true,
		Data:     response,
		Metadata: models.Metadata{},
	})
}

// GetPersonalInfoByUser retrieves all personal information for a user
// @Summary Get all personal information for a user
// @Description Retrieve all personal information entries for a specific user
// @Tags personal-info
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} models.APIResponse "Personal info list retrieved successfully"
// @Failure 400 {object} models.APIResponse "Invalid request"
// @Failure 500 {object} models.APIResponse "Server error"
// @Router /api/rag/personal-info/user/{user_id} [get]
func (pih *PersonalInfoHandler) GetPersonalInfoByUser(c *gin.Context) {
	startTime := time.Now()
	userID := c.Param("user_id")

	if userID == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.ErrorInfo{
				Code:    "INVALID_REQUEST",
				Message: "user_id is required",
				Details: map[string]interface{}{
					"field":  "user_id",
					"reason": "required parameter missing",
				},
			},
			Metadata: models.Metadata{},
		})
		return
	}

	// Get all personal info for user
	personalInfoList, err := pih.personalInfoStore.GetPersonalInfoByUser(context.Background(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: &models.ErrorInfo{
				Code:    "INTERNAL_ERROR",
				Message: "failed to get personal information",
				Details: map[string]interface{}{
					"error": err.Error(),
				},
			},
			Metadata: models.Metadata{},
		})
		return
	}

	processingTimeMs := time.Since(startTime).Milliseconds()

	// Build response
	items := make([]models.PersonalInfoResponse, 0)
	if personalInfoList != nil {
		for _, info := range personalInfoList {
			items = append(items, models.PersonalInfoResponse{
				ID:         info.ID,
				UserID:     info.UserID,
				Content:    info.Content,
				Category:   info.Category,
				Importance: info.Importance,
				CreatedAt:  info.CreatedAt.UTC().Format(time.RFC3339),
				UpdatedAt:  info.UpdatedAt.UTC().Format(time.RFC3339),
			})
		}
	}

	listResp := models.PersonalInfoListResponse{
		Items:  items,
		Total:  len(items),
		UserID: userID,
	}

	response := map[string]interface{}{
		"personal_info_list":   listResp,
		"processing_time_ms":   processingTimeMs,
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success:  true,
		Data:     response,
		Metadata: models.Metadata{},
	})
}

// UpdatePersonalInfo updates a personal information entry
// @Summary Update personal information
// @Description Update an existing personal information entry
// @Tags personal-info
// @Accept json
// @Produce json
// @Param info_id path string true "Personal info ID"
// @Param request body models.PersonalInfoUpdateRequest true "Personal info update request"
// @Success 200 {object} models.APIResponse "Personal info updated successfully"
// @Failure 400 {object} models.APIResponse "Invalid request"
// @Failure 404 {object} models.APIResponse "Personal info not found"
// @Failure 500 {object} models.APIResponse "Server error"
// @Router /api/rag/personal-info/{info_id} [put]
func (pih *PersonalInfoHandler) UpdatePersonalInfo(c *gin.Context) {
	startTime := time.Now()
	infoID := c.Param("info_id")

	if infoID == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.ErrorInfo{
				Code:    "INVALID_REQUEST",
				Message: "info_id is required",
				Details: map[string]interface{}{
					"field":  "info_id",
					"reason": "required parameter missing",
				},
			},
			Metadata: models.Metadata{},
		})
		return
	}

	var req models.PersonalInfoUpdateRequest

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

	// Get existing personal info
	personalInfo, err := pih.personalInfoStore.GetPersonalInfo(context.Background(), infoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: &models.ErrorInfo{
				Code:    "INTERNAL_ERROR",
				Message: "failed to get personal information",
				Details: map[string]interface{}{
					"error": err.Error(),
				},
			},
			Metadata: models.Metadata{},
		})
		return
	}

	if personalInfo == nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error: &models.ErrorInfo{
				Code:    "PERSONAL_INFO_NOT_FOUND",
				Message: "personal information not found",
				Details: map[string]interface{}{
					"info_id": infoID,
				},
			},
			Metadata: models.Metadata{},
		})
		return
	}

	// Update fields
	if req.Content != "" {
		personalInfo.Content = req.Content
	}
	if req.Category != "" {
		personalInfo.Category = req.Category
	}
	if req.Importance != "" {
		personalInfo.Importance = req.Importance
	}
	personalInfo.UpdatedAt = time.Now()

	// Save updated personal info
	if err := pih.personalInfoStore.UpdatePersonalInfo(context.Background(), personalInfo); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: &models.ErrorInfo{
				Code:    "INTERNAL_ERROR",
				Message: "failed to update personal information",
				Details: map[string]interface{}{
					"error": err.Error(),
				},
			},
			Metadata: models.Metadata{},
		})
		return
	}

	processingTimeMs := time.Since(startTime).Milliseconds()

	// Build response
	infoResp := models.PersonalInfoResponse{
		ID:         personalInfo.ID,
		UserID:     personalInfo.UserID,
		Content:    personalInfo.Content,
		Category:   personalInfo.Category,
		Importance: personalInfo.Importance,
		CreatedAt:  personalInfo.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:  personalInfo.UpdatedAt.UTC().Format(time.RFC3339),
	}

	response := map[string]interface{}{
		"personal_info":        infoResp,
		"processing_time_ms":   processingTimeMs,
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success:  true,
		Data:     response,
		Metadata: models.Metadata{},
	})
}

// DeletePersonalInfo deletes a personal information entry
// @Summary Delete personal information
// @Description Delete a personal information entry by ID
// @Tags personal-info
// @Produce json
// @Param info_id path string true "Personal info ID"
// @Success 200 {object} models.APIResponse "Personal info deleted successfully"
// @Failure 404 {object} models.APIResponse "Personal info not found"
// @Failure 500 {object} models.APIResponse "Server error"
// @Router /api/rag/personal-info/{info_id} [delete]
func (pih *PersonalInfoHandler) DeletePersonalInfo(c *gin.Context) {
	startTime := time.Now()
	infoID := c.Param("info_id")

	if infoID == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.ErrorInfo{
				Code:    "INVALID_REQUEST",
				Message: "info_id is required",
				Details: map[string]interface{}{
					"field":  "info_id",
					"reason": "required parameter missing",
				},
			},
			Metadata: models.Metadata{},
		})
		return
	}

	// Check if personal info exists
	personalInfo, err := pih.personalInfoStore.GetPersonalInfo(context.Background(), infoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: &models.ErrorInfo{
				Code:    "INTERNAL_ERROR",
				Message: "failed to get personal information",
				Details: map[string]interface{}{
					"error": err.Error(),
				},
			},
			Metadata: models.Metadata{},
		})
		return
	}

	if personalInfo == nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error: &models.ErrorInfo{
				Code:    "PERSONAL_INFO_NOT_FOUND",
				Message: "personal information not found",
				Details: map[string]interface{}{
					"info_id": infoID,
				},
			},
			Metadata: models.Metadata{},
		})
		return
	}

	// Delete personal info
	if err := pih.personalInfoStore.DeletePersonalInfo(context.Background(), infoID); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: &models.ErrorInfo{
				Code:    "INTERNAL_ERROR",
				Message: "failed to delete personal information",
				Details: map[string]interface{}{
					"error": err.Error(),
				},
			},
			Metadata: models.Metadata{},
		})
		return
	}

	processingTimeMs := time.Since(startTime).Milliseconds()

	response := map[string]interface{}{
		"deleted_info_id":      infoID,
		"processing_time_ms":   processingTimeMs,
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success:  true,
		Data:     response,
		Metadata: models.Metadata{},
	})
}