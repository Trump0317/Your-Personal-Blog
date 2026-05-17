package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ypb/your-personal-blog/internal/blog/model"
	"github.com/ypb/your-personal-blog/internal/blog/usecase"
)

type visitorController struct {
	postUC     usecase.Post
	tagUC      usecase.Tag
	categoryUC usecase.Category
}

func NewVisitorRouter(group *gin.RouterGroup, postUC usecase.Post, tagUC usecase.Tag, categoryUC usecase.Category) {
	vc := &visitorController{
		postUC:     postUC,
		tagUC:      tagUC,
		categoryUC: categoryUC,
	}

	// 1. 页面路由：用户浏览器直接访问的内容
	group.GET("/index", vc.Index)
	group.GET("/post/:idOrSlug/page", vc.PostPage)
	group.GET("/categories/page", vc.CategoriesPage)
	group.GET("/tags/page", vc.TagsPage)
	group.GET("/archives/page", vc.ArchivesPage)

	// 2. 数据接口：前端页面通过 JavaScript 请求的内容
	group.GET("/posts", vc.ListPosts)
	group.GET("/post/:idOrSlug", vc.GetPost)
	group.GET("/tags", vc.ListTags)
	group.GET("/categories", vc.ListCategories)
}

// Index 返回博客首页 HTML
func (vc *visitorController) Index(c *gin.Context) {
	// 直接返回静态 HTML 文件
	c.File("web/blog/index.html")
}

// PostPage 返回文章详情 HTML
func (vc *visitorController) PostPage(c *gin.Context) {
	c.File("web/blog/post.html")
}

// CategoriesPage 返回分类列表页 HTML
func (vc *visitorController) CategoriesPage(c *gin.Context) {
	c.File("web/blog/categories.html")
}

// TagsPage 返回标签列表页 HTML
func (vc *visitorController) TagsPage(c *gin.Context) {
	c.File("web/blog/tags.html")
}

// ArchivesPage 返回文章归档页 HTML
func (vc *visitorController) ArchivesPage(c *gin.Context) {
	c.File("web/blog/archives.html")
}

// GetPost API：返回单篇文章详情
func (vc *visitorController) GetPost(c *gin.Context) {
	idOrSlug := c.Param("idOrSlug")
	if idOrSlug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	post, err := vc.postUC.Get(c.Request.Context(), idOrSlug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 在 Controller 层聚合数据
	var category *model.Category
	if post.CategoryID != "" {
		category, _ = vc.categoryUC.Get(c.Request.Context(), post.CategoryID)
	}

	var tags []*model.Tag
	if len(post.TagIDs) > 0 {
		tags, _ = vc.tagUC.ListByIDs(c.Request.Context(), post.TagIDs)
	}

	output := &usecase.PostDetailOutput{
		ID:          post.ID,
		Title:       post.Title,
		Slug:        post.Slug,
		Content:     post.Content,
		HTMLContent: post.HTMLContent,
		Summary:     post.Summary,
		Status:      post.Status,
		ViewCount:   post.ViewCount,
		CategoryID:  post.CategoryID,
		Category:    category,
		Tags:        tagsToValueSlice(tags),
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		PublishedAt: post.PublishedAt,
	}

	c.JSON(http.StatusOK, output)
}

// ListPosts API：返回文章列表数据 (JSON)
func (vc *visitorController) ListPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	status := model.PostPublished
	in := &usecase.PostListInput{
		Page:     page,
		PageSize: pageSize,
		Status:   &status, // 仅展示已发布的文章
	}

	// 增加按分类或标签过滤的支持
	catID := c.Query("category_id")
	if catID != "" {
		in.CategoryID = &catID
	}
	tagID := c.Query("tag_id")
	if tagID != "" {
		in.TagID = &tagID
	}
	query := c.Query("q")
	if query != "" {
		in.Query = query
	}

	posts, total, err := vc.postUC.List(c.Request.Context(), in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章失败"})
		return
	}

	// 在 Controller 层聚合数据
	items := make([]*usecase.PostItem, 0, len(posts))
	for _, p := range posts {
		var category *model.Category
		if p.CategoryID != "" {
			category, _ = vc.categoryUC.Get(c.Request.Context(), p.CategoryID)
		}

		var tags []*model.Tag
		if len(p.TagIDs) > 0 {
			tags, _ = vc.tagUC.ListByIDs(c.Request.Context(), p.TagIDs)
		}

		items = append(items, &usecase.PostItem{
			ID:          p.ID,
			Title:       p.Title,
			Slug:        p.Slug,
			Summary:     p.Summary,
			Status:      p.Status,
			ViewCount:   p.ViewCount,
			Category:    category,
			Tags:        tagsToValueSlice(tags),
			CreatedAt:   p.CreatedAt,
			PublishedAt: p.PublishedAt,
		})
	}

	c.JSON(http.StatusOK, usecase.PostListOutput{
		Posts: items,
		Total: total,
	})
}

// 辅助函数：将对象切片转为模型中的非指针切片
func tagsToValueSlice(tags []*model.Tag) []model.Tag {
	res := make([]model.Tag, 0, len(tags))
	for _, t := range tags {
		if t != nil {
			res = append(res, *t)
		}
	}
	return res
}

// ListTags API：返回所有标签供展示
func (vc *visitorController) ListTags(c *gin.Context) {
	tags, err := vc.tagUC.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取标签失败"})
		return
	}
	c.JSON(http.StatusOK, tags)
}

// ListCategories API：返回所有分类供展示
func (vc *visitorController) ListCategories(c *gin.Context) {
	categories, err := vc.categoryUC.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取分类失败"})
		return
	}
	c.JSON(http.StatusOK, categories)
}
