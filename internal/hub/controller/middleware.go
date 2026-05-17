package controller

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ypb/your-personal-blog/internal/hub/usercase"
	"github.com/ypb/your-personal-blog/internal/pkg/logger"
)

// AdminAuthMiddleware 管理员专属鉴权中间件
func AdminAuthMiddleware(adminUC usercase.Admin) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-Admin-API-Key")
		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "X-Admin-API-Key header is required"})
			return
		}

		c.Next()
	}
}

const ContextUserKey = "user_info"

// AuthMiddleware 基于 APIKey 的鉴权中间件
func AuthMiddleware(user usercase.User) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "X-API-Key header is required"})
			return
		}

		user, err := user.GetByAPIKey(c.Request.Context(), apiKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid API Key"})
			return
		}

		// 将用户信息存入 context，方便后续逻辑（如配额检查）使用
		c.Set(ContextUserKey, user)
		c.Next()
	}
}

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

// RecoveryMiddleware 恢复中间件
func RecoveryMiddleware(l logger.Interface) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				l.Error(fmt.Sprintf("panic recovered: %v\nstack: %s", err, string(debug.Stack())))
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
