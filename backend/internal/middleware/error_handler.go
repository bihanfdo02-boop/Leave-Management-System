package middleware

import (
	"net/http"

	"github.com/KalinduBihan/leave-management-api/internal/domain"
	"github.com/KalinduBihan/leave-management-api/internal/dto"
	"github.com/gin-gonic/gin"
)

// ErrorHandlerMiddleware handles panics and errors
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromAppError(domain.ErrInternalServer))
				c.Abort()
			}
		}()
		c.Next()
	}
}