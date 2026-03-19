package handler

import (
	"net/http"
	"strconv"

	"github.com/KalinduBihan/leave-management-api/internal/domain"
	"github.com/KalinduBihan/leave-management-api/internal/dto"
	"github.com/KalinduBihan/leave-management-api/internal/service"
	"github.com/KalinduBihan/leave-management-api/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// EmployeeHandler handles employee requests
type EmployeeHandler struct {
	employeeService service.EmployeeService
}

// NewEmployeeHandler creates a new employee handler
func NewEmployeeHandler(employeeService service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{
		employeeService: employeeService,
	}
}

// CreateEmployee creates a new employee
// @Summary Create Employee
// @Description Create a new employee (admin only)
// @Tags Employees
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateEmployeeRequest true "Create employee request"
// @Success 201 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Router /api/v1/employees [post]
func (h *EmployeeHandler) CreateEmployee(c *gin.Context) {
	var req dto.CreateEmployeeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseFromAppError(
			domain.NewAppError(domain.ErrCodeValidation, err.Error(), 400, nil),
		))
		return
	}

	response, err := h.employeeService.CreateEmployee(c.Request.Context(), &req)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			c.JSON(appErr.Status, dto.ErrorResponseFromAppError(appErr))
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromAppError(domain.ErrInternalServer))
		}
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse("Employee created successfully", response))
}

// GetEmployee retrieves an employee
// @Summary Get Employee
// @Description Get employee details by ID
// @Tags Employees
// @Security BearerAuth
// @Produce json
// @Param id path string true "Employee ID"
// @Success 200 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Router /api/v1/employees/{id} [get]
func (h *EmployeeHandler) GetEmployee(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseFromAppError(
			domain.NewAppError(domain.ErrCodeValidation, "invalid employee ID", 400, nil),
		))
		return
	}

	response, err := h.employeeService.GetEmployee(c.Request.Context(), id)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			c.JSON(appErr.Status, dto.ErrorResponseFromAppError(appErr))
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromAppError(domain.ErrInternalServer))
		}
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Employee retrieved successfully", response))
}

// UpdateEmployee updates an employee
// @Summary Update Employee
// @Description Update employee details
// @Tags Employees
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Employee ID"
// @Param request body dto.UpdateEmployeeRequest true "Update employee request"
// @Success 200 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Router /api/v1/employees/{id} [put]
func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseFromAppError(
			domain.NewAppError(domain.ErrCodeValidation, "invalid employee ID", 400, nil),
		))
		return
	}

	var req dto.UpdateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseFromAppError(
			domain.NewAppError(domain.ErrCodeValidation, err.Error(), 400, nil),
		))
		return
	}

	response, err := h.employeeService.UpdateEmployee(c.Request.Context(), id, &req)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			c.JSON(appErr.Status, dto.ErrorResponseFromAppError(appErr))
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromAppError(domain.ErrInternalServer))
		}
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Employee updated successfully", response))
}

// DeleteEmployee deletes an employee
// @Summary Delete Employee
// @Description Delete an employee (soft delete)
// @Tags Employees
// @Security BearerAuth
// @Produce json
// @Param id path string true "Employee ID"
// @Success 204
// @Failure 404 {object} dto.APIResponse
// @Router /api/v1/employees/{id} [delete]
func (h *EmployeeHandler) DeleteEmployee(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseFromAppError(
			domain.NewAppError(domain.ErrCodeValidation, "invalid employee ID", 400, nil),
		))
		return
	}

	err = h.employeeService.DeleteEmployee(c.Request.Context(), id)
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

// ListEmployees lists all employees
// @Summary List Employees
// @Description Get list of all employees with pagination
// @Tags Employees
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/employees [get]
func (h *EmployeeHandler) ListEmployees(c *gin.Context) {
	page, pageSize := utils.GetPaginationParams(c.DefaultQuery("page", "1"), c.DefaultQuery("page_size", "10"))

	response, err := h.employeeService.ListEmployees(c.Request.Context(), page, pageSize)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			c.JSON(appErr.Status, dto.ErrorResponseFromAppError(appErr))
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromAppError(domain.ErrInternalServer))
		}
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Employees retrieved successfully", response))
}

// GetLeaveBalance retrieves employee's leave balance
// @Summary Get Leave Balance
// @Description Get employee's leave balance for a specific year
// @Tags Employees
// @Security BearerAuth
// @Produce json
// @Param id path string true "Employee ID"
// @Param year query int false "Year" default(2026)
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/employees/{id}/balance [get]
func (h *EmployeeHandler) GetLeaveBalance(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseFromAppError(
			domain.NewAppError(domain.ErrCodeValidation, "invalid employee ID", 400, nil),
		))
		return
	}

	yearStr := c.DefaultQuery("year", strconv.Itoa(2026))
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		year = 2026
	}

	response, err := h.employeeService.GetLeaveBalance(c.Request.Context(), id, year)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			c.JSON(appErr.Status, dto.ErrorResponseFromAppError(appErr))
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromAppError(domain.ErrInternalServer))
		}
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Leave balance retrieved successfully", response))
}