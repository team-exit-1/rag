package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthResponse represents a health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

// HealthCheck handles health check requests
// @Summary Health check
// @Description Check if the RAG server is healthy
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status:  "ok",
		Service: "rag-server",
	})
}
