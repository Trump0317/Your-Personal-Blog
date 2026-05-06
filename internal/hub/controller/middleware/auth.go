package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ypb/your-personal-blog/internal/hub/usercase"
)

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
