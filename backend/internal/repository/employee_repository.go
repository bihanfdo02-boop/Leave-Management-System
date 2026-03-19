package repository

import (
	"context"
	"errors"

	"github.com/KalinduBihan/leave-management-api/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PostgresEmployeeRepository implements EmployeeRepository interface
type PostgresEmployeeRepository struct {
	db *gorm.DB
}

// NewEmployeeRepository creates a new employee repository
func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &PostgresEmployeeRepository{db: db}
}

// Create creates a new employee
func (r *PostgresEmployeeRepository) Create(ctx context.Context, employee *domain.Employee) (*domain.Employee, error) {
	if err := r.db.WithContext(ctx).Create(employee).Error; err != nil {
		return nil, err
	}
	return employee, nil
}

// GetByID retrieves an employee by ID with relationships
func (r *PostgresEmployeeRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Employee, error) {
	var employee domain.Employee
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Department").
		Preload("Manager").
		Preload("LeaveBalances").
		Where("id = ?", id).
		First(&employee).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrEmployeeNotFound
		}
		return nil, err
	}
	return &employee, nil
}

// GetByUserID retrieves an employee by user ID
func (r *PostgresEmployeeRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Employee, error) {
	var employee domain.Employee
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Department").
		Where("user_id = ?", userID).
		First(&employee).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrEmployeeNotFound
		}
		return nil, err
	}
	return &employee, nil
}

// GetByEmail retrieves an employee by email
func (r *PostgresEmployeeRepository) GetByEmail(ctx context.Context, email string) (*domain.Employee, error) {
	var employee domain.Employee
	if err := r.db.WithContext(ctx).
		Preload("User").
		Where("email = ?", email).
		First(&employee).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrEmployeeNotFound
		}
		return nil, err
	}
	return &employee, nil
}

// Update updates an employee
func (r *PostgresEmployeeRepository) Update(ctx context.Context, employee *domain.Employee) (*domain.Employee, error) {
	if err := r.db.WithContext(ctx).Save(employee).Error; err != nil {
		return nil, err
	}
	return employee, nil
}

// Delete soft deletes an employee
func (r *PostgresEmployeeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&domain.Employee{}).Where("id = ?", id).Update("is_active", false).Error
}

// List retrieves all employees with pagination
func (r *PostgresEmployeeRepository) List(ctx context.Context, page, pageSize int) ([]domain.Employee, int64, error) {
	var employees []domain.Employee
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&domain.Employee{}).Where("is_active = ?", true).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Department").
		Preload("Manager").
		Where("is_active = ?", true).
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&employees).Error; err != nil {
		return nil, 0, err
	}

	return employees, total, nil
}

// ListByDepartment retrieves employees by department ID
func (r *PostgresEmployeeRepository) ListByDepartment(ctx context.Context, deptID uuid.UUID) ([]domain.Employee, error) {
	var employees []domain.Employee
	if err := r.db.WithContext(ctx).
		Preload("User").
		Where("department_id = ? AND is_active = ?", deptID, true).
		Order("first_name ASC, last_name ASC").
		Find(&employees).Error; err != nil {
		return nil, err
	}
	return employees, nil
}

// ListByManager retrieves employees by manager ID
func (r *PostgresEmployeeRepository) ListByManager(ctx context.Context, managerID uuid.UUID) ([]domain.Employee, error) {
	var employees []domain.Employee
	if err := r.db.WithContext(ctx).
		Preload("User").
		Where("manager_id = ? AND is_active = ?", managerID, true).
		Order("first_name ASC, last_name ASC").
		Find(&employees).Error; err != nil {
		return nil, err
	}
	return employees, nil
}