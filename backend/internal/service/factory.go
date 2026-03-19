package service

import (
	"github.com/KalinduBihan/leave-management-api/config"
	"github.com/KalinduBihan/leave-management-api/internal/repository"
	"gorm.io/gorm"
)

// Services holds all service implementations
type Services struct {
	Auth        AuthService
	Employee    EmployeeService
	Department  DepartmentService
	LeaveType   LeaveTypeService
	LeaveRequest LeaveRequestService
	Dashboard   DashboardService
}

// NewServices creates and returns all services
func NewServices(repos *repository.Repositories, cfg *config.Config, db *gorm.DB) *Services {
	return &Services{
		Auth: NewAuthService(repos.User, cfg),
		Employee: NewEmployeeService(
			repos.Employee,
			repos.User,
			repos.Department,
			repos.LeaveBalance,
			repos.LeaveType,
		),
		Department: NewDepartmentService(repos.Department),
		LeaveType:  NewLeaveTypeService(repos.LeaveType),
		LeaveRequest: NewLeaveRequestService(
			repos.LeaveRequest,
			repos.Employee,
			repos.LeaveType,
			repos.LeaveBalance,
		),
		Dashboard: NewDashboardService(
			repos.LeaveRequest,
			repos.Employee,
			repos.LeaveBalance,
		),
	}
}