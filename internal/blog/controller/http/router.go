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
	group.POST("/posts", admin.CreatePost)
	group.PUT("/posts/:id", admin.UpdatePost)
	group.DELETE("/posts/:id", admin.DeletePost)

	// 元数据列表 (用于下拉框/选择器)
	group.GET("/categories", admin.ListCategories)
	group.GET("/tags", admin.ListTags)
}

/*
已暴露的接口清单:

I. 访客端 API (/blog/visitor)
---------------------------------------
1. 页面展示:
   - GET  /blog/visitor/index          : 博客首页 HTML
   - GET  /blog/visitor/admin          : 管理后台 HTML

2. 文章接口:
   - GET  /blog/visitor/posts          : 获取已发布文章列表 (分页/过滤/搜索)
   - GET  /blog/visitor/post/:idOrSlug : 获取指定文章详情 (支持 ID 或 Slug)

3. 元数据接口:
   - GET  /blog/visitor/categories     : 获取全量分类
   - GET  /blog/visitor/tags           : 获取全量标签


II. 管理端 API (/blog/admin)
---------------------------------------
1. 文章管理:
   - GET    /blog/admin/posts          : 获取全量文章列表 (含草稿/已发布)
   - POST   /blog/admin/posts          : 创建新文章 (默认强制为草稿)
   - PUT    /blog/admin/posts/:id      : 更新文章信息 (支持修改标题/Slug/内容/分类/状态/标签)
   - DELETE /blog/admin/posts/:id      : 删除文章

2. 辅助接口 (用于编辑器配置):
   - GET    /blog/admin/categories     : 获取所有分类列表
   - GET    /blog/admin/tags           : 获取所有标签列表
*/
