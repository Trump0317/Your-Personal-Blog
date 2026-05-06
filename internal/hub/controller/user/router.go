package user

import (
	"github.com/gin-gonic/gin"
	"github.com/ypb/your-personal-blog/internal/hub/controller/middleware"
	"github.com/ypb/your-personal-blog/internal/hub/usercase"
	"github.com/ypb/your-personal-blog/internal/pkg/logger"
)

func NewFileRouter(group *gin.RouterGroup, l logger.Interface, file usercase.File, user usercase.User) {
	fileHandler := NewFileController(file, user)

	// V1 API 规范
	fileGroup := group.Group("/file")
	{
		// 1. 上传文件 - API Key 鉴权
		fileGroup.POST("/upload", middleware.AuthMiddleware(user), fileHandler.Upload)

		// 2. 获取单个文件 - 公开/Token
		fileGroup.GET("/:id", fileHandler.GetByID)

		// 3. 获取用户全部文件 - API Key 鉴权
		fileGroup.GET("/all", middleware.AuthMiddleware(user), fileHandler.ListAll)

		// 4. 删除文件 - API Key 鉴权
		fileGroup.DELETE("/:id", middleware.AuthMiddleware(user), fileHandler.Delete)
	}
}

func NewUserRouter(group *gin.RouterGroup, user usercase.User) {
	userHandler := NewUserController(user)

	userGroup := group.Group("/user")
	{
		// 注册接口 - 公开
		userGroup.POST("/register", userHandler.Register)
	}
}
