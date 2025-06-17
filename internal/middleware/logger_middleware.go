package middleware

import (
	"time"

	"mock-server/pkg/logger"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		// Process request
		c.Next()

		// Log after request completes
		end := time.Now()
		latency := end.Sub(start)

		tmpLogger := logger.GetLogger()

		tmpLogger.InfoW("api executed", map[string]any{
			"request_id": c.GetHeader("X-Request-ID"),
			"user_id":    c.GetHeader("X-User-ID"),
			"status":     c.Writer.Status(),
			"method":     c.Request.Method,
			"path":       path,
			"query":      query,
			"ip":         c.ClientIP(),
			"user-agent": c.Request.UserAgent(),
			"latency":    latency,
			"errors":     c.Errors.ByType(gin.ErrorTypePrivate).String(),
		})

	}
}
