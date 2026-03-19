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

// leaveRequestService implements LeaveRequestService interface
type leaveRequestService struct {
	leaveRequestRepo repository.LeaveRequestRepository
	employeeRepo     repository.EmployeeRepository
	leaveTypeRepo    repository.LeaveTypeRepository
	leaveBalanceRepo repository.LeaveBalanceRepository
}

// NewLeaveRequestService creates a new leave request service
func NewLeaveRequestService(
	leaveRequestRepo repository.LeaveRequestRepository,
	employeeRepo repository.EmployeeRepository,
	leaveTypeRepo repository.LeaveTypeRepository,
	leaveBalanceRepo repository.LeaveBalanceRepository,
) LeaveRequestService {
	return &leaveRequestService{
		leaveRequestRepo: leaveRequestRepo,
		employeeRepo:     employeeRepo,
		leaveTypeRepo:    leaveTypeRepo,
		leaveBalanceRepo: leaveBalanceRepo,
	}
}

// CreateLeaveRequest creates a new leave request
func (s *leaveRequestService) CreateLeaveRequest(ctx context.Context, employeeID uuid.UUID, req *dto.CreateLeaveRequestRequest) (*dto.LeaveRequestDetailResponse, error) {
	// Verify employee exists
	employee, err := s.employeeRepo.GetByID(ctx, employeeID)
	if err != nil || employee == nil {
		return nil, domain.ErrEmployeeNotFound
	}

	// Verify leave type exists
	leaveType, err := s.leaveTypeRepo.GetByID(ctx, req.LeaveTypeID)
	if err != nil || leaveType == nil {
		return nil, domain.ErrLeaveTypeNotFound
	}

	// Parse dates
	startDate, err := utils.ParseDateString(req.StartDate)
	if err != nil {
		return nil, domain.NewAppError(
			domain.ErrCodeValidation,
			"invalid start date format",
			400,
			nil,
		)
	}

	endDate, err := utils.ParseDateString(req.EndDate)
	if err != nil {
		return nil, domain.NewAppError(
			domain.ErrCodeValidation,
			"invalid end date format",
			400,
			nil,
		)
	}

	// Validate date range
	if !utils.IsValidDateRange(startDate, endDate) {
		return nil, domain.NewAppError(
			domain.ErrCodeValidation,
			"end date must be after or equal to start date",
			400,
			nil,
		)
	}

	// Calculate number of days
	numberOfDays := utils.CalculateDaysBetween(startDate, endDate)
	if numberOfDays <= 0 {
		return nil, domain.NewAppError(
			domain.ErrCodeValidation,
			"invalid leave duration",
			400,
			nil,
		)
	}

	// Check for overlapping leave requests
	existingRequests, err := s.leaveRequestRepo.GetByEmployeeAndDateRange(ctx, employeeID, req.StartDate, req.EndDate)
	if err == nil && len(existingRequests) > 0 {
		return nil, domain.NewAppError(
			domain.ErrCodeConflict,
			"employee already has leave during this period",
			409,
			nil,
		)
	}

	// Check leave balance
	year := startDate.Year()
	balance, err := s.leaveBalanceRepo.GetByEmployeeAndType(ctx, employeeID, req.LeaveTypeID, year)
	if err != nil || balance == nil {
		return nil, domain.ErrInsufficientLeave
	}

	if balance.RemainingDays < int(numberOfDays) {
		return nil, domain.NewAppError(
			domain.ErrCodeValidation,
			"insufficient leave balance",
			400,
			map[string]interface{}{
				"required": numberOfDays,
				"available": balance.RemainingDays,
			},
		)
	}

	// Create leave request
	leaveRequest := &domain.LeaveRequest{
		ID:           uuid.New(),
		EmployeeID:   employeeID,
		LeaveTypeID:  req.LeaveTypeID,
		StartDate:    startDate,
		EndDate:      endDate,
		NumberOfDays: numberOfDays,
		Reason:       req.Reason,
		AttachmentURL: req.Attachment,
		Status:       domain.StatusPending,
	}

	leaveRequest, err = s.leaveRequestRepo.Create(ctx, leaveRequest)
	if err != nil {
		return nil, domain.NewAppError(
			domain.ErrCodeInternalServer,
			"failed to create leave request",
			500,
			err,
		)
	}

	// Reload with relationships
	leaveRequest, _ = s.leaveRequestRepo.GetByID(ctx, leaveRequest.ID)

	return dto.ToLeaveRequestDetailResponse(leaveRequest), nil
}

// GetLeaveRequest retrieves a leave request
func (s *leaveRequestService) GetLeaveRequest(ctx context.Context, id uuid.UUID) (*dto.LeaveRequestDetailResponse, error) {
	leaveRequest, err := s.leaveRequestRepo.GetByID(ctx, id)
	if err != nil || leaveRequest == nil {
		return nil, domain.ErrLeaveRequestNotFound
	}

	return dto.ToLeaveRequestDetailResponse(leaveRequest), nil
}

