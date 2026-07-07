package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 简单的请求日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		log.Printf("[%s] %s %s %v",
			c.Request.Method,
			c.Request.URL.Path,
			statusColor(c.Writer.Status()),
			time.Since(start),
		)
	}
}

func statusColor(code int) string {
	switch {
	case code >= 200 && code < 300:
		return "200"
	case code >= 300 && code < 400:
		return "3xx"
	case code >= 400 && code < 500:
		return "4xx"
	default:
		return "5xx"
	}
}
