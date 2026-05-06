package controller

import (
	"github.com/gin-gonic/gin"
	config "github.com/ypb/your-personal-blog/config/hub"
	adminController "github.com/ypb/your-personal-blog/internal/hub/controller/admin"
	"github.com/ypb/your-personal-blog/internal/hub/controller/middleware"
	userController "github.com/ypb/your-personal-blog/internal/hub/controller/user"
	"github.com/ypb/your-personal-blog/internal/hub/usercase"
	"github.com/ypb/your-personal-blog/internal/pkg/logger"
)

func NewRouter(ginEngine *gin.Engine, cfg *config.Config, l logger.Interface, file usercase.File, user usercase.User, admin usercase.Admin) {
	// 注册中间件
	ginEngine.Use(middleware.LoggerMiddleware(l))
	ginEngine.Use(middleware.RecoveryMiddleware(l))

	// 注册全局路由 /api
	apiGroup := ginEngine.Group("/api")

	// 注册 v1 路由 /api/v1
	apiV1Group := apiGroup.Group("/v1")
	{
		// 管理员路由 /api/v1/admin/...
		adminController.NewAdminRouter(apiV1Group, admin)

		// 业务路由 /api/v1/file/...
		userController.NewFileRouter(apiV1Group, l, file, user)
		// 业务路由 /api/v1/user/...
		userController.NewUserRouter(apiV1Group, user)
	}
}