// UpdateLeaveRequest updates a leave request (only if pending)
func (s *leaveRequestService) UpdateLeaveRequest(ctx context.Context, id uuid.UUID, req *dto.UpdateLeaveRequestRequest) (*dto.LeaveRequestDetailResponse, error) {
	// Get existing leave request
	leaveRequest, err := s.leaveRequestRepo.GetByID(ctx, id)
	if err != nil || leaveRequest == nil {
		return nil, domain.ErrLeaveRequestNotFound
	}

	// Can only update if pending
	if !leaveRequest.IsPending() {
		return nil, domain.NewAppError(
			domain.ErrCodeValidation,
			"can only update pending leave requests",
			400,
			nil,
		)
	}

	// Update dates if provided
	if req.StartDate != "" && req.EndDate != "" {
		startDate, err := utils.ParseDateString(req.StartDate)
		if err != nil {
			return nil, domain.NewAppError(
				domain.ErrCodeValidation,
				"invalid start date format",
				400,
				nil,
			)
		}

		endDate, err := utils.ParseDateString(req.EndDate)
		if err != nil {
			return nil, domain.NewAppError(
				domain.ErrCodeValidation,
				"invalid end date format",
				400,
				nil,
			)
		}

		if !utils.IsValidDateRange(startDate, endDate) {
			return nil, domain.NewAppError(
				domain.ErrCodeValidation,
				"end date must be after or equal to start date",
				400,
				nil,
			)
		}

		leaveRequest.StartDate = startDate
		leaveRequest.EndDate = endDate
		leaveRequest.NumberOfDays = utils.CalculateDaysBetween(startDate, endDate)
	}

	// Update reason if provided
	if req.Reason != "" {
		leaveRequest.Reason = req.Reason
	}

	// Update attachment if provided
	if req.Attachment != "" {
		leaveRequest.AttachmentURL = req.Attachment
	}

	// Update in database
	leaveRequest, err = s.leaveRequestRepo.Update(ctx, leaveRequest)
	if err != nil {
		return nil, domain.NewAppError(
			domain.ErrCodeInternalServer,
			"failed to update leave request",
			500,
			err,
		)
	}

	// Reload with relationships
	leaveRequest, _ = s.leaveRequestRepo.GetByID(ctx, leaveRequest.ID)

	return dto.ToLeaveRequestDetailResponse(leaveRequest), nil
}

// DeleteLeaveRequest cancels a leave request (only if pending)
func (s *leaveRequestService) DeleteLeaveRequest(ctx context.Context, id uuid.UUID) error {
	leaveRequest, err := s.leaveRequestRepo.GetByID(ctx, id)
	if err != nil || leaveRequest == nil {
		return domain.ErrLeaveRequestNotFound
	}

	if !leaveRequest.IsPending() {
		return domain.NewAppError(
			domain.ErrCodeValidation,
			"can only delete pending leave requests",
			400,
			nil,
		)
	}

	return s.leaveRequestRepo.Delete(ctx, id)
}

