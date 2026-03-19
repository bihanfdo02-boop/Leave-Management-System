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

// LeaveRequestHandler handles leave request requests
type LeaveRequestHandler struct {
	leaveRequestService service.LeaveRequestService
}

// NewLeaveRequestHandler creates a new leave request handler
func NewLeaveRequestHandler(leaveRequestService service.LeaveRequestService) *LeaveRequestHandler {
	return &LeaveRequestHandler{
		leaveRequestService: leaveRequestService,
	}
}

// CreateLeaveRequest creates a new leave request
// @Summary Create Leave Request
// @Description Create a new leave request
// @Tags Leave Requests
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateLeaveRequestRequest true "Create leave request"
// @Success 201 {object} dto.APIResponse
// @Router /api/v1/leaves [post]
func (h *LeaveRequestHandler) CreateLeaveRequest(c *gin.Context) {
	var req dto.CreateLeaveRequestRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseFromAppError(
			domain.NewAppError(domain.ErrCodeValidation, err.Error(), 400, nil),
		))
		return
	}

	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponseFromAppError(domain.ErrUnauthorized))
		return
	}

	_, err := parseUserID(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponseFromAppError(domain.ErrUnauthorized))
		return
	}

	// Get employee ID from user ID (would need employee service call)
	// For now, we'll use a parameter or extract from employee lookup
	employeeIDStr := c.Query("employee_id")
	if employeeIDStr == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseFromAppError(
			domain.NewAppError(domain.ErrCodeValidation, "employee_id is required", 400, nil),
		))
		return
	}

	employeeID, err := uuid.Parse(employeeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseFromAppError(
			domain.NewAppError(domain.ErrCodeValidation, "invalid employee_id", 400, nil),
		))
		return
	}

	response, err := h.leaveRequestService.CreateLeaveRequest(c.Request.Context(), employeeID, &req)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			c.JSON(appErr.Status, dto.ErrorResponseFromAppError(appErr))
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromAppError(domain.ErrInternalServer))
		}
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse("Leave request created successfully", response))
}

// GetLeaveRequest retrieves a leave request
// @Summary Get Leave Request
// @Description Get leave request details by ID
// @Tags Leave Requests
// @Security BearerAuth
// @Produce json
// @Param id path string true "Leave Request ID"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/leaves/{id} [get]
func (h *LeaveRequestHandler) GetLeaveRequest(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseFromAppError(
			domain.NewAppError(domain.ErrCodeValidation, "invalid leave request ID", 400, nil),
		))
		return
	}

	response, err := h.leaveRequestService.GetLeaveRequest(c.Request.Context(), id)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			c.JSON(appErr.Status, dto.ErrorResponseFromAppError(appErr))
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromAppError(domain.ErrInternalServer))
		}
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Leave request retrieved successfully", response))
}

// UpdateLeaveRequest updates a leave request
// @Summary Update Leave Request
// @Description Update leave request (only if pending)
// @Tags Leave Requests
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Leave Request ID"
// @Param request body dto.UpdateLeaveRequestRequest true "Update leave request"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/leaves/{id} [put]
func (h *LeaveRequestHandler) UpdateLeaveRequest(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseFromAppError(
			domain.NewAppError(domain.ErrCodeValidation, "invalid leave request ID", 400, nil),
		))
		return
	}

	var req dto.UpdateLeaveRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseFromAppError(
			domain.NewAppError(domain.ErrCodeValidation, err.Error(), 400, nil),
		))
		return
	}

	response, err := h.leaveRequestService.UpdateLeaveRequest(c.Request.Context(), id, &req)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			c.JSON(appErr.Status, dto.ErrorResponseFromAppError(appErr))
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromAppError(domain.ErrInternalServer))
		}
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Leave request updated successfully", response))
}

// DeleteLeaveRequest cancels a leave request
// @Summary Delete Leave Request
// @Description Cancel a leave request (only if pending)
// @Tags Leave Requests
// @Security BearerAuth
// @Produce json
// @Param id path string true "Leave Request ID"
// @Success 204
// @Router /api/v1/leaves/{id} [delete]
func (h *LeaveRequestHandler) DeleteLeaveRequest(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseFromAppError(
			domain.NewAppError(domain.ErrCodeValidation, "invalid leave request ID", 400, nil),
		))
		return
	}

	err = h.leaveRequestService.DeleteLeaveRequest(c.Request.Context(), id)
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

// ListLeaveRequests lists all leave requests
// @Summary List Leave Requests
// @Description Get list of all leave requests with pagination
// @Tags Leave Requests
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/leaves [get]
func (h *LeaveRequestHandler) ListLeaveRequests(c *gin.Context) {
	page, pageSize := utils.GetPaginationParams(c.DefaultQuery("page", "1"), c.DefaultQuery("page_size", "10"))

	response, err := h.leaveRequestService.ListLeaveRequests(c.Request.Context(), page, pageSize)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			c.JSON(appErr.Status, dto.ErrorResponseFromAppError(appErr))
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromAppError(domain.ErrInternalServer))
		}
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Leave requests retrieved successfully", response))
}

