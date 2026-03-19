package dto

import (
	"time"

	"github.com/KalinduBihan/leave-management-api/internal/domain"
	"github.com/google/uuid"
)

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required,min=2"`
	LastName  string `json:"last_name" binding:"required,min=2"`
	Phone     string `json:"phone" binding:"omitempty,min=10,max=20"`
}

// LoginRequest represents a user login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents a successful login response
type LoginResponse struct {
	AccessToken string        `json:"access_token"`
	ExpiresIn   int64         `json:"expires_in"`
	User        *UserResponse `json:"user"`
}

// UserResponse represents a user in responses
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	LastLogin *time.Time `json:"last_login,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToUserResponse converts a User domain model to UserResponse
func ToUserResponse(user *domain.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Role:      string(user.Role),
		IsActive:  user.IsActive,
		LastLogin: user.LastLogin,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}