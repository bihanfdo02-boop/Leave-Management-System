package service

import (
	"context"
	"time"

	"github.com/KalinduBihan/leave-management-api/internal/domain"
	"github.com/KalinduBihan/leave-management-api/internal/repository"
	"github.com/google/uuid"
)

// dashboardService implements DashboardService interface
type dashboardService struct {
	leaveRequestRepo repository.LeaveRequestRepository
	employeeRepo     repository.EmployeeRepository
	leaveBalanceRepo repository.LeaveBalanceRepository
}

// NewDashboardService creates a new dashboard service
func NewDashboardService(
	leaveRequestRepo repository.LeaveRequestRepository,
	employeeRepo repository.EmployeeRepository,
	leaveBalanceRepo repository.LeaveBalanceRepository,
) DashboardService {
	return &dashboardService{
		leaveRequestRepo: leaveRequestRepo,
		employeeRepo:     employeeRepo,
		leaveBalanceRepo: leaveBalanceRepo,
	}
}

// GetDashboardStats retrieves dashboard statistics
func (s *dashboardService) GetDashboardStats(ctx context.Context) (map[string]interface{}, error) {
	// Get pending requests count
	pendingRequests, total, _ := s.leaveRequestRepo.GetByStatus(ctx, domain.StatusPending, 1, 1000)
	pendingCount := len(pendingRequests)

	// Get approved count (current month)
	approvedRequests, _, _ := s.leaveRequestRepo.GetByStatus(ctx, domain.StatusApproved, 1, 1000)
	monthStart := time.Now().AddDate(0, -1, 0)
	monthlyApproved := 0
	for _, req := range approvedRequests {
		if req.ApprovalDate != nil && req.ApprovalDate.After(monthStart) {
			monthlyApproved++
		}
	}

	// Get employees count
	employees, totalEmployees, _ := s.employeeRepo.List(ctx, 1, 1000)

	// Get leave type distribution
	year := time.Now().Year()
	leaveDistribution := make(map[string]int)
	for _, req := range approvedRequests {
		if req.ApprovalDate != nil && req.ApprovalDate.Year() == year {
			if req.LeaveType != nil {
				leaveDistribution[req.LeaveType.Name]++
			}
		}
	}

	stats := map[string]interface{}{
		"pending_approvals":        pendingCount,
		"total_pending":            total,
		"monthly_approved":         monthlyApproved,
		"total_employees":          totalEmployees,
		"active_employees":         len(employees),
		"leave_type_distribution":  leaveDistribution,
		"timestamp":                time.Now(),
	}

	return stats, nil
}

// GetEmployeeStats retrieves employee-specific statistics
func (s *dashboardService) GetEmployeeStats(ctx context.Context, employeeID uuid.UUID) (map[string]interface{}, error) {
	// Verify employee exists
	employee, err := s.employeeRepo.GetByID(ctx, employeeID)
	if err != nil || employee == nil {
		return nil, domain.ErrEmployeeNotFound
	}

	// Get leave requests for this employee
	requests, _, _ := s.leaveRequestRepo.GetByEmployee(ctx, employeeID, 1, 1000)

	// Count by status
	pendingCount := 0
	approvedCount := 0
	rejectedCount := 0
	totalDaysUsed := 0.0

	for _, req := range requests {
		switch req.Status {
		case domain.StatusPending:
			pendingCount++
		case domain.StatusApproved:
			approvedCount++
			totalDaysUsed += float64(req.NumberOfDays)
		case domain.StatusRejected:
			rejectedCount++
		}
	}

	// Get leave balance
	year := time.Now().Year()
	balances, _ := s.leaveBalanceRepo.GetByEmployee(ctx, employeeID, year)

	totalAvailable := 0
	totalRemaining := 0
	balanceDetails := make([]map[string]interface{}, 0)

	for _, balance := range balances {
		totalAvailable += balance.TotalDays
		totalRemaining += balance.RemainingDays
		leaveTypeName := "Unknown"
		if balance.LeaveType != nil {
			leaveTypeName = balance.LeaveType.Name
		}
		balanceDetails = append(balanceDetails, map[string]interface{}{
			"leave_type":     leaveTypeName,
			"total_days":     balance.TotalDays,
			"used_days":      balance.UsedDays,
			"remaining_days": balance.RemainingDays,
		})
	}

	departmentName := "N/A"
	if employee.Department != nil {
		departmentName = employee.Department.Name
	}

	stats := map[string]interface{}{
		"employee": map[string]interface{}{
			"id":         employee.ID,
			"name":       employee.GetFullName(),
			"email":      employee.Email,
			"department": departmentName,
		},
		"leave_requests": map[string]interface{}{
			"pending":  pendingCount,
			"approved": approvedCount,
			"rejected": rejectedCount,
			"total":    len(requests),
		},
		"leave_balance": map[string]interface{}{
			"total_available": totalAvailable,
			"total_used":      totalDaysUsed,
			"total_remaining": totalRemaining,
			"details":         balanceDetails,
		},
		"timestamp": time.Now(),
	}

	return stats, nil
}