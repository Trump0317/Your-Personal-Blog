package http

import (
	"github.com/gin-gonic/gin"

	config "github.com/ypb/your-personal-blog/config/blog"
	"github.com/ypb/your-personal-blog/internal/blog/controller/http/visitor"
	"github.com/ypb/your-personal-blog/internal/blog/usecase"
	"github.com/ypb/your-personal-blog/internal/pkg/logger"
)

func NewRouter(ginEngine *gin.Engine, cfg *config.Config, l logger.Interface, postUC usecase.Post) {
	// 注册博客相关路由
	blogGroup := ginEngine.Group("/blog")

	// 访问端路由
	visitorGroup := blogGroup.Group("/visitor")
	visitor.NewVisitorRouter(visitorGroup, postUC)

	// 管理端路由

}
