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

// DepartmentHandler handles department requests
type DepartmentHandler struct {
	departmentService service.DepartmentService
}

// NewDepartmentHandler creates a new department handler
func NewDepartmentHandler(departmentService service.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{
		departmentService: departmentService,
	}
}

// CreateDepartment creates a new department
// @Summary Create Department
// @Description Create a new department (admin only)
// @Tags Departments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateDepartmentRequest true "Create department request"
// @Success 201 {object} dto.APIResponse
// @Router /api/v1/departments [post]
func (h *DepartmentHandler) CreateDepartment(c *gin.Context) {
	var req dto.CreateDepartmentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseFromAppError(
			domain.NewAppError(domain.ErrCodeValidation, err.Error(), 400, nil),
		))
		return
	}

	response, err := h.departmentService.CreateDepartment(c.Request.Context(), &req)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			c.JSON(appErr.Status, dto.ErrorResponseFromAppError(appErr))
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromAppError(domain.ErrInternalServer))
		}
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse("Department created successfully", response))
}

// GetDepartment retrieves a department
// @Summary Get Department
// @Description Get department details by ID
// @Tags Departments
// @Produce json
// @Param id path string true "Department ID"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/departments/{id} [get]
func (h *DepartmentHandler) GetDepartment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseFromAppError(
			domain.NewAppError(domain.ErrCodeValidation, "invalid department ID", 400, nil),
		))
		return
	}

	response, err := h.departmentService.GetDepartment(c.Request.Context(), id)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			c.JSON(appErr.Status, dto.ErrorResponseFromAppError(appErr))
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromAppError(domain.ErrInternalServer))
		}
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Department retrieved successfully", response))
}

// UpdateDepartment updates a department
// @Summary Update Department
// @Description Update department details (admin only)
// @Tags Departments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Department ID"
// @Param request body dto.UpdateDepartmentRequest true "Update department request"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/departments/{id} [put]
func (h *DepartmentHandler) UpdateDepartment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseFromAppError(
			domain.NewAppError(domain.ErrCodeValidation, "invalid department ID", 400, nil),
		))
		return
	}

	var req dto.UpdateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseFromAppError(
			domain.NewAppError(domain.ErrCodeValidation, err.Error(), 400, nil),
		))
		return
	}

	response, err := h.departmentService.UpdateDepartment(c.Request.Context(), id, &req)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			c.JSON(appErr.Status, dto.ErrorResponseFromAppError(appErr))
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromAppError(domain.ErrInternalServer))
		}
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Department updated successfully", response))
}

// DeleteDepartment deletes a department
// @Summary Delete Department
// @Description Delete a department (admin only)
// @Tags Departments
// @Security BearerAuth
// @Produce json
// @Param id path string true "Department ID"
// @Success 204
// @Router /api/v1/departments/{id} [delete]
func (h *DepartmentHandler) DeleteDepartment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseFromAppError(
			domain.NewAppError(domain.ErrCodeValidation, "invalid department ID", 400, nil),
		))
		return
	}

	err = h.departmentService.DeleteDepartment(c.Request.Context(), id)
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

// ListDepartments lists all departments
// @Summary List Departments
// @Description Get list of all departments with pagination
// @Tags Departments
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/departments [get]
func (h *DepartmentHandler) ListDepartments(c *gin.Context) {
	page, pageSize := utils.GetPaginationParams(c.DefaultQuery("page", "1"), c.DefaultQuery("page_size", "10"))

	response, err := h.departmentService.ListDepartments(c.Request.Context(), page, pageSize)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			c.JSON(appErr.Status, dto.ErrorResponseFromAppError(appErr))
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromAppError(domain.ErrInternalServer))
		}
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Departments retrieved successfully", response))
}