package service

import (
	"context"
	"time"

	"github.com/KalinduBihan/leave-management-api/internal/domain"
	"github.com/KalinduBihan/leave-management-api/internal/dto"
	"github.com/KalinduBihan/leave-management-api/internal/repository"
	"github.com/KalinduBihan/leave-management-api/internal/utils"
	"github.com/google/uuid"
)

// employeeService implements EmployeeService interface
type employeeService struct {
	employeeRepo    repository.EmployeeRepository
	userRepo        repository.UserRepository
	departmentRepo  repository.DepartmentRepository
	leaveBalanceRepo repository.LeaveBalanceRepository
	leaveTypeRepo   repository.LeaveTypeRepository
}

// NewEmployeeService creates a new employee service
func NewEmployeeService(
	employeeRepo repository.EmployeeRepository,
	userRepo repository.UserRepository,
	departmentRepo repository.DepartmentRepository,
	leaveBalanceRepo repository.LeaveBalanceRepository,
	leaveTypeRepo repository.LeaveTypeRepository,
) EmployeeService {
	return &employeeService{
		employeeRepo:    employeeRepo,
		userRepo:        userRepo,
		departmentRepo:  departmentRepo,
		leaveBalanceRepo: leaveBalanceRepo,
		leaveTypeRepo:   leaveTypeRepo,
	}
}

// CreateEmployee creates a new employee
func (s *employeeService) CreateEmployee(ctx context.Context, req *dto.CreateEmployeeRequest) (*dto.EmployeeDetailResponse, error) {
	// Validate user exists
	user, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil || user == nil {
		return nil, domain.ErrUserNotFound
	}

	// Validate department exists
	dept, err := s.departmentRepo.GetByID(ctx, req.DepartmentID)
	if err != nil || dept == nil {
		return nil, domain.ErrDepartmentNotFound
	}

	// Validate manager if provided
	if req.ManagerID != nil {
		manager, err := s.employeeRepo.GetByID(ctx, *req.ManagerID)
		if err != nil || manager == nil {
			return nil, domain.ErrEmployeeNotFound
		}
	}

	// Parse hire date
	hireDate, err := utils.ParseDateString(req.HireDate)
	if err != nil {
		return nil, domain.NewAppError(
			domain.ErrCodeValidation,
			"invalid hire date format",
			400,
			nil,
		)
	}

	// Create employee
	employee := &domain.Employee{
		ID:           uuid.New(),
		UserID:       req.UserID,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		Phone:        req.Phone,
		DepartmentID: &req.DepartmentID,
		ManagerID:    req.ManagerID,
		HireDate:     hireDate,
		IsActive:     true,
	}

	employee, err = s.employeeRepo.Create(ctx, employee)
	if err != nil {
		return nil, domain.NewAppError(
			domain.ErrCodeInternalServer,
			"failed to create employee",
			500,
			err,
		)
	}

	// Initialize leave balances for all active leave types
	leaveTypes, err := s.leaveTypeRepo.ListActive(ctx)
	if err == nil {
		currentYear := time.Now().Year()
		for _, lt := range leaveTypes {
			balance := &domain.LeaveBalance{
				ID:          uuid.New(),
				EmployeeID:  employee.ID,
				LeaveTypeID: lt.ID,
				TotalDays:   lt.DefaultDaysPerYear,
				UsedDays:    0,
				Year:        currentYear,
			}
			_, _ = s.leaveBalanceRepo.Create(ctx, balance)
		}
	}

	// Load relationships
	employee, _ = s.employeeRepo.GetByID(ctx, employee.ID)

	return dto.ToEmployeeDetailResponse(employee), nil
}

// GetEmployee retrieves an employee
func (s *employeeService) GetEmployee(ctx context.Context, id uuid.UUID) (*dto.EmployeeDetailResponse, error) {
	employee, err := s.employeeRepo.GetByID(ctx, id)
	if err != nil || employee == nil {
		return nil, domain.ErrEmployeeNotFound
	}

	return dto.ToEmployeeDetailResponse(employee), nil
}

