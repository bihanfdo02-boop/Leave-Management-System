package repository

import (
	"context"
	"errors"
	"time"

	"github.com/KalinduBihan/leave-management-api/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PostgresLeaveRequestRepository implements LeaveRequestRepository interface
type PostgresLeaveRequestRepository struct {
	db *gorm.DB
}

// NewLeaveRequestRepository creates a new leave request repository
func NewLeaveRequestRepository(db *gorm.DB) LeaveRequestRepository {
	return &PostgresLeaveRequestRepository{db: db}
}

// Create creates a new leave request
func (r *PostgresLeaveRequestRepository) Create(ctx context.Context, request *domain.LeaveRequest) (*domain.LeaveRequest, error) {
	if err := r.db.WithContext(ctx).Create(request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

// GetByID retrieves a leave request by ID with relationships
func (r *PostgresLeaveRequestRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.LeaveRequest, error) {
	var request domain.LeaveRequest
	if err := r.db.WithContext(ctx).
		Preload("Employee").
		Preload("LeaveType").
		Preload("ApprovedBy").
		Where("id = ?", id).
		First(&request).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrLeaveRequestNotFound
		}
		return nil, err
	}
	return &request, nil
}

// GetByEmployee retrieves all requests for an employee
func (r *PostgresLeaveRequestRepository) GetByEmployee(ctx context.Context, employeeID uuid.UUID, page, pageSize int) ([]domain.LeaveRequest, int64, error) {
	var requests []domain.LeaveRequest
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).
		Model(&domain.LeaveRequest{}).
		Where("employee_id = ?", employeeID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).
		Preload("LeaveType").
		Where("employee_id = ?", employeeID).
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&requests).Error; err != nil {
		return nil, 0, err
	}

	return requests, total, nil
}

// GetPending retrieves all pending leave requests
func (r *PostgresLeaveRequestRepository) GetPending(ctx context.Context, page, pageSize int) ([]domain.LeaveRequest, int64, error) {
	var requests []domain.LeaveRequest
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).
		Model(&domain.LeaveRequest{}).
		Where("status = ?", domain.StatusPending).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).
		Preload("Employee").
		Preload("LeaveType").
		Where("status = ?", domain.StatusPending).
		Offset(offset).
		Limit(pageSize).
		Order("created_at ASC").
		Find(&requests).Error; err != nil {
		return nil, 0, err
	}

	return requests, total, nil
}

// GetByStatus retrieves leave requests by status
func (r *PostgresLeaveRequestRepository) GetByStatus(ctx context.Context, status domain.LeaveStatus, page, pageSize int) ([]domain.LeaveRequest, int64, error) {
	var requests []domain.LeaveRequest
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).
		Model(&domain.LeaveRequest{}).
		Where("status = ?", status).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).
		Preload("Employee").
		Preload("LeaveType").
		Where("status = ?", status).
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&requests).Error; err != nil {
		return nil, 0, err
	}

	return requests, total, nil
}

// Update updates a leave request
func (r *PostgresLeaveRequestRepository) Update(ctx context.Context, request *domain.LeaveRequest) (*domain.LeaveRequest, error) {
	if err := r.db.WithContext(ctx).Save(request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

// Delete deletes a leave request
func (r *PostgresLeaveRequestRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.LeaveRequest{}, "id = ?", id).Error
}

// List retrieves all leave requests with pagination
func (r *PostgresLeaveRequestRepository) List(ctx context.Context, page, pageSize int) ([]domain.LeaveRequest, int64, error) {
	var requests []domain.LeaveRequest
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&domain.LeaveRequest{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).
		Preload("Employee").
		Preload("LeaveType").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&requests).Error; err != nil {
		return nil, 0, err
	}

	return requests, total, nil
}

// GetByEmployeeAndDateRange retrieves requests for an employee within date range
func (r *PostgresLeaveRequestRepository) GetByEmployeeAndDateRange(ctx context.Context, employeeID uuid.UUID, startDate, endDate string) ([]domain.LeaveRequest, error) {
	var requests []domain.LeaveRequest

	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, err
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).
		Where("employee_id = ? AND start_date <= ? AND end_date >= ? AND status != ?",
			employeeID, end, start, domain.StatusRejected).
		Order("start_date ASC").
		Find(&requests).Error; err != nil {
		return nil, err
	}

	return requests, nil
}