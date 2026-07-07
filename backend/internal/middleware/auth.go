package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// BasicAuth 简单的 HTTP Basic 认证中间件
func BasicAuth(username, password string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, pwd, ok := c.Request.BasicAuth()
		if !ok || user != username || pwd != password {
			c.Header("WWW-Authenticate", `Basic realm="Blog Admin"`)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}