// GetEmployeeLeaveRequests gets leave requests for an employee
// @Summary Get Employee Leave Requests
// @Description Get leave requests for a specific employee
// @Tags Leave Requests
// @Security BearerAuth
// @Produce json
// @Param employee_id path string true "Employee ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/employees/{employee_id}/leaves [get]
func (h *LeaveRequestHandler) GetEmployeeLeaveRequests(c *gin.Context) {
	employeeID, err := uuid.Parse(c.Param("employee_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseFromAppError(
			domain.NewAppError(domain.ErrCodeValidation, "invalid employee ID", 400, nil),
		))
		return
	}

	page, pageSize := utils.GetPaginationParams(c.DefaultQuery("page", "1"), c.DefaultQuery("page_size", "10"))

	response, err := h.leaveRequestService.GetEmployeeLeaveRequests(c.Request.Context(), employeeID, page, pageSize)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			c.JSON(appErr.Status, dto.ErrorResponseFromAppError(appErr))
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromAppError(domain.ErrInternalServer))
		}
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Leave requests retrieved successfully", response))
}

// GetPendingApprovals gets pending leave requests for a manager
// @Summary Get Pending Approvals
// @Description Get pending leave requests for manager approval
// @Tags Leave Requests
// @Security BearerAuth
// @Produce json
// @Param manager_id path string true "Manager ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/managers/{manager_id}/pending-approvals [get]
func (h *LeaveRequestHandler) GetPendingApprovals(c *gin.Context) {
	managerID, err := uuid.Parse(c.Param("manager_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseFromAppError(
			domain.NewAppError(domain.ErrCodeValidation, "invalid manager ID", 400, nil),
		))
		return
	}

	page, pageSize := utils.GetPaginationParams(c.DefaultQuery("page", "1"), c.DefaultQuery("page_size", "10"))

	response, err := h.leaveRequestService.GetPendingApprovals(c.Request.Context(), managerID, page, pageSize)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			c.JSON(appErr.Status, dto.ErrorResponseFromAppError(appErr))
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromAppError(domain.ErrInternalServer))
		}
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Pending approvals retrieved successfully", response))
}

// ApproveLeaveRequest approves a leave request
// @Summary Approve Leave Request
// @Description Approve a pending leave request
// @Tags Leave Requests
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Leave Request ID"
// @Param request body dto.ApproveLeaveRequestRequest true "Approval request"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/leaves/{id}/approve [post]
func (h *LeaveRequestHandler) ApproveLeaveRequest(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseFromAppError(
			domain.NewAppError(domain.ErrCodeValidation, "invalid leave request ID", 400, nil),
		))
		return
	}

	var req dto.ApproveLeaveRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseFromAppError(
			domain.NewAppError(domain.ErrCodeValidation, err.Error(), 400, nil),
		))
		return
	}

	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponseFromAppError(domain.ErrUnauthorized))
		return
	}

	approverID, err := parseUserID(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponseFromAppError(domain.ErrUnauthorized))
		return
	}

	// Get approver employee ID from user ID (would need employee service call)
	approverEmployeeIDStr := c.Query("approver_employee_id")
	if approverEmployeeIDStr == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseFromAppError(
			domain.NewAppError(domain.ErrCodeValidation, "approver_employee_id is required", 400, nil),
		))
		return
	}

	approverEmployeeID, err := uuid.Parse(approverEmployeeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseFromAppError(
			domain.NewAppError(domain.ErrCodeValidation, "invalid approver_employee_id", 400, nil),
		))
		return
	}

	_ = approverID // Use userID if needed

	response, err := h.leaveRequestService.ApproveLeaveRequest(c.Request.Context(), id, approverEmployeeID, &req)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			c.JSON(appErr.Status, dto.ErrorResponseFromAppError(appErr))
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromAppError(domain.ErrInternalServer))
		}
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Leave request approved successfully", response))
}

// RejectLeaveRequest rejects a leave request
// @Summary Reject Leave Request
// @Description Reject a pending leave request
// @Tags Leave Requests
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Leave Request ID"
// @Param request body dto.RejectLeaveRequestRequest true "Rejection request"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/leaves/{id}/reject [post]
func (h *LeaveRequestHandler) RejectLeaveRequest(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseFromAppError(
			domain.NewAppError(domain.ErrCodeValidation, "invalid leave request ID", 400, nil),
		))
		return
	}

	var req dto.RejectLeaveRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseFromAppError(
			domain.NewAppError(domain.ErrCodeValidation, err.Error(), 400, nil),
		))
		return
	}

	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponseFromAppError(domain.ErrUnauthorized))
		return
	}

	approverID, err := parseUserID(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponseFromAppError(domain.ErrUnauthorized))
		return
	}

	// Get approver employee ID from user ID
	approverEmployeeIDStr := c.Query("approver_employee_id")
	if approverEmployeeIDStr == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseFromAppError(
			domain.NewAppError(domain.ErrCodeValidation, "approver_employee_id is required", 400, nil),
		))
		return
	}

	approverEmployeeID, err := uuid.Parse(approverEmployeeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseFromAppError(
			domain.NewAppError(domain.ErrCodeValidation, "invalid approver_employee_id", 400, nil),
		))
		return
	}

	_ = approverID // Use userID if needed

	response, err := h.leaveRequestService.RejectLeaveRequest(c.Request.Context(), id, approverEmployeeID, &req)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			c.JSON(appErr.Status, dto.ErrorResponseFromAppError(appErr))
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromAppError(domain.ErrInternalServer))
		}
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Leave request rejected successfully", response))
}