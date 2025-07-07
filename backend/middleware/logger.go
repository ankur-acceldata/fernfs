package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger returns a gin middleware for request logging
func Logger(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		if len(c.Errors) > 0 {
			// Log errors
			for _, e := range c.Errors.Errors() {
				log.Error("request error",
					zap.String("path", path),
					zap.String("query", query),
					zap.String("error", e),
				)
			}
		} else {
			log.Info("request completed",
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.Int("status", c.Writer.Status()),
				zap.Duration("latency", latency),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
			)
		}
	}
} 