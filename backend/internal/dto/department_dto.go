package dto

import (
	"time"

	"github.com/KalinduBihan/leave-management-api/internal/domain"
	"github.com/google/uuid"
)

// CreateDepartmentRequest represents a request to create a department
type CreateDepartmentRequest struct {
	Name        string `json:"name" binding:"required,min=2"`
	Description string `json:"description" binding:"omitempty"`
}

// UpdateDepartmentRequest represents a request to update a department
type UpdateDepartmentRequest struct {
	Name        string `json:"name" binding:"omitempty,min=2"`
	Description string `json:"description" binding:"omitempty"`
	IsActive    *bool  `json:"is_active,omitempty"`
}

// DepartmentResponse represents a department in responses
type DepartmentResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// DepartmentDetailResponse includes employee count
type DepartmentDetailResponse struct {
	*DepartmentResponse
	EmployeeCount int `json:"employee_count"`
}

// ToDepartmentResponse converts Department domain model to DepartmentResponse
func ToDepartmentResponse(dept *domain.Department) *DepartmentResponse {
	return &DepartmentResponse{
		ID:          dept.ID,
		Name:        dept.Name,
		Description: dept.Description,
		IsActive:    dept.IsActive,
		CreatedAt:   dept.CreatedAt,
		UpdatedAt:   dept.UpdatedAt,
	}
}

// ToDepartmentDetailResponse converts Department to detailed response
func ToDepartmentDetailResponse(dept *domain.Department) *DepartmentDetailResponse {
	return &DepartmentDetailResponse{
		DepartmentResponse: ToDepartmentResponse(dept),
		EmployeeCount:      len(dept.Employees),
	}
}