// ListLeaveRequests retrieves all leave requests
func (s *leaveRequestService) ListLeaveRequests(ctx context.Context, page, pageSize int) (*dto.PaginatedResponse, error) {
	leaveRequests, total, err := s.leaveRequestRepo.List(ctx, page, pageSize)
	if err != nil {
		return nil, domain.NewAppError(
			domain.ErrCodeInternalServer,
			"failed to list leave requests",
			500,
			err,
		)
	}

	// Convert to DTOs
	leaveRequestDTOs := make([]dto.LeaveRequestDetailResponse, len(leaveRequests))
	for i, lr := range leaveRequests {
		leaveRequestDTOs[i] = *dto.ToLeaveRequestDetailResponse(&lr)
	}

	totalPages := utils.CalculateTotalPages(total, pageSize)

	return &dto.PaginatedResponse{
		Data:       leaveRequestDTOs,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetEmployeeLeaveRequests retrieves leave requests for an employee
func (s *leaveRequestService) GetEmployeeLeaveRequests(ctx context.Context, employeeID uuid.UUID, page, pageSize int) (*dto.PaginatedResponse, error) {
	// Verify employee exists
	employee, err := s.employeeRepo.GetByID(ctx, employeeID)
	if err != nil || employee == nil {
		return nil, domain.ErrEmployeeNotFound
	}

	leaveRequests, total, err := s.leaveRequestRepo.GetByEmployee(ctx, employeeID, page, pageSize)
	if err != nil {
		return nil, domain.NewAppError(
			domain.ErrCodeInternalServer,
			"failed to list leave requests",
			500,
			err,
		)
	}

	// Convert to DTOs
	leaveRequestDTOs := make([]dto.LeaveRequestResponse, len(leaveRequests))
	for i, lr := range leaveRequests {
		leaveRequestDTOs[i] = *dto.ToLeaveRequestResponse(&lr)
	}

	totalPages := utils.CalculateTotalPages(total, pageSize)

	return &dto.PaginatedResponse{
		Data:       leaveRequestDTOs,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetPendingApprovals retrieves pending leave requests for manager approval
func (s *leaveRequestService) GetPendingApprovals(ctx context.Context, managerID uuid.UUID, page, pageSize int) (*dto.PaginatedResponse, error) {
	// Verify manager exists
	manager, err := s.employeeRepo.GetByID(ctx, managerID)
	if err != nil || manager == nil {
		return nil, domain.ErrEmployeeNotFound
	}

	// Get all employees reporting to this manager
	subordinates, err := s.employeeRepo.ListByManager(ctx, managerID)
	if err != nil {
		return nil, domain.NewAppError(
			domain.ErrCodeInternalServer,
			"failed to get subordinates",
			500,
			err,
		)
	}

	// Get pending requests for all subordinates
	var allPendingRequests []domain.LeaveRequest
	for _, emp := range subordinates {
		requests, _, err := s.leaveRequestRepo.GetByStatus(ctx, domain.StatusPending, 1, 1000)
		if err == nil {
			for _, req := range requests {
				if req.EmployeeID == emp.ID {
					allPendingRequests = append(allPendingRequests, req)
				}
			}
		}
	}

	// Convert to DTOs
	leaveRequestDTOs := make([]dto.LeaveRequestDetailResponse, len(allPendingRequests))
	for i, lr := range allPendingRequests {
		leaveRequestDTOs[i] = *dto.ToLeaveRequestDetailResponse(&lr)
	}

	totalPages := utils.CalculateTotalPages(int64(len(allPendingRequests)), pageSize)

	return &dto.PaginatedResponse{
		Data:       leaveRequestDTOs,
		Total:      int64(len(allPendingRequests)),
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// ApproveLeaveRequest approves a leave request
func (s *leaveRequestService) ApproveLeaveRequest(ctx context.Context, requestID uuid.UUID, approverID uuid.UUID, req *dto.ApproveLeaveRequestRequest) (*dto.LeaveRequestDetailResponse, error) {
	// Get leave request
	leaveRequest, err := s.leaveRequestRepo.GetByID(ctx, requestID)
	if err != nil || leaveRequest == nil {
		return nil, domain.ErrLeaveRequestNotFound
	}

	// Must be pending
	if !leaveRequest.IsPending() {
		return nil, domain.NewAppError(
			domain.ErrCodeValidation,
			"can only approve pending leave requests",
			400,
			nil,
		)
	}

	// Verify approver exists
	approver, err := s.employeeRepo.GetByID(ctx, approverID)
	if err != nil || approver == nil {
		return nil, domain.ErrEmployeeNotFound
	}

	// Update leave request
	now := time.Now()
	leaveRequest.Status = domain.StatusApproved
	leaveRequest.ApprovedByID = &approverID
	leaveRequest.ApprovalDate = &now

	leaveRequest, err = s.leaveRequestRepo.Update(ctx, leaveRequest)
	if err != nil {
		return nil, domain.NewAppError(
			domain.ErrCodeInternalServer,
			"failed to approve leave request",
			500,
			err,
		)
	}

	// Update leave balance
	year := leaveRequest.StartDate.Year()
	balance, _ := s.leaveBalanceRepo.GetByEmployeeAndType(ctx, leaveRequest.EmployeeID, leaveRequest.LeaveTypeID, year)
	if balance != nil {
		_ = s.leaveBalanceRepo.UpdateUsedDays(ctx, balance.ID, leaveRequest.NumberOfDays)
	}

	// Reload with relationships
	leaveRequest, _ = s.leaveRequestRepo.GetByID(ctx, leaveRequest.ID)

	return dto.ToLeaveRequestDetailResponse(leaveRequest), nil
}

// RejectLeaveRequest rejects a leave request
func (s *leaveRequestService) RejectLeaveRequest(ctx context.Context, requestID uuid.UUID, approverID uuid.UUID, req *dto.RejectLeaveRequestRequest) (*dto.LeaveRequestDetailResponse, error) {
	// Get leave request
	leaveRequest, err := s.leaveRequestRepo.GetByID(ctx, requestID)
	if err != nil || leaveRequest == nil {
		return nil, domain.ErrLeaveRequestNotFound
	}

	// Must be pending
	if !leaveRequest.IsPending() {
		return nil, domain.NewAppError(
			domain.ErrCodeValidation,
			"can only reject pending leave requests",
			400,
			nil,
		)
	}

	// Verify approver exists
	approver, err := s.employeeRepo.GetByID(ctx, approverID)
	if err != nil || approver == nil {
		return nil, domain.ErrEmployeeNotFound
	}

	// Update leave request
	now := time.Now()
	leaveRequest.Status = domain.StatusRejected
	leaveRequest.ApprovedByID = &approverID
	leaveRequest.ApprovalDate = &now
	leaveRequest.RejectionReason = req.RejectionReason

	leaveRequest, err = s.leaveRequestRepo.Update(ctx, leaveRequest)
	if err != nil {
		return nil, domain.NewAppError(
			domain.ErrCodeInternalServer,
			"failed to reject leave request",
			500,
			err,
		)
	}

	// Reload with relationships
	leaveRequest, _ = s.leaveRequestRepo.GetByID(ctx, leaveRequest.ID)

	return dto.ToLeaveRequestDetailResponse(leaveRequest), nil
}