package service

import (
	"context"

	"github.com/KalinduBihan/leave-management-api/internal/domain"
	"github.com/KalinduBihan/leave-management-api/internal/dto"
	"github.com/KalinduBihan/leave-management-api/internal/repository"
	"github.com/KalinduBihan/leave-management-api/internal/utils"
	"github.com/google/uuid"
)

// departmentService implements DepartmentService interface
type departmentService struct {
	departmentRepo repository.DepartmentRepository
}

// NewDepartmentService creates a new department service
func NewDepartmentService(departmentRepo repository.DepartmentRepository) DepartmentService {
	return &departmentService{
		departmentRepo: departmentRepo,
	}
}

// CreateDepartment creates a new department
func (s *departmentService) CreateDepartment(ctx context.Context, req *dto.CreateDepartmentRequest) (*dto.DepartmentDetailResponse, error) {
	// Validate input
	name := utils.TrimString(req.Name)
	if utils.IsEmptyString(name) {
		return nil, domain.NewAppError(
			domain.ErrCodeValidation,
			"department name is required",
			400,
			nil,
		)
	}

	// Check if department already exists
	existing, _ := s.departmentRepo.GetByName(ctx, name)
	if existing != nil {
		return nil, domain.ErrDepartmentExists
	}

	// Create department
	department := &domain.Department{
		ID:          uuid.New(),
		Name:        name,
		Description: req.Description,
		IsActive:    true,
	}

	department, err := s.departmentRepo.Create(ctx, department)
	if err != nil {
		return nil, domain.NewAppError(
			domain.ErrCodeInternalServer,
			"failed to create department",
			500,
			err,
		)
	}

	return dto.ToDepartmentDetailResponse(department), nil
}

// GetDepartment retrieves a department
func (s *departmentService) GetDepartment(ctx context.Context, id uuid.UUID) (*dto.DepartmentDetailResponse, error) {
	department, err := s.departmentRepo.GetByID(ctx, id)
	if err != nil || department == nil {
		return nil, domain.ErrDepartmentNotFound
	}

	return dto.ToDepartmentDetailResponse(department), nil
}

// UpdateDepartment updates a department
func (s *departmentService) UpdateDepartment(ctx context.Context, id uuid.UUID, req *dto.UpdateDepartmentRequest) (*dto.DepartmentDetailResponse, error) {
	// Get existing department
	department, err := s.departmentRepo.GetByID(ctx, id)
	if err != nil || department == nil {
		return nil, domain.ErrDepartmentNotFound
	}

	// Update fields
	if req.Name != "" {
		department.Name = utils.TrimString(req.Name)
	}
	if req.Description != "" {
		department.Description = req.Description
	}
	if req.IsActive != nil {
		department.IsActive = *req.IsActive
	}

	// Update in database
	department, err = s.departmentRepo.Update(ctx, department)
	if err != nil {
		return nil, domain.NewAppError(
			domain.ErrCodeInternalServer,
			"failed to update department",
			500,
			err,
		)
	}

	// Reload with employees
	department, _ = s.departmentRepo.GetByID(ctx, department.ID)

	return dto.ToDepartmentDetailResponse(department), nil
}

// DeleteDepartment deletes a department
func (s *departmentService) DeleteDepartment(ctx context.Context, id uuid.UUID) error {
	department, err := s.departmentRepo.GetByID(ctx, id)
	if err != nil || department == nil {
		return domain.ErrDepartmentNotFound
	}

	return s.departmentRepo.Delete(ctx, id)
}

// ListDepartments retrieves all departments
func (s *departmentService) ListDepartments(ctx context.Context, page, pageSize int) (*dto.PaginatedResponse, error) {
	departments, total, err := s.departmentRepo.List(ctx, page, pageSize)
	if err != nil {
		return nil, domain.NewAppError(
			domain.ErrCodeInternalServer,
			"failed to list departments",
			500,
			err,
		)
	}

	// Convert to DTOs
	departmentDTOs := make([]dto.DepartmentDetailResponse, len(departments))
	for i, dept := range departments {
		departmentDTOs[i] = *dto.ToDepartmentDetailResponse(&dept)
	}

	totalPages := utils.CalculateTotalPages(total, pageSize)

	return &dto.PaginatedResponse{
		Data:       departmentDTOs,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}