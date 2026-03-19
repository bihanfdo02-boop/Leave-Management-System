package repository

import (
	"context"
	"errors"

	"github.com/KalinduBihan/leave-management-api/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PostgresAuditLogRepository implements AuditLogRepository interface
type PostgresAuditLogRepository struct {
	db *gorm.DB
}

// NewAuditLogRepository creates a new audit log repository
func NewAuditLogRepository(db *gorm.DB) AuditLogRepository {
	return &PostgresAuditLogRepository{db: db}
}

// Create creates a new audit log
func (r *PostgresAuditLogRepository) Create(ctx context.Context, log *domain.AuditLog) (*domain.AuditLog, error) {
	if err := r.db.WithContext(ctx).Create(log).Error; err != nil {
		return nil, err
	}
	return log, nil
}

// GetByID retrieves an audit log by ID
func (r *PostgresAuditLogRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.AuditLog, error) {
	var log domain.AuditLog
	if err := r.db.WithContext(ctx).
		Preload("User").
		Where("id = ?", id).
		First(&log).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.NewAppError(domain.ErrCodeNotFound, "audit log not found", 404, nil)
		}
		return nil, err
	}
	return &log, nil
}

// GetByUser retrieves audit logs by user ID
func (r *PostgresAuditLogRepository) GetByUser(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]domain.AuditLog, int64, error) {
	var logs []domain.AuditLog
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).
		Model(&domain.AuditLog{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).
		Preload("User").
		Where("user_id = ?", userID).
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetByEntity retrieves audit logs by entity type and ID
func (r *PostgresAuditLogRepository) GetByEntity(ctx context.Context, entityType string, entityID uuid.UUID) ([]domain.AuditLog, error) {
	var logs []domain.AuditLog
	if err := r.db.WithContext(ctx).
		Preload("User").
		Where("entity_type = ? AND entity_id = ?", entityType, entityID).
		Order("created_at DESC").
		Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

// List retrieves all audit logs with pagination
func (r *PostgresAuditLogRepository) List(ctx context.Context, page, pageSize int) ([]domain.AuditLog, int64, error) {
	var logs []domain.AuditLog
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&domain.AuditLog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).
		Preload("User").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}