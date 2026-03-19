package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware logs request details
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		duration := time.Since(startTime)
		statusCode := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.RequestURI

		log.Printf("[%s] %s %s - %d - %v",
			time.Now().Format("2006-01-02 15:04:05"),
			method,
			path,
			statusCode,
			duration,
		)
	}
}