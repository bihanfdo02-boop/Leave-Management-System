package repository

import (
	"context"
	"errors"

	"github.com/KalinduBihan/leave-management-api/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PostgresDepartmentRepository implements DepartmentRepository interface
type PostgresDepartmentRepository struct {
	db *gorm.DB
}

// NewDepartmentRepository creates a new department repository
func NewDepartmentRepository(db *gorm.DB) DepartmentRepository {
	return &PostgresDepartmentRepository{db: db}
}

// Create creates a new department
func (r *PostgresDepartmentRepository) Create(ctx context.Context, department *domain.Department) (*domain.Department, error) {
	if err := r.db.WithContext(ctx).Create(department).Error; err != nil {
		return nil, err
	}
	return department, nil
}

// GetByID retrieves a department by ID with employees
func (r *PostgresDepartmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Department, error) {
	var department domain.Department
	if err := r.db.WithContext(ctx).
		Preload("Employees", "is_active = ?", true).
		Where("id = ?", id).
		First(&department).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrDepartmentNotFound
		}
		return nil, err
	}
	return &department, nil
}

// GetByName retrieves a department by name
func (r *PostgresDepartmentRepository) GetByName(ctx context.Context, name string) (*domain.Department, error) {
	var department domain.Department
	if err := r.db.WithContext(ctx).
		Where("name = ?", name).
		First(&department).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrDepartmentNotFound
		}
		return nil, err
	}
	return &department, nil
}

// Update updates a department
func (r *PostgresDepartmentRepository) Update(ctx context.Context, department *domain.Department) (*domain.Department, error) {
	if err := r.db.WithContext(ctx).Save(department).Error; err != nil {
		return nil, err
	}
	return department, nil
}

// Delete soft deletes a department
func (r *PostgresDepartmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&domain.Department{}).Where("id = ?", id).Update("is_active", false).Error
}

// List retrieves all departments with pagination
func (r *PostgresDepartmentRepository) List(ctx context.Context, page, pageSize int) ([]domain.Department, int64, error) {
	var departments []domain.Department
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&domain.Department{}).Where("is_active = ?", true).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).
		Preload("Employees", "is_active = ?", true).
		Where("is_active = ?", true).
		Offset(offset).
		Limit(pageSize).
		Order("name ASC").
		Find(&departments).Error; err != nil {
		return nil, 0, err
	}

	return departments, total, nil
}