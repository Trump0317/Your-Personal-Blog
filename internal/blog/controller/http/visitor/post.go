package visitor

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ypb/your-personal-blog/internal/blog/model"
	"github.com/ypb/your-personal-blog/internal/blog/usecase"
	"github.com/ypb/your-personal-blog/internal/blog/usecase/input"
)

type visitorController struct {
	postUC usecase.Post
}

func newVisitorController(postUC usecase.Post) *visitorController {
	return &visitorController{postUC: postUC}
}

// Index 返回博客首页 HTML
func (vc *visitorController) Index(c *gin.Context) {
	// 直接返回静态 HTML 文件
	c.File("web/blog/index.html")
}

// ListPosts API：返回文章列表数据 (JSON)
func (vc *visitorController) ListPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	status := model.PostPublished
	in := &input.PostList{
		Page:     page,
		PageSize: pageSize,
		Status:   &status, // 仅展示已发布的文章
	}

	posts, total, err := vc.postUC.List(c.Request.Context(), in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"total": total,
	})
}
