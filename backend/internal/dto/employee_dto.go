package dto

import (
	"time"

	"github.com/KalinduBihan/leave-management-api/internal/domain"
	"github.com/google/uuid"
)

// CreateEmployeeRequest represents a request to create an employee
type CreateEmployeeRequest struct {
	UserID       uuid.UUID `json:"user_id" binding:"required"`
	FirstName    string    `json:"first_name" binding:"required,min=2"`
	LastName     string    `json:"last_name" binding:"required,min=2"`
	Email        string    `json:"email" binding:"required,email"`
	Phone        string    `json:"phone" binding:"omitempty"`
	DepartmentID uuid.UUID `json:"department_id" binding:"required"`
	ManagerID    *uuid.UUID `json:"manager_id,omitempty"`
	HireDate     string    `json:"hire_date" binding:"required"` // ISO 8601 format
}

// UpdateEmployeeRequest represents a request to update an employee
type UpdateEmployeeRequest struct {
	FirstName    string    `json:"first_name" binding:"omitempty,min=2"`
	LastName     string    `json:"last_name" binding:"omitempty,min=2"`
	Phone        string    `json:"phone" binding:"omitempty"`
	DepartmentID *uuid.UUID `json:"department_id,omitempty"`
	ManagerID    *uuid.UUID `json:"manager_id,omitempty"`
	IsActive     *bool     `json:"is_active,omitempty"`
}

// EmployeeResponse represents an employee in responses
type EmployeeResponse struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	DepartmentID *uuid.UUID `json:"department_id,omitempty"`
	ManagerID    *uuid.UUID `json:"manager_id,omitempty"`
	HireDate     time.Time `json:"hire_date"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// EmployeeDetailResponse includes related data
type EmployeeDetailResponse struct {
	*EmployeeResponse
	Department *DepartmentResponse `json:"department,omitempty"`
	Manager    *EmployeeResponse   `json:"manager,omitempty"`
	LeaveBalance []LeaveBalanceResponse `json:"leave_balance,omitempty"`
}

// ToEmployeeResponse converts Employee domain model to EmployeeResponse
func ToEmployeeResponse(emp *domain.Employee) *EmployeeResponse {
	return &EmployeeResponse{
		ID:           emp.ID,
		UserID:       emp.UserID,
		FirstName:    emp.FirstName,
		LastName:     emp.LastName,
		Email:        emp.Email,
		Phone:        emp.Phone,
		DepartmentID: emp.DepartmentID,
		ManagerID:    emp.ManagerID,
		HireDate:     emp.HireDate,
		IsActive:     emp.IsActive,
		CreatedAt:    emp.CreatedAt,
		UpdatedAt:    emp.UpdatedAt,
	}
}

// ToEmployeeDetailResponse converts Employee to detailed response
func ToEmployeeDetailResponse(emp *domain.Employee) *EmployeeDetailResponse {
	detail := &EmployeeDetailResponse{
		EmployeeResponse: ToEmployeeResponse(emp),
	}

	if emp.Department != nil {
		detail.Department = ToDepartmentResponse(emp.Department)
	}

	if emp.Manager != nil {
		detail.Manager = ToEmployeeResponse(emp.Manager)
	}

	if len(emp.LeaveBalances) > 0 {
		detail.LeaveBalance = make([]LeaveBalanceResponse, len(emp.LeaveBalances))
		for i, lb := range emp.LeaveBalances {
			detail.LeaveBalance[i] = *ToLeaveBalanceResponse(&lb)
		}
	}

	return detail
}