package service

import (
	"context"

	"github.com/KalinduBihan/leave-management-api/internal/dto"
	"github.com/google/uuid"
)

// AuthService defines authentication operations
type AuthService interface {
	// Register registers a new user
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.LoginResponse, error)

	// Login authenticates a user
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)

	// ValidateToken validates a JWT token
	ValidateToken(token string) (uuid.UUID, error)

	// RefreshToken refreshes a JWT token
	RefreshToken(ctx context.Context, userID uuid.UUID) (string, error)
}

// EmployeeService defines employee operations
type EmployeeService interface {
	// CreateEmployee creates a new employee
	CreateEmployee(ctx context.Context, req *dto.CreateEmployeeRequest) (*dto.EmployeeDetailResponse, error)

	// GetEmployee retrieves an employee
	GetEmployee(ctx context.Context, id uuid.UUID) (*dto.EmployeeDetailResponse, error)

	// UpdateEmployee updates an employee
	UpdateEmployee(ctx context.Context, id uuid.UUID, req *dto.UpdateEmployeeRequest) (*dto.EmployeeDetailResponse, error)

	// DeleteEmployee deletes an employee
	DeleteEmployee(ctx context.Context, id uuid.UUID) error

	// ListEmployees retrieves all employees
	ListEmployees(ctx context.Context, page, pageSize int) (*dto.PaginatedResponse, error)

	// GetEmployeeByUserID retrieves an employee by user ID
	GetEmployeeByUserID(ctx context.Context, userID uuid.UUID) (*dto.EmployeeDetailResponse, error)

	// GetLeaveBalance retrieves employee's leave balance
	GetLeaveBalance(ctx context.Context, employeeID uuid.UUID, year int) ([]dto.LeaveBalanceResponse, error)
}

// DepartmentService defines department operations
type DepartmentService interface {
	// CreateDepartment creates a new department
	CreateDepartment(ctx context.Context, req *dto.CreateDepartmentRequest) (*dto.DepartmentDetailResponse, error)

	// GetDepartment retrieves a department
	GetDepartment(ctx context.Context, id uuid.UUID) (*dto.DepartmentDetailResponse, error)

	// UpdateDepartment updates a department
	UpdateDepartment(ctx context.Context, id uuid.UUID, req *dto.UpdateDepartmentRequest) (*dto.DepartmentDetailResponse, error)

	// DeleteDepartment deletes a department
	DeleteDepartment(ctx context.Context, id uuid.UUID) error

	// ListDepartments retrieves all departments
	ListDepartments(ctx context.Context, page, pageSize int) (*dto.PaginatedResponse, error)
}

// LeaveTypeService defines leave type operations
type LeaveTypeService interface {
	// CreateLeaveType creates a new leave type
	CreateLeaveType(ctx context.Context, req *dto.CreateLeaveTypeRequest) (*dto.LeaveTypeResponse, error)

	// GetLeaveType retrieves a leave type
	GetLeaveType(ctx context.Context, id uuid.UUID) (*dto.LeaveTypeResponse, error)

	// UpdateLeaveType updates a leave type
	UpdateLeaveType(ctx context.Context, id uuid.UUID, req *dto.UpdateLeaveTypeRequest) (*dto.LeaveTypeResponse, error)

	// DeleteLeaveType deletes a leave type
	DeleteLeaveType(ctx context.Context, id uuid.UUID) error

	// ListLeaveTypes retrieves all leave types
	ListLeaveTypes(ctx context.Context, page, pageSize int) (*dto.PaginatedResponse, error)
}

// LeaveRequestService defines leave request operations
type LeaveRequestService interface {
	// CreateLeaveRequest creates a new leave request
	CreateLeaveRequest(ctx context.Context, employeeID uuid.UUID, req *dto.CreateLeaveRequestRequest) (*dto.LeaveRequestDetailResponse, error)

	// GetLeaveRequest retrieves a leave request
	GetLeaveRequest(ctx context.Context, id uuid.UUID) (*dto.LeaveRequestDetailResponse, error)

	// UpdateLeaveRequest updates a leave request (only if pending)
	UpdateLeaveRequest(ctx context.Context, id uuid.UUID, req *dto.UpdateLeaveRequestRequest) (*dto.LeaveRequestDetailResponse, error)

	// DeleteLeaveRequest cancels a leave request (only if pending)
	DeleteLeaveRequest(ctx context.Context, id uuid.UUID) error

	// ListLeaveRequests retrieves all leave requests
	ListLeaveRequests(ctx context.Context, page, pageSize int) (*dto.PaginatedResponse, error)

	// GetEmployeeLeaveRequests retrieves leave requests for an employee
	GetEmployeeLeaveRequests(ctx context.Context, employeeID uuid.UUID, page, pageSize int) (*dto.PaginatedResponse, error)

	// GetPendingApprovals retrieves pending leave requests for manager approval
	GetPendingApprovals(ctx context.Context, managerID uuid.UUID, page, pageSize int) (*dto.PaginatedResponse, error)

	// ApproveLeaveRequest approves a leave request
	ApproveLeaveRequest(ctx context.Context, requestID uuid.UUID, approverID uuid.UUID, req *dto.ApproveLeaveRequestRequest) (*dto.LeaveRequestDetailResponse, error)

	// RejectLeaveRequest rejects a leave request
	RejectLeaveRequest(ctx context.Context, requestID uuid.UUID, approverID uuid.UUID, req *dto.RejectLeaveRequestRequest) (*dto.LeaveRequestDetailResponse, error)
}

// DashboardService defines dashboard operations
type DashboardService interface {
	// GetDashboardStats retrieves dashboard statistics
	GetDashboardStats(ctx context.Context) (map[string]interface{}, error)

	// GetEmployeeStats retrieves employee-specific statistics
	GetEmployeeStats(ctx context.Context, employeeID uuid.UUID) (map[string]interface{}, error)
}