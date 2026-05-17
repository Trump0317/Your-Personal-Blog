package controller

import (
	"github.com/gin-gonic/gin"
	config "github.com/ypb/your-personal-blog/config/hub"
	"github.com/ypb/your-personal-blog/internal/hub/usercase"
	"github.com/ypb/your-personal-blog/internal/pkg/logger"
)

func NewRouter(ginEngine *gin.Engine, cfg *config.Config, l logger.Interface, file usercase.File, user usercase.User, admin usercase.Admin) {
	// 注册中间件
	ginEngine.Use(LoggerMiddleware(l))
	ginEngine.Use(RecoveryMiddleware(l))

	// 注册全局路由 /api
	apiGroup := ginEngine.Group("/api")

	// 注册 v1 路由 /api/v1
	apiV1Group := apiGroup.Group("/v1")
	{
		// 管理员路由 /api/v1/admin/...
		NewAdminRouter(apiV1Group, admin)

		// 业务路由 /api/v1/file/... 与 /api/v1/user/...
		NewUserRouter(apiV1Group, l, file, user)
	}
}
