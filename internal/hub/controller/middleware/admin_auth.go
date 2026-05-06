package middleware

import (
"net/http"

"github.com/gin-gonic/gin"
"github.com/ypb/your-personal-blog/internal/hub/usercase"
)

// AdminAuthMiddleware 管理员专属鉴权中间件
func AdminAuthMiddleware(adminUC usercase.Admin) gin.HandlerFunc {
return func(c *gin.Context) {
apiKey := c.GetHeader("X-Admin-API-Key")
if apiKey == "" {
c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "X-Admin-API-Key header is required"})
return
}

if !adminUC.VerifyAPIKey(apiKey) {
c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "invalid Admin API Key"})
return
}

c.Next()
}
}
