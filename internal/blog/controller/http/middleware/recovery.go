package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/ypb/your-personal-blog/internal/pkg/logger"
)

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
