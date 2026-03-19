package dto

import (
	"github.com/KalinduBihan/leave-management-api/internal/domain"
	"github.com/google/uuid"
)

// LeaveBalanceResponse represents a leave balance in responses
type LeaveBalanceResponse struct {
	ID            uuid.UUID `json:"id"`
	EmployeeID    uuid.UUID `json:"employee_id"`
	LeaveTypeID   uuid.UUID `json:"leave_type_id"`
	TotalDays     int       `json:"total_days"`
	UsedDays      int       `json:"used_days"`
	RemainingDays int       `json:"remaining_days"`
	Year          int       `json:"year"`
}

// ToLeaveBalanceResponse converts LeaveBalance domain model to LeaveBalanceResponse
func ToLeaveBalanceResponse(lb *domain.LeaveBalance) *LeaveBalanceResponse {
	return &LeaveBalanceResponse{
		ID:            lb.ID,
		EmployeeID:    lb.EmployeeID,
		LeaveTypeID:   lb.LeaveTypeID,
		TotalDays:     lb.TotalDays,
		UsedDays:      lb.UsedDays,
		RemainingDays: lb.RemainingDays,
		Year:          lb.Year,
	}
}