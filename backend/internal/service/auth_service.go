package service

import (
	"context"
	"time"

	"github.com/KalinduBihan/leave-management-api/config"
	"github.com/KalinduBihan/leave-management-api/internal/domain"
	"github.com/KalinduBihan/leave-management-api/internal/dto"
	"github.com/KalinduBihan/leave-management-api/internal/repository"
	"github.com/KalinduBihan/leave-management-api/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// authService implements AuthService interface
type authService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

// Register registers a new user
func (s *authService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.LoginResponse, error) {
	// Validate input
	if !utils.IsValidEmail(req.Email) {
		return nil, domain.NewAppError(
			domain.ErrCodeValidation,
			"invalid email format",
			400,
			nil,
		)
	}

	if !utils.IsValidPassword(req.Password) {
		return nil, domain.NewAppError(
			domain.ErrCodeValidation,
			"password must be at least 8 characters",
			400,
			nil,
		)
	}

	// Check if email already exists
	email := utils.NormalizeEmail(req.Email)
	existingUser, _ := s.userRepo.GetByEmail(ctx, email)
	if existingUser != nil {
		return nil, domain.ErrDuplicateEmail
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, domain.NewAppError(
			domain.ErrCodeInternalServer,
			"failed to hash password",
			500,
			nil,
		)
	}

	// Create user
	user := &domain.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(passwordHash),
		Role:         domain.RoleEmployee,
		IsActive:     true,
	}

	user, err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, domain.NewAppError(
			domain.ErrCodeInternalServer,
			"failed to create user",
			500,
			err,
		)
	}

	// Generate token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	// Parse expiration duration
	expiration, _ := time.ParseDuration(s.cfg.JWT.Expiration)

	return &dto.LoginResponse{
		AccessToken: token,
		ExpiresIn:   int64(expiration.Seconds()),
		User:        dto.ToUserResponse(user),
	}, nil
}

// Login authenticates a user
func (s *authService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// Validate input
	email := utils.NormalizeEmail(req.Email)
	if !utils.IsValidEmail(email) {
		return nil, domain.ErrInvalidCredentials
	}

	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil || user == nil {
		return nil, domain.ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive {
		return nil, domain.NewAppError(
			domain.ErrCodeForbidden,
			"user account is inactive",
			403,
			nil,
		)
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	// Update last login
	now := time.Now()
	user.LastLogin = &now
	_, _ = s.userRepo.Update(ctx, user)

	// Generate token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	// Parse expiration duration
	expiration, _ := time.ParseDuration(s.cfg.JWT.Expiration)

	return &dto.LoginResponse{
		AccessToken: token,
		ExpiresIn:   int64(expiration.Seconds()),
		User:        dto.ToUserResponse(user),
	}, nil
}

// ValidateToken validates a JWT token and returns the user ID
func (s *authService) ValidateToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.NewAppError(
				domain.ErrCodeUnauthorized,
				"invalid token signing method",
				401,
				nil,
			)
		}
		return []byte(s.cfg.JWT.Secret), nil
	})

	if err != nil || !token.Valid {
		return uuid.UUID{}, domain.ErrUnauthorized
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return uuid.UUID{}, domain.ErrUnauthorized
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.UUID{}, domain.ErrUnauthorized
	}

	return userID, nil
}

// RefreshToken refreshes a JWT token
func (s *authService) RefreshToken(ctx context.Context, userID uuid.UUID) (string, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return "", domain.ErrUnauthorized
	}

	if !user.IsActive {
		return "", domain.NewAppError(
			domain.ErrCodeForbidden,
			"user account is inactive",
			403,
			nil,
		)
	}

	return s.generateToken(user)
}

// generateToken generates a JWT token
func (s *authService) generateToken(user *domain.User) (string, error) {
	expiration, err := time.ParseDuration(s.cfg.JWT.Expiration)
	if err != nil {
		expiration = 24 * time.Hour
	}

	claims := &jwt.RegisteredClaims{
		Subject:   user.ID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		return "", domain.NewAppError(
			domain.ErrCodeInternalServer,
			"failed to generate token",
			500,
			err,
		)
	}

	return tokenString, nil
}