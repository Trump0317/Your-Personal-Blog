package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ypb/your-personal-blog/internal/hub/usercase"
)

type AdminController struct {
	adminUC usercase.Admin
}

func NewAdminController(adminUC usercase.Admin) *AdminController {
	return &AdminController{
		adminUC: adminUC,
	}
}

func NewAdminRouter(group *gin.RouterGroup, adminUC usercase.Admin) {
	handler := NewAdminController(adminUC)

	adminGroup := group.Group("/admin")
	adminGroup.Use(AdminAuthMiddleware(adminUC))
	{
		adminGroup.GET("/stats", handler.GetStats)
	}
}

// GetStats 获取系统全量统计
func (a *AdminController) GetStats(c *gin.Context) {
	stats, err := a.adminUC.GetSystemStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
