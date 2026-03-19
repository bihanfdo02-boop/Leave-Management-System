package handler

import (
	"net/http"

	"github.com/KalinduBihan/leave-management-api/internal/domain"
	"github.com/KalinduBihan/leave-management-api/internal/dto"
	"github.com/KalinduBihan/leave-management-api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// DashboardHandler handles dashboard requests
type DashboardHandler struct {
	dashboardService service.DashboardService
}

// NewDashboardHandler creates a new dashboard handler
func NewDashboardHandler(dashboardService service.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

// GetDashboardStats retrieves dashboard statistics
// @Summary Get Dashboard Stats
// @Description Get overall dashboard statistics
// @Tags Dashboard
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/dashboard/stats [get]
func (h *DashboardHandler) GetDashboardStats(c *gin.Context) {
	stats, err := h.dashboardService.GetDashboardStats(c.Request.Context())
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			c.JSON(appErr.Status, dto.ErrorResponseFromAppError(appErr))
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromAppError(domain.ErrInternalServer))
		}
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Dashboard stats retrieved successfully", stats))
}

// GetEmployeeStats retrieves employee-specific statistics
// @Summary Get Employee Stats
// @Description Get statistics for a specific employee
// @Tags Dashboard
// @Security BearerAuth
// @Produce json
// @Param employee_id path string true "Employee ID"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/dashboard/employees/{employee_id}/stats [get]
func (h *DashboardHandler) GetEmployeeStats(c *gin.Context) {
	employeeIDStr := c.Param("employee_id")
	employeeID, err := uuid.Parse(employeeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseFromAppError(
			domain.NewAppError(domain.ErrCodeValidation, "invalid employee ID", 400, nil),
		))
		return
	}

	stats, err := h.dashboardService.GetEmployeeStats(c.Request.Context(), employeeID)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			c.JSON(appErr.Status, dto.ErrorResponseFromAppError(appErr))
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromAppError(domain.ErrInternalServer))
		}
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Employee stats retrieved successfully", stats))
}