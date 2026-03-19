package repository

import (
	"context"
	"errors"

	"github.com/KalinduBihan/leave-management-api/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PostgresUserRepository implements UserRepository interface
type PostgresUserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &PostgresUserRepository{db: db}
}

// Create creates a new user
func (r *PostgresUserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// GetByID retrieves a user by ID
func (r *PostgresUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// Update updates a user
func (r *PostgresUserRepository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// Delete soft deletes a user by setting is_active to false
func (r *PostgresUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).Update("is_active", false).Error
}

// List retrieves all users with pagination
func (r *PostgresUserRepository) List(ctx context.Context, page, pageSize int) ([]domain.User, int64, error) {
	var users []domain.User
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&domain.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}