package handler

import (
	"net/http"

	"github.com/KalinduBihan/leave-management-api/internal/dto"
	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check requests
type HealthHandler struct {
}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthCheck returns server health status
// @Summary Health Check
// @Description Check if server is running
// @Tags Health
// @Produce json
// @Success 200 {object} dto.APIResponse
// @Router /health [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, dto.SuccessResponse("Server is running", map[string]string{
		"status": "healthy",
	}))
}