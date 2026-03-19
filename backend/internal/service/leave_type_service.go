package service

import (
	"context"

	"github.com/KalinduBihan/leave-management-api/internal/domain"
	"github.com/KalinduBihan/leave-management-api/internal/dto"
	"github.com/KalinduBihan/leave-management-api/internal/repository"
	"github.com/KalinduBihan/leave-management-api/internal/utils"
	"github.com/google/uuid"
)

// leaveTypeService implements LeaveTypeService interface
type leaveTypeService struct {
	leaveTypeRepo repository.LeaveTypeRepository
}

// NewLeaveTypeService creates a new leave type service
func NewLeaveTypeService(leaveTypeRepo repository.LeaveTypeRepository) LeaveTypeService {
	return &leaveTypeService{
		leaveTypeRepo: leaveTypeRepo,
	}
}

// CreateLeaveType creates a new leave type
func (s *leaveTypeService) CreateLeaveType(ctx context.Context, req *dto.CreateLeaveTypeRequest) (*dto.LeaveTypeResponse, error) {
	// Validate input
	name := utils.TrimString(req.Name)
	if utils.IsEmptyString(name) {
		return nil, domain.NewAppError(
			domain.ErrCodeValidation,
			"leave type name is required",
			400,
			nil,
		)
	}

	if req.DefaultDaysPerYear < 1 {
		return nil, domain.NewAppError(
			domain.ErrCodeValidation,
			"default days per year must be at least 1",
			400,
			nil,
		)
	}

	// Check if leave type already exists
	existing, _ := s.leaveTypeRepo.GetByName(ctx, name)
	if existing != nil {
		return nil, domain.NewAppError(
			domain.ErrCodeConflict,
			"leave type already exists",
			409,
			nil,
		)
	}

	// Create leave type
	leaveType := &domain.LeaveType{
		ID:                 uuid.New(),
		Name:               name,
		Description:        req.Description,
		DefaultDaysPerYear: req.DefaultDaysPerYear,
		IsPaid:             req.IsPaid,
		RequiresApproval:   req.RequiresApproval,
		IsActive:           true,
	}

	leaveType, err := s.leaveTypeRepo.Create(ctx, leaveType)
	if err != nil {
		return nil, domain.NewAppError(
			domain.ErrCodeInternalServer,
			"failed to create leave type",
			500,
			err,
		)
	}

	return dto.ToLeaveTypeResponse(leaveType), nil
}

// GetLeaveType retrieves a leave type
func (s *leaveTypeService) GetLeaveType(ctx context.Context, id uuid.UUID) (*dto.LeaveTypeResponse, error) {
	leaveType, err := s.leaveTypeRepo.GetByID(ctx, id)
	if err != nil || leaveType == nil {
		return nil, domain.ErrLeaveTypeNotFound
	}

	return dto.ToLeaveTypeResponse(leaveType), nil
}

// UpdateLeaveType updates a leave type
func (s *leaveTypeService) UpdateLeaveType(ctx context.Context, id uuid.UUID, req *dto.UpdateLeaveTypeRequest) (*dto.LeaveTypeResponse, error) {
	// Get existing leave type
	leaveType, err := s.leaveTypeRepo.GetByID(ctx, id)
	if err != nil || leaveType == nil {
		return nil, domain.ErrLeaveTypeNotFound
	}

	// Update fields
	if req.Name != "" {
		leaveType.Name = utils.TrimString(req.Name)
	}
	if req.Description != "" {
		leaveType.Description = req.Description
	}
	if req.DefaultDaysPerYear > 0 {
		leaveType.DefaultDaysPerYear = req.DefaultDaysPerYear
	}
	if req.IsPaid != nil {
		leaveType.IsPaid = *req.IsPaid
	}
	if req.RequiresApproval != nil {
		leaveType.RequiresApproval = *req.RequiresApproval
	}
	if req.IsActive != nil {
		leaveType.IsActive = *req.IsActive
	}

	// Update in database
	leaveType, err = s.leaveTypeRepo.Update(ctx, leaveType)
	if err != nil {
		return nil, domain.NewAppError(
			domain.ErrCodeInternalServer,
			"failed to update leave type",
			500,
			err,
		)
	}

	return dto.ToLeaveTypeResponse(leaveType), nil
}

// DeleteLeaveType deletes a leave type
func (s *leaveTypeService) DeleteLeaveType(ctx context.Context, id uuid.UUID) error {
	leaveType, err := s.leaveTypeRepo.GetByID(ctx, id)
	if err != nil || leaveType == nil {
		return domain.ErrLeaveTypeNotFound
	}

	return s.leaveTypeRepo.Delete(ctx, id)
}

// ListLeaveTypes retrieves all leave types
func (s *leaveTypeService) ListLeaveTypes(ctx context.Context, page, pageSize int) (*dto.PaginatedResponse, error) {
	leaveTypes, total, err := s.leaveTypeRepo.List(ctx, page, pageSize)
	if err != nil {
		return nil, domain.NewAppError(
			domain.ErrCodeInternalServer,
			"failed to list leave types",
			500,
			err,
		)
	}

	// Convert to DTOs
	leaveTypeDTOs := make([]dto.LeaveTypeResponse, len(leaveTypes))
	for i, lt := range leaveTypes {
		leaveTypeDTOs[i] = *dto.ToLeaveTypeResponse(&lt)
	}

	totalPages := utils.CalculateTotalPages(total, pageSize)

	return &dto.PaginatedResponse{
		Data:       leaveTypeDTOs,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}