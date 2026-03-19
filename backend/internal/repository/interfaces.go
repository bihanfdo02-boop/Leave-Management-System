package repository

import (
	"context"

	"github.com/KalinduBihan/leave-management-api/internal/domain"
	"github.com/google/uuid"
)

// UserRepository defines user data access methods
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *domain.User) (*domain.User, error)

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*domain.User, error)

	// Update updates a user
	Update(ctx context.Context, user *domain.User) (*domain.User, error)

	// Delete soft deletes a user
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves all users with pagination
	List(ctx context.Context, page, pageSize int) ([]domain.User, int64, error)
}

// EmployeeRepository defines employee data access methods
type EmployeeRepository interface {
	// Create creates a new employee
	Create(ctx context.Context, employee *domain.Employee) (*domain.Employee, error)

	// GetByID retrieves an employee by ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Employee, error)

	// GetByUserID retrieves an employee by user ID
	GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Employee, error)

	// GetByEmail retrieves an employee by email
	GetByEmail(ctx context.Context, email string) (*domain.Employee, error)

	// Update updates an employee
	Update(ctx context.Context, employee *domain.Employee) (*domain.Employee, error)

	// Delete soft deletes an employee
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves all employees with pagination
	List(ctx context.Context, page, pageSize int) ([]domain.Employee, int64, error)

	// ListByDepartment retrieves employees by department ID
	ListByDepartment(ctx context.Context, deptID uuid.UUID) ([]domain.Employee, error)

	// ListByManager retrieves employees by manager ID
	ListByManager(ctx context.Context, managerID uuid.UUID) ([]domain.Employee, error)
}

// DepartmentRepository defines department data access methods
type DepartmentRepository interface {
	// Create creates a new department
	Create(ctx context.Context, department *domain.Department) (*domain.Department, error)

	// GetByID retrieves a department by ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Department, error)

	// GetByName retrieves a department by name
	GetByName(ctx context.Context, name string) (*domain.Department, error)

	// Update updates a department
	Update(ctx context.Context, department *domain.Department) (*domain.Department, error)

	// Delete soft deletes a department
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves all departments with pagination
	List(ctx context.Context, page, pageSize int) ([]domain.Department, int64, error)
}

// LeaveTypeRepository defines leave type data access methods
type LeaveTypeRepository interface {
	// Create creates a new leave type
	Create(ctx context.Context, leaveType *domain.LeaveType) (*domain.LeaveType, error)

	// GetByID retrieves a leave type by ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.LeaveType, error)

	// GetByName retrieves a leave type by name
	GetByName(ctx context.Context, name string) (*domain.LeaveType, error)

	// Update updates a leave type
	Update(ctx context.Context, leaveType *domain.LeaveType) (*domain.LeaveType, error)

	// Delete soft deletes a leave type
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves all leave types with pagination
	List(ctx context.Context, page, pageSize int) ([]domain.LeaveType, int64, error)

	// ListActive retrieves all active leave types
	ListActive(ctx context.Context) ([]domain.LeaveType, error)
}

// LeaveBalanceRepository defines leave balance data access methods
type LeaveBalanceRepository interface {
	// Create creates a new leave balance
	Create(ctx context.Context, balance *domain.LeaveBalance) (*domain.LeaveBalance, error)

	// GetByID retrieves a leave balance by ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.LeaveBalance, error)

	// GetByEmployeeAndType retrieves balance for an employee and leave type in a year
	GetByEmployeeAndType(ctx context.Context, employeeID, leaveTypeID uuid.UUID, year int) (*domain.LeaveBalance, error)

	// GetByEmployee retrieves all balances for an employee in a year
	GetByEmployee(ctx context.Context, employeeID uuid.UUID, year int) ([]domain.LeaveBalance, error)

	// Update updates a leave balance
	Update(ctx context.Context, balance *domain.LeaveBalance) (*domain.LeaveBalance, error)

	// Delete deletes a leave balance
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves all leave balances with pagination
	List(ctx context.Context, page, pageSize int) ([]domain.LeaveBalance, int64, error)

	// UpdateUsedDays updates the used days for a balance
	UpdateUsedDays(ctx context.Context, balanceID uuid.UUID, daysToAdd float32) error
}

// LeaveRequestRepository defines leave request data access methods
type LeaveRequestRepository interface {
	// Create creates a new leave request
	Create(ctx context.Context, request *domain.LeaveRequest) (*domain.LeaveRequest, error)

	// GetByID retrieves a leave request by ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.LeaveRequest, error)

	// GetByEmployee retrieves all requests for an employee
	GetByEmployee(ctx context.Context, employeeID uuid.UUID, page, pageSize int) ([]domain.LeaveRequest, int64, error)

	// GetPending retrieves all pending leave requests
	GetPending(ctx context.Context, page, pageSize int) ([]domain.LeaveRequest, int64, error)

	// GetByStatus retrieves leave requests by status
	GetByStatus(ctx context.Context, status domain.LeaveStatus, page, pageSize int) ([]domain.LeaveRequest, int64, error)

	// Update updates a leave request
	Update(ctx context.Context, request *domain.LeaveRequest) (*domain.LeaveRequest, error)

	// Delete deletes a leave request
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves all leave requests with pagination
	List(ctx context.Context, page, pageSize int) ([]domain.LeaveRequest, int64, error)

	// GetByEmployeeAndDateRange retrieves requests for an employee within date range
	GetByEmployeeAndDateRange(ctx context.Context, employeeID uuid.UUID, startDate, endDate string) ([]domain.LeaveRequest, error)
}

// AuditLogRepository defines audit log data access methods
type AuditLogRepository interface {
	// Create creates a new audit log
	Create(ctx context.Context, log *domain.AuditLog) (*domain.AuditLog, error)

	// GetByID retrieves an audit log by ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.AuditLog, error)

	// GetByUser retrieves audit logs by user ID
	GetByUser(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]domain.AuditLog, int64, error)

	// GetByEntity retrieves audit logs by entity type and ID
	GetByEntity(ctx context.Context, entityType string, entityID uuid.UUID) ([]domain.AuditLog, error)

	// List retrieves all audit logs with pagination
	List(ctx context.Context, page, pageSize int) ([]domain.AuditLog, int64, error)
}