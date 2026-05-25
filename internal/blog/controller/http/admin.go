package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	config "github.com/ypb/your-personal-blog/config/blog"
	"github.com/ypb/your-personal-blog/internal/blog/model"
	"github.com/ypb/your-personal-blog/internal/blog/usecase"
)

type adminController struct {
	postUC     usecase.Post
	categoryUC usecase.Category
	tagUC      usecase.Tag
	cfg        *config.Config
}

func NewAdminController(p usecase.Post, c usecase.Category, t usecase.Tag, cfg *config.Config) *adminController {
	return &adminController{
		postUC:     p,
		categoryUC: c,
		tagUC:      t,
		cfg:        cfg,
	}
}

// ListPosts 列出所有文章（含草稿）
func (ctrl *adminController) ListPosts(c *gin.Context) {
	var input usecase.PostListInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	output, err := ctrl.postUC.ListBy(c.Request.Context(), &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	items := make([]*PostAdminItemResponse, 0)
	for _, p := range output.Posts {
		var cat *CategoryBrief
		if p.Category != nil {
			cat = &CategoryBrief{
				ID:   p.Category.ID,
				Name: p.Category.Name,
				Slug: p.Category.Slug,
			}
		}

		tags := make([]TagBrief, 0)
		for _, t := range p.Tags {
			tags = append(tags, TagBrief{
				ID:   t.ID,
				Name: t.Name,
				Slug: t.Slug,
			})
		}

		items = append(items, &PostAdminItemResponse{
			ID:          p.ID,
			Title:       p.Title,
			Slug:        p.Slug,
			Status:      int(p.Status),
			ViewCount:   p.ViewCount,
			Category:    cat,
			Tags:        tags,
			CreatedAt:   p.CreatedAt,
			PublishedAt: p.PublishedAt,
		})
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data: gin.H{
			"total": output.Total,
			"items": items,
		},
	})
}

// GUI 返回博客管理后台 HTML
func (ctrl *adminController) GUI(c *gin.Context) {
	c.File("web/blog/admin.html")
}

// LoginGUI 返回登录页面 HTML
func (ctrl *adminController) LoginGUI(c *gin.Context) {
	c.File("web/blog/login.html")
}

// Login 登录接口 (手动验证账号密码)
func (ctrl *adminController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err == nil {
		if req.Username == ctrl.cfg.Admin.Username && req.Password == ctrl.cfg.Admin.Password {
			c.JSON(http.StatusOK, Response{Code: 0, Message: "Login successful"})
			return
		}
	}

	user, pwd, ok := c.Request.BasicAuth()
	if ok && user == ctrl.cfg.Admin.Username && pwd == ctrl.cfg.Admin.Password {
		c.JSON(http.StatusOK, Response{Code: 0, Message: "Login successful"})
		return
	}

	c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
	c.JSON(http.StatusUnauthorized, Response{Code: http.StatusUnauthorized, Message: "Unauthorized"})
}

// CreatePost 创建文章
func (ctrl *adminController) CreatePost(c *gin.Context) {
	var body CreatePostRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	input := usecase.PostCreateInput{
		Title:      body.Title,
		Slug:       body.Slug,
		Content:    body.Content,
		Summary:    body.Summary,
		CategoryID: body.CategoryID,
		TagIDs:     body.TagIDs,
	}

	switch body.Status {
	case "published":
		input.Status = model.PostPublished
	default:
		input.Status = model.PostDraft
	}

	id, err := ctrl.postUC.Create(c.Request.Context(), &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, Response{Code: 0, Message: "success", Data: gin.H{"id": id}})
}

// GetPost 获取文章详情 (管理端支持获取草稿)
func (ctrl *adminController) GetPost(c *gin.Context) {
	id := c.Param("id")
	post, err := ctrl.postUC.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, Response{Code: http.StatusNotFound, Message: "post not found"})
		return
	}
	c.JSON(http.StatusOK, Response{Code: 0, Message: "success", Data: post})
}

// UpdatePost 更新文章
func (ctrl *adminController) UpdatePost(c *gin.Context) {
	var body UpdatePostRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	input := usecase.PostUpdateInput{
		ID:         c.Param("id"),
		Title:      body.Title,
		Slug:       body.Slug,
		Content:    body.Content,
		Summary:    body.Summary,
		CategoryID: body.CategoryID,
		TagIDs:     body.TagIDs,
	}

	if body.Status != nil {
		var s model.PostStatus
		switch *body.Status {
		case "published":
			s = model.PostPublished
		default:
			s = model.PostDraft
		}
		input.Status = &s
	}

	if err := ctrl.postUC.Update(c.Request.Context(), &input); err != nil {
		c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, Response{Code: 0, Message: "updated"})
}

// DeletePost 删除文章
func (ctrl *adminController) DeletePost(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.postUC.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, Response{Code: 0, Message: "deleted"})
}

// ListCategories 获取全量分类
func (ctrl *adminController) ListCategories(c *gin.Context) {
	cats, err := ctrl.categoryUC.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, Response{Code: 0, Message: "success", Data: cats})
}

// CreateCategory 创建分类
func (ctrl *adminController) CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	in := usecase.CategoryCreateInput{
		Name: req.Name,
		Slug: req.Slug,
	}
	id, err := ctrl.categoryUC.Create(c.Request.Context(), &in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, Response{Code: 0, Message: "success", Data: gin.H{"id": id}})
}

// UpdateCategory 更新分类
func (ctrl *adminController) UpdateCategory(c *gin.Context) {
	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	in := usecase.CategoryUpdateInput{
		ID:   c.Param("id"),
		Name: req.Name,
		Slug: req.Slug,
	}
	if err := ctrl.categoryUC.Update(c.Request.Context(), &in); err != nil {
		c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, Response{Code: 0, Message: "success"})
}

// DeleteCategory 删除分类
func (ctrl *adminController) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.categoryUC.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, Response{Code: 0, Message: "success"})
}

// ListTags 获取全量标签
func (ctrl *adminController) ListTags(c *gin.Context) {
	tags, err := ctrl.tagUC.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, Response{Code: 0, Message: "success", Data: tags})
}

// CreateTag 管理端：创建标签
func (ctrl *adminController) CreateTag(c *gin.Context) {
	var req CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	in := usecase.TagCreateInput{
		Name: req.Name,
		Slug: req.Slug,
	}
	if err := ctrl.tagUC.Create(c.Request.Context(), &in); err != nil {
		c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, Response{Code: 0, Message: "success", Data: in})
}

// DeleteTag 删除标签
func (ctrl *adminController) DeleteTag(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.tagUC.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, Response{Code: 0, Message: "success"})
}
