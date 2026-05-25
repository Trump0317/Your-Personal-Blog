package http

import (
	"github.com/gin-gonic/gin"

	config "github.com/ypb/your-personal-blog/config/blog"
	"github.com/ypb/your-personal-blog/internal/blog/usecase"
	"github.com/ypb/your-personal-blog/internal/pkg/logger"
)

func NewRouter(ginEngine *gin.Engine, cfg *config.Config, l logger.Interface, postUC usecase.Post, tagUC usecase.Tag, categoryUC usecase.Category) {
	// 全局中间件
	ginEngine.Use(LoggerMiddleware(l))
	ginEngine.Use(RecoveryMiddleware(l))

	// 注册博客相关路由
	blogGroup := ginEngine.Group("/blog")

	// 访问端路由
	visitorGroup := blogGroup.Group("/visitor")
	NewVisitorRouter(visitorGroup, postUC, tagUC, categoryUC)

	// 管理端控制层初始化 (部分路由需要鉴权，部分不需要)
	adminCtrl := NewAdminController(postUC, categoryUC, tagUC, cfg)

	// 无需鉴权的管理端路由
	blogGroup.GET("/admin/login", adminCtrl.LoginGUI)
	blogGroup.POST("/admin/login", adminCtrl.Login)

	// 管理端路由 (需要鉴权)
	adminGroup := blogGroup.Group("/admin")
	adminGroup.Use(AuthMiddleware(cfg.Admin.Username, cfg.Admin.Password))
	NewAdminRouter(adminGroup, adminCtrl)
}

func NewAdminRouter(group *gin.RouterGroup, admin *adminController) {
	// 管理端入口界面
	group.GET("/gui", admin.GUI)

	// 文章管理
	group.GET("/posts", admin.ListPosts)
	group.GET("/posts/:id", admin.GetPost)
	group.POST("/posts", admin.CreatePost)
	group.PUT("/posts/:id", admin.UpdatePost)
	group.DELETE("/posts/:id", admin.DeletePost)

	// 元数据列表 (用于下拉框/选择器)
	group.GET("/categories", admin.ListCategories)
	group.POST("/categories", admin.CreateCategory)
	group.PUT("/categories/:id", admin.UpdateCategory)
	group.DELETE("/categories/:id", admin.DeleteCategory)

	group.GET("/tags", admin.ListTags)
	group.POST("/tags", admin.CreateTag)
	group.DELETE("/tags/:id", admin.DeleteTag)
}
