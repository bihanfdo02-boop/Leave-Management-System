package middleware

import (
	"net/http"
	"strings"

	"github.com/KalinduBihan/leave-management-api/internal/domain"
	"github.com/KalinduBihan/leave-management-api/internal/dto"
	"github.com/KalinduBihan/leave-management-api/internal/service"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT token and extracts user ID
func AuthMiddleware(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponseFromAppError(domain.ErrUnauthorized))
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponseFromAppError(domain.ErrUnauthorized))
			c.Abort()
			return
		}

		token := parts[1]

		// Validate token
		userID, err := authService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponseFromAppError(domain.ErrUnauthorized))
			c.Abort()
			return
		}

		// Store user ID in context
		c.Set("userID", userID.String())
		c.Next()
	}
}

// OptionalAuthMiddleware validates JWT token if provided
func OptionalAuthMiddleware(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		token := parts[1]
		userID, err := authService.ValidateToken(token)
		if err == nil {
			c.Set("userID", userID.String())
		}

		c.Next()
	}
}