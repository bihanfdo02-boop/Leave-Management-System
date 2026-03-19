package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware handles CORS headers
func CORSMiddleware(allowedOrigins string) gin.HandlerFunc {
	origins := strings.Split(allowedOrigins, ",")
	originMap := make(map[string]bool)
	for _, origin := range origins {
		originMap[strings.TrimSpace(origin)] = true
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if originMap[origin] {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}