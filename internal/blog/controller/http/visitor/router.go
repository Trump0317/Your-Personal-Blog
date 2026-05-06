package visitor

import (
	"github.com/gin-gonic/gin"

	"github.com/ypb/your-personal-blog/internal/blog/usecase"
)

func NewVisitorRouter(group *gin.RouterGroup, postUC usecase.Post) {
	vc := newVisitorController(postUC)

	// 1. 页面路由：用户浏览器直接访问的内容
	group.GET("/index", vc.Index)

	// 2. 数据接口：前端页面通过 JavaScript 请求的内容
	group.GET("/posts", vc.ListPosts)
}