// UpdateEmployee updates an employee
func (s *employeeService) UpdateEmployee(ctx context.Context, id uuid.UUID, req *dto.UpdateEmployeeRequest) (*dto.EmployeeDetailResponse, error) {
	// Get existing employee
	employee, err := s.employeeRepo.GetByID(ctx, id)
	if err != nil || employee == nil {
		return nil, domain.ErrEmployeeNotFound
	}

	// Validate department if provided
	if req.DepartmentID != nil {
		dept, err := s.departmentRepo.GetByID(ctx, *req.DepartmentID)
		if err != nil || dept == nil {
			return nil, domain.ErrDepartmentNotFound
		}
	}

	// Validate manager if provided
	if req.ManagerID != nil {
		manager, err := s.employeeRepo.GetByID(ctx, *req.ManagerID)
		if err != nil || manager == nil {
			return nil, domain.ErrEmployeeNotFound
		}
	}

	// Update fields
	if req.FirstName != "" {
		employee.FirstName = req.FirstName
	}
	if req.LastName != "" {
		employee.LastName = req.LastName
	}
	if req.Phone != "" {
		employee.Phone = req.Phone
	}
	if req.DepartmentID != nil {
		employee.DepartmentID = req.DepartmentID
	}
	if req.ManagerID != nil {
		employee.ManagerID = req.ManagerID
	}
	if req.IsActive != nil {
		employee.IsActive = *req.IsActive
	}

	// Update in database
	employee, err = s.employeeRepo.Update(ctx, employee)
	if err != nil {
		return nil, domain.NewAppError(
			domain.ErrCodeInternalServer,
			"failed to update employee",
			500,
			err,
		)
	}

	// Reload with relationships
	employee, _ = s.employeeRepo.GetByID(ctx, employee.ID)

	return dto.ToEmployeeDetailResponse(employee), nil
}

// DeleteEmployee deletes an employee
func (s *employeeService) DeleteEmployee(ctx context.Context, id uuid.UUID) error {
	employee, err := s.employeeRepo.GetByID(ctx, id)
	if err != nil || employee == nil {
		return domain.ErrEmployeeNotFound
	}

	return s.employeeRepo.Delete(ctx, id)
}

// ListEmployees retrieves all employees
func (s *employeeService) ListEmployees(ctx context.Context, page, pageSize int) (*dto.PaginatedResponse, error) {
	employees, total, err := s.employeeRepo.List(ctx, page, pageSize)
	if err != nil {
		return nil, domain.NewAppError(
			domain.ErrCodeInternalServer,
			"failed to list employees",
			500,
			err,
		)
	}

	// Convert to DTOs
	employeeDTOs := make([]dto.EmployeeResponse, len(employees))
	for i, emp := range employees {
		employeeDTOs[i] = *dto.ToEmployeeResponse(&emp)
	}

	totalPages := utils.CalculateTotalPages(total, pageSize)

	return &dto.PaginatedResponse{
		Data:       employeeDTOs,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetEmployeeByUserID retrieves an employee by user ID
func (s *employeeService) GetEmployeeByUserID(ctx context.Context, userID uuid.UUID) (*dto.EmployeeDetailResponse, error) {
	employee, err := s.employeeRepo.GetByUserID(ctx, userID)
	if err != nil || employee == nil {
		return nil, domain.ErrEmployeeNotFound
	}

	return dto.ToEmployeeDetailResponse(employee), nil
}

// GetLeaveBalance retrieves employee's leave balance
func (s *employeeService) GetLeaveBalance(ctx context.Context, employeeID uuid.UUID, year int) ([]dto.LeaveBalanceResponse, error) {
	// Verify employee exists
	employee, err := s.employeeRepo.GetByID(ctx, employeeID)
	if err != nil || employee == nil {
		return nil, domain.ErrEmployeeNotFound
	}

	// Get leave balances
	balances, err := s.leaveBalanceRepo.GetByEmployee(ctx, employeeID, year)
	if err != nil {
		return nil, domain.NewAppError(
			domain.ErrCodeInternalServer,
			"failed to get leave balance",
			500,
			err,
		)
	}

	// Convert to DTOs
	balanceDTOs := make([]dto.LeaveBalanceResponse, len(balances))
	for i, balance := range balances {
		balanceDTOs[i] = *dto.ToLeaveBalanceResponse(&balance)
	}

	return balanceDTOs, nil
}