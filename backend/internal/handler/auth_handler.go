package handler

import (
	"net/http"

	"github.com/KalinduBihan/leave-management-api/internal/domain"
	"github.com/KalinduBihan/leave-management-api/internal/dto"
	"github.com/KalinduBihan/leave-management-api/internal/service"
	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register registers a new user
// @Summary Register new user
// @Description Register a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register request"
// @Success 201 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseFromAppError(
			domain.NewAppError(domain.ErrCodeValidation, err.Error(), 400, nil),
		))
		return
	}

	response, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			c.JSON(appErr.Status, dto.ErrorResponseFromAppError(appErr))
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromAppError(domain.ErrInternalServer))
		}
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse("User registered successfully", response))
}

// Login authenticates a user
// @Summary Login
// @Description Authenticate user and get access token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login request"
// @Success 200 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseFromAppError(
			domain.NewAppError(domain.ErrCodeValidation, err.Error(), 400, nil),
		))
		return
	}

	response, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			c.JSON(appErr.Status, dto.ErrorResponseFromAppError(appErr))
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromAppError(domain.ErrInternalServer))
		}
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Login successful", response))
}

// RefreshToken refreshes the access token
// @Summary Refresh Token
// @Description Refresh access token using current token
// @Tags Auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponseFromAppError(domain.ErrUnauthorized))
		return
	}

	userID, err := parseUserID(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponseFromAppError(domain.ErrUnauthorized))
		return
	}

	token, err := h.authService.RefreshToken(c.Request.Context(), userID)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			c.JSON(appErr.Status, dto.ErrorResponseFromAppError(appErr))
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromAppError(domain.ErrInternalServer))
		}
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Token refreshed successfully", map[string]string{
		"access_token": token,
	}))
}