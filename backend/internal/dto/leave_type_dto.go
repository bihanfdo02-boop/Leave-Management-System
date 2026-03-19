package dto

import (
	"time"

	"github.com/KalinduBihan/leave-management-api/internal/domain"
	"github.com/google/uuid"
)

// CreateLeaveTypeRequest represents a request to create a leave type
type CreateLeaveTypeRequest struct {
	Name               string `json:"name" binding:"required,min=2"`
	Description        string `json:"description" binding:"omitempty"`
	DefaultDaysPerYear int    `json:"default_days_per_year" binding:"required,min=1"`
	IsPaid             bool   `json:"is_paid" binding:""`
	RequiresApproval   bool   `json:"requires_approval" binding:""`
}

// UpdateLeaveTypeRequest represents a request to update a leave type
type UpdateLeaveTypeRequest struct {
	Name               string `json:"name" binding:"omitempty,min=2"`
	Description        string `json:"description" binding:"omitempty"`
	DefaultDaysPerYear int    `json:"default_days_per_year" binding:"omitempty,min=1"`
	IsPaid             *bool  `json:"is_paid,omitempty"`
	RequiresApproval   *bool  `json:"requires_approval,omitempty"`
	IsActive           *bool  `json:"is_active,omitempty"`
}

// LeaveTypeResponse represents a leave type in responses
type LeaveTypeResponse struct {
	ID                 uuid.UUID `json:"id"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	DefaultDaysPerYear int       `json:"default_days_per_year"`
	IsPaid             bool      `json:"is_paid"`
	RequiresApproval   bool      `json:"requires_approval"`
	IsActive           bool      `json:"is_active"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// ToLeaveTypeResponse converts LeaveType domain model to LeaveTypeResponse
func ToLeaveTypeResponse(lt *domain.LeaveType) *LeaveTypeResponse {
	return &LeaveTypeResponse{
		ID:                 lt.ID,
		Name:               lt.Name,
		Description:        lt.Description,
		DefaultDaysPerYear: lt.DefaultDaysPerYear,
		IsPaid:             lt.IsPaid,
		RequiresApproval:   lt.RequiresApproval,
		IsActive:           lt.IsActive,
		CreatedAt:          lt.CreatedAt,
		UpdatedAt:          lt.UpdatedAt,
	}
}