package dto

import (
	"time"

	"github.com/KalinduBihan/leave-management-api/internal/domain"
	"github.com/google/uuid"
)

// CreateLeaveRequestRequest represents a request to create a leave request
type CreateLeaveRequestRequest struct {
	LeaveTypeID uuid.UUID `json:"leave_type_id" binding:"required"`
	StartDate   string    `json:"start_date" binding:"required"` // ISO 8601 format
	EndDate     string    `json:"end_date" binding:"required"`   // ISO 8601 format
	Reason      string    `json:"reason" binding:"required,min=10"`
	Attachment  string    `json:"attachment,omitempty"`
}

// UpdateLeaveRequestRequest represents a request to update a leave request (only pending)
type UpdateLeaveRequestRequest struct {
	StartDate  string `json:"start_date" binding:"omitempty"` // ISO 8601 format
	EndDate    string `json:"end_date" binding:"omitempty"`   // ISO 8601 format
	Reason     string `json:"reason" binding:"omitempty,min=10"`
	Attachment string `json:"attachment,omitempty"`
}

// ApproveLeaveRequestRequest represents a request to approve a leave request
type ApproveLeaveRequestRequest struct {
	ApprovalComment string `json:"approval_comment,omitempty"`
}

// RejectLeaveRequestRequest represents a request to reject a leave request
type RejectLeaveRequestRequest struct {
	RejectionReason string `json:"rejection_reason" binding:"required,min=5"`
}

// LeaveRequestResponse represents a leave request in responses
type LeaveRequestResponse struct {
	ID              uuid.UUID `json:"id"`
	EmployeeID      uuid.UUID `json:"employee_id"`
	LeaveTypeID     uuid.UUID `json:"leave_type_id"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	NumberOfDays    float32   `json:"number_of_days"`
	Reason          string    `json:"reason"`
	AttachmentURL   string    `json:"attachment_url,omitempty"`
	Status          string    `json:"status"`
	ApprovedByID    *uuid.UUID `json:"approved_by_id,omitempty"`
	ApprovalDate    *time.Time `json:"approval_date,omitempty"`
	RejectionReason string    `json:"rejection_reason,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// LeaveRequestDetailResponse includes related data
type LeaveRequestDetailResponse struct {
	*LeaveRequestResponse
	Employee  *EmployeeResponse `json:"employee,omitempty"`
	LeaveType *LeaveTypeResponse `json:"leave_type,omitempty"`
	ApprovedBy *EmployeeResponse `json:"approved_by,omitempty"`
}

// ToLeaveRequestResponse converts LeaveRequest domain model to LeaveRequestResponse
func ToLeaveRequestResponse(lr *domain.LeaveRequest) *LeaveRequestResponse {
	return &LeaveRequestResponse{
		ID:              lr.ID,
		EmployeeID:      lr.EmployeeID,
		LeaveTypeID:     lr.LeaveTypeID,
		StartDate:       lr.StartDate,
		EndDate:         lr.EndDate,
		NumberOfDays:    lr.NumberOfDays,
		Reason:          lr.Reason,
		AttachmentURL:   lr.AttachmentURL,
		Status:          string(lr.Status),
		ApprovedByID:    lr.ApprovedByID,
		ApprovalDate:    lr.ApprovalDate,
		RejectionReason: lr.RejectionReason,
		CreatedAt:       lr.CreatedAt,
		UpdatedAt:       lr.UpdatedAt,
	}
}

// ToLeaveRequestDetailResponse converts LeaveRequest to detailed response
func ToLeaveRequestDetailResponse(lr *domain.LeaveRequest) *LeaveRequestDetailResponse {
	detail := &LeaveRequestDetailResponse{
		LeaveRequestResponse: ToLeaveRequestResponse(lr),
	}

	if lr.Employee != nil {
		detail.Employee = ToEmployeeResponse(lr.Employee)
	}

	if lr.LeaveType != nil {
		detail.LeaveType = ToLeaveTypeResponse(lr.LeaveType)
	}

	if lr.ApprovedBy != nil {
		detail.ApprovedBy = ToEmployeeResponse(lr.ApprovedBy)
	}

	return detail
}