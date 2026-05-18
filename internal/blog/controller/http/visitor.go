package http

import (
	"net/http"

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
		c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Message: "参数错误"})
		return
	}

	post, err := vc.postUC.Get(c.Request.Context(), idOrSlug)
	if err != nil {
		c.JSON(http.StatusNotFound, Response{Code: http.StatusNotFound, Message: "文章不存在"})
		return
	}

	tags := make([]TagBenefit, 0)
	for _, t := range post.Tags {
		tags = append(tags, TagBenefit{
			ID:   t.ID,
			Name: t.Name,
			Slug: t.Slug,
		})
	}

	var cat *CategoryBenefit
	if post.Category != nil {
		cat = &CategoryBenefit{
			ID:   post.Category.ID,
			Name: post.Category.Name,
			Slug: post.Category.Slug,
		}
	}

	res := PostDetailResponse{
		ID:          post.ID,
		Title:       post.Title,
		Slug:        post.Slug,
		Content:     post.Content,
		HTMLContent: post.HTMLContent,
		Summary:     post.Summary,
		ViewCount:   post.ViewCount,
		Category:    cat,
		Tags:        tags,
		PublishedAt: post.PublishedAt,
	}

	c.JSON(http.StatusOK, Response{Code: 0, Message: "success", Data: res})
}

// ListPosts API：返回文章列表数据 (JSON)
func (vc *visitorController) ListPosts(c *gin.Context) {
	var req PostListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Message: "参数错误"})
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	status := model.PostPublished
	in := &usecase.PostListInput{
		Page:     req.Page,
		PageSize: req.PageSize,
		Status:   &status,
		Query:    req.Query,
	}

	if req.CategoryID != "" {
		in.CategoryID = &req.CategoryID
	}
	if req.TagID != "" {
		in.TagID = &req.TagID
	}

	output, err := vc.postUC.List(c.Request.Context(), in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: "获取文章失败"})
		return
	}

	items := make([]*PostItemResponse, 0)
	for _, p := range output.Posts {
		tags := make([]TagBenefit, 0)
		for _, t := range p.Tags {
			tags = append(tags, TagBenefit{
				ID:   t.ID,
				Name: t.Name,
				Slug: t.Slug,
			})
		}

		var cat *CategoryBenefit
		if p.Category != nil {
			cat = &CategoryBenefit{
				ID:   p.Category.ID,
				Name: p.Category.Name,
				Slug: p.Category.Slug,
			}
		}

		items = append(items, &PostItemResponse{
			ID:          p.ID,
			Title:       p.Title,
			Slug:        p.Slug,
			Summary:     p.Summary,
			ViewCount:   p.ViewCount,
			Category:    cat,
			Tags:        tags,
			PublishedAt: p.PublishedAt,
		})
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data: PostListResponse{
			Total: output.Total,
			Items: items,
		},
	})
}

// ListTags API：返回所有标签供展示
func (vc *visitorController) ListTags(c *gin.Context) {
	tags, err := vc.tagUC.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: "获取标签失败"})
		return
	}

	res := make([]TagBenefit, 0)
	for _, t := range tags {
		res = append(res, TagBenefit{
			ID:   t.ID,
			Name: t.Name,
			Slug: t.Slug,
		})
	}

	c.JSON(http.StatusOK, Response{Code: 0, Message: "success", Data: res})
}

// ListCategories API：返回所有分类供展示
func (vc *visitorController) ListCategories(c *gin.Context) {
	categories, err := vc.categoryUC.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: "获取分类失败"})
		return
	}

	res := make([]CategoryBenefit, 0)
	for _, cat := range categories {
		res = append(res, CategoryBenefit{
			ID:   cat.ID,
			Name: cat.Name,
			Slug: cat.Slug,
		})
	}

	c.JSON(http.StatusOK, Response{Code: 0, Message: "success", Data: res})
}
