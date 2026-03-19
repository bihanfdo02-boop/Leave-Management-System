package repository

import "gorm.io/gorm"

// Repositories holds all repository implementations
type Repositories struct {
	User        UserRepository
	Employee    EmployeeRepository
	Department  DepartmentRepository
	LeaveType   LeaveTypeRepository
	LeaveBalance LeaveBalanceRepository
	LeaveRequest LeaveRequestRepository
	AuditLog    AuditLogRepository
}

// NewRepositories creates and returns all repositories
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User:         NewUserRepository(db),
		Employee:     NewEmployeeRepository(db),
		Department:   NewDepartmentRepository(db),
		LeaveType:    NewLeaveTypeRepository(db),
		LeaveBalance: NewLeaveBalanceRepository(db),
		LeaveRequest: NewLeaveRequestRepository(db),
		AuditLog:     NewAuditLogRepository(db),
	}
}