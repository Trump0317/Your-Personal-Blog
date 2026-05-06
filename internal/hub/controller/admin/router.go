package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/ypb/your-personal-blog/internal/hub/controller/middleware"
	"github.com/ypb/your-personal-blog/internal/hub/usercase"
)

func NewAdminRouter(group *gin.RouterGroup, adminUC usercase.Admin) {
	handler := NewAdminController(adminUC)

	adminGroup := group.Group("/admin")
	adminGroup.Use(middleware.AdminAuthMiddleware(adminUC))
	{
		adminGroup.GET("/stats", handler.GetStats)
	}
}
