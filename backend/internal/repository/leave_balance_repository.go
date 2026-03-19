package repository

import (
	"context"
	"errors"
	"time"

	"github.com/KalinduBihan/leave-management-api/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PostgresLeaveBalanceRepository implements LeaveBalanceRepository interface
type PostgresLeaveBalanceRepository struct {
	db *gorm.DB
}

// NewLeaveBalanceRepository creates a new leave balance repository
func NewLeaveBalanceRepository(db *gorm.DB) LeaveBalanceRepository {
	return &PostgresLeaveBalanceRepository{db: db}
}

// Create creates a new leave balance
func (r *PostgresLeaveBalanceRepository) Create(ctx context.Context, balance *domain.LeaveBalance) (*domain.LeaveBalance, error) {
	if err := r.db.WithContext(ctx).Create(balance).Error; err != nil {
		return nil, err
	}
	return balance, nil
}

// GetByID retrieves a leave balance by ID
func (r *PostgresLeaveBalanceRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.LeaveBalance, error) {
	var balance domain.LeaveBalance
	if err := r.db.WithContext(ctx).
		Preload("Employee").
		Preload("LeaveType").
		Where("id = ?", id).
		First(&balance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrLeaveBalanceNotFound
		}
		return nil, err
	}
	return &balance, nil
}

// GetByEmployeeAndType retrieves balance for an employee and leave type in a year
func (r *PostgresLeaveBalanceRepository) GetByEmployeeAndType(ctx context.Context, employeeID, leaveTypeID uuid.UUID, year int) (*domain.LeaveBalance, error) {
	var balance domain.LeaveBalance
	if err := r.db.WithContext(ctx).
		Preload("Employee").
		Preload("LeaveType").
		Where("employee_id = ? AND leave_type_id = ? AND year = ?", employeeID, leaveTypeID, year).
		First(&balance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrLeaveBalanceNotFound
		}
		return nil, err
	}
	return &balance, nil
}

// GetByEmployee retrieves all balances for an employee in a year
func (r *PostgresLeaveBalanceRepository) GetByEmployee(ctx context.Context, employeeID uuid.UUID, year int) ([]domain.LeaveBalance, error) {
	var balances []domain.LeaveBalance
	if err := r.db.WithContext(ctx).
		Preload("LeaveType").
		Where("employee_id = ? AND year = ?", employeeID, year).
		Find(&balances).Error; err != nil {
		return nil, err
	}
	return balances, nil
}

// Update updates a leave balance
func (r *PostgresLeaveBalanceRepository) Update(ctx context.Context, balance *domain.LeaveBalance) (*domain.LeaveBalance, error) {
	if err := r.db.WithContext(ctx).Save(balance).Error; err != nil {
		return nil, err
	}
	return balance, nil
}

// Delete deletes a leave balance
func (r *PostgresLeaveBalanceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.LeaveBalance{}, "id = ?", id).Error
}

// List retrieves all leave balances with pagination
func (r *PostgresLeaveBalanceRepository) List(ctx context.Context, page, pageSize int) ([]domain.LeaveBalance, int64, error) {
	var balances []domain.LeaveBalance
	var total int64

	currentYear := time.Now().Year()

	// Get total count
	if err := r.db.WithContext(ctx).
		Model(&domain.LeaveBalance{}).
		Where("year = ?", currentYear).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).
		Preload("Employee").
		Preload("LeaveType").
		Where("year = ?", currentYear).
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&balances).Error; err != nil {
		return nil, 0, err
	}

	return balances, total, nil
}

// UpdateUsedDays updates the used days for a balance
func (r *PostgresLeaveBalanceRepository) UpdateUsedDays(ctx context.Context, balanceID uuid.UUID, daysToAdd float32) error {
	return r.db.WithContext(ctx).
		Model(&domain.LeaveBalance{}).
		Where("id = ?", balanceID).
		Update("used_days", gorm.Expr("used_days + ?", daysToAdd)).
		Error
}