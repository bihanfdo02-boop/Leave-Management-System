package repository

import (
	"context"
	"errors"

	"github.com/KalinduBihan/leave-management-api/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PostgresLeaveTypeRepository implements LeaveTypeRepository interface
type PostgresLeaveTypeRepository struct {
	db *gorm.DB
}

// NewLeaveTypeRepository creates a new leave type repository
func NewLeaveTypeRepository(db *gorm.DB) LeaveTypeRepository {
	return &PostgresLeaveTypeRepository{db: db}
}

// Create creates a new leave type
func (r *PostgresLeaveTypeRepository) Create(ctx context.Context, leaveType *domain.LeaveType) (*domain.LeaveType, error) {
	if err := r.db.WithContext(ctx).Create(leaveType).Error; err != nil {
		return nil, err
	}
	return leaveType, nil
}

// GetByID retrieves a leave type by ID
func (r *PostgresLeaveTypeRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.LeaveType, error) {
	var leaveType domain.LeaveType
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&leaveType).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrLeaveTypeNotFound
		}
		return nil, err
	}
	return &leaveType, nil
}

// GetByName retrieves a leave type by name
func (r *PostgresLeaveTypeRepository) GetByName(ctx context.Context, name string) (*domain.LeaveType, error) {
	var leaveType domain.LeaveType
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&leaveType).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrLeaveTypeNotFound
		}
		return nil, err
	}
	return &leaveType, nil
}

// Update updates a leave type
func (r *PostgresLeaveTypeRepository) Update(ctx context.Context, leaveType *domain.LeaveType) (*domain.LeaveType, error) {
	if err := r.db.WithContext(ctx).Save(leaveType).Error; err != nil {
		return nil, err
	}
	return leaveType, nil
}

// Delete soft deletes a leave type
func (r *PostgresLeaveTypeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&domain.LeaveType{}).Where("id = ?", id).Update("is_active", false).Error
}

// List retrieves all leave types with pagination
func (r *PostgresLeaveTypeRepository) List(ctx context.Context, page, pageSize int) ([]domain.LeaveType, int64, error) {
	var leaveTypes []domain.LeaveType
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&domain.LeaveType{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(pageSize).
		Order("name ASC").
		Find(&leaveTypes).Error; err != nil {
		return nil, 0, err
	}

	return leaveTypes, total, nil
}

// ListActive retrieves all active leave types
func (r *PostgresLeaveTypeRepository) ListActive(ctx context.Context) ([]domain.LeaveType, error) {
	var leaveTypes []domain.LeaveType
	if err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("name ASC").
		Find(&leaveTypes).Error; err != nil {
		return nil, err
	}
	return leaveTypes, nil
}