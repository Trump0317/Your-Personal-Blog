package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ypb/your-personal-blog/internal/pkg/logger"
)

// LoggerMiddleware 日志中间件
func LoggerMiddleware(l logger.Interface) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				l.Error(e)
			}
		} else {
			l.Info("request",
				"status", c.Writer.Status(),
				"method", c.Request.Method,
				"path", path,
				"query", query,
				"ip", c.ClientIP(),
				"latency", latency,
			)
		}
	}
}
