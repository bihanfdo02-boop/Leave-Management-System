package handler

import (
	"net/http"

	"github.com/KalinduBihan/leave-management-api/internal/domain"
	"github.com/KalinduBihan/leave-management-api/internal/dto"
	"github.com/KalinduBihan/leave-management-api/internal/service"
	"github.com/KalinduBihan/leave-management-api/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// LeaveTypeHandler handles leave type requests
type LeaveTypeHandler struct {
	leaveTypeService service.LeaveTypeService
}

// NewLeaveTypeHandler creates a new leave type handler
func NewLeaveTypeHandler(leaveTypeService service.LeaveTypeService) *LeaveTypeHandler {
	return &LeaveTypeHandler{
		leaveTypeService: leaveTypeService,
	}
}

// CreateLeaveType creates a new leave type
// @Summary Create Leave Type
// @Description Create a new leave type (admin only)
// @Tags Leave Types
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateLeaveTypeRequest true "Create leave type request"
// @Success 201 {object} dto.APIResponse
// @Router /api/v1/leave-types [post]
func (h *LeaveTypeHandler) CreateLeaveType(c *gin.Context) {
	var req dto.CreateLeaveTypeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseFromAppError(
			domain.NewAppError(domain.ErrCodeValidation, err.Error(), 400, nil),
		))
		return
	}

	response, err := h.leaveTypeService.CreateLeaveType(c.Request.Context(), &req)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			c.JSON(appErr.Status, dto.ErrorResponseFromAppError(appErr))
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromAppError(domain.ErrInternalServer))
		}
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse("Leave type created successfully", response))
}

// GetLeaveType retrieves a leave type
// @Summary Get Leave Type
// @Description Get leave type details by ID
// @Tags Leave Types
// @Produce json
// @Param id path string true "Leave Type ID"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/leave-types/{id} [get]
func (h *LeaveTypeHandler) GetLeaveType(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseFromAppError(
			domain.NewAppError(domain.ErrCodeValidation, "invalid leave type ID", 400, nil),
		))
		return
	}

	response, err := h.leaveTypeService.GetLeaveType(c.Request.Context(), id)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			c.JSON(appErr.Status, dto.ErrorResponseFromAppError(appErr))
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromAppError(domain.ErrInternalServer))
		}
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Leave type retrieved successfully", response))
}

// UpdateLeaveType updates a leave type
// @Summary Update Leave Type
// @Description Update leave type details (admin only)
// @Tags Leave Types
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Leave Type ID"
// @Param request body dto.UpdateLeaveTypeRequest true "Update leave type request"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/leave-types/{id} [put]
func (h *LeaveTypeHandler) UpdateLeaveType(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseFromAppError(
			domain.NewAppError(domain.ErrCodeValidation, "invalid leave type ID", 400, nil),
		))
		return
	}

	var req dto.UpdateLeaveTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseFromAppError(
			domain.NewAppError(domain.ErrCodeValidation, err.Error(), 400, nil),
		))
		return
	}

	response, err := h.leaveTypeService.UpdateLeaveType(c.Request.Context(), id, &req)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			c.JSON(appErr.Status, dto.ErrorResponseFromAppError(appErr))
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromAppError(domain.ErrInternalServer))
		}
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Leave type updated successfully", response))
}

// DeleteLeaveType deletes a leave type
// @Summary Delete Leave Type
// @Description Delete a leave type (admin only)
// @Tags Leave Types
// @Security BearerAuth
// @Produce json
// @Param id path string true "Leave Type ID"
// @Success 204
// @Router /api/v1/leave-types/{id} [delete]
func (h *LeaveTypeHandler) DeleteLeaveType(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseFromAppError(
			domain.NewAppError(domain.ErrCodeValidation, "invalid leave type ID", 400, nil),
		))
		return
	}

	err = h.leaveTypeService.DeleteLeaveType(c.Request.Context(), id)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			c.JSON(appErr.Status, dto.ErrorResponseFromAppError(appErr))
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromAppError(domain.ErrInternalServer))
		}
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ListLeaveTypes lists all leave types
// @Summary List Leave Types
// @Description Get list of all leave types with pagination
// @Tags Leave Types
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/leave-types [get]
func (h *LeaveTypeHandler) ListLeaveTypes(c *gin.Context) {
	page, pageSize := utils.GetPaginationParams(c.DefaultQuery("page", "1"), c.DefaultQuery("page_size", "10"))

	response, err := h.leaveTypeService.ListLeaveTypes(c.Request.Context(), page, pageSize)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			c.JSON(appErr.Status, dto.ErrorResponseFromAppError(appErr))
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromAppError(domain.ErrInternalServer))
		}
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Leave types retrieved successfully", response))
}