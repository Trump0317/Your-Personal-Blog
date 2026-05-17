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

// ListPosts 列出所有文章
func (ctrl *adminController) ListPosts(c *gin.Context) {
	var input usecase.PostListInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	posts, total, err := ctrl.postUC.List(c.Request.Context(), &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 在 Controller 层聚合数据
	items := make([]*usecase.PostItem, 0, len(posts))
	for _, p := range posts {
		var category *model.Category
		if p.CategoryID != "" {
			category, _ = ctrl.categoryUC.Get(c.Request.Context(), p.CategoryID)
		}

		var tags []*model.Tag
		if len(p.TagIDs) > 0 {
			tags, _ = ctrl.tagUC.ListByIDs(c.Request.Context(), p.TagIDs)
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
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&body); err == nil {
		if body.Username == ctrl.cfg.Admin.Username && body.Password == ctrl.cfg.Admin.Password {
			c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
			return
		}
	}

	user, pwd, ok := c.Request.BasicAuth()
	if ok && user == ctrl.cfg.Admin.Username && pwd == ctrl.cfg.Admin.Password {
		c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
		return
	}

	c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
	c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
}

// CreatePost 创建文章
func (ctrl *adminController) CreatePost(c *gin.Context) {
	var body struct {
		Title       string   `json:"title"`
		Slug        string   `json:"slug"`
		Content     string   `json:"content"`
		Summary     string   `json:"summary"`
		CategoryID  string   `json:"category_id"`
		Status      string   `json:"status"` // 接收字符串
		TagIDs      []string `json:"tag_ids"`
		NewTagNames []string `json:"new_tag_names"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "details": "binding error"})
		return
	}

	// 在 Controller 层处理标签：将名称转换为 ID
	var finalTagIDs []string
	if len(body.NewTagNames) > 0 {
		ids, err := ctrl.tagUC.GetOrCreates(c.Request.Context(), body.NewTagNames)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process tags"})
			return
		}
		finalTagIDs = ids
	}

	// 合并已有 ID
	if len(body.TagIDs) > 0 {
		tagMap := make(map[string]struct{})
		for _, id := range finalTagIDs {
			tagMap[id] = struct{}{}
		}
		for _, id := range body.TagIDs {
			tagMap[id] = struct{}{}
		}
		finalTagIDs = make([]string, 0, len(tagMap))
		for id := range tagMap {
			finalTagIDs = append(finalTagIDs, id)
		}
	}

	input := usecase.PostCreateInput{
		Title:      body.Title,
		Slug:       body.Slug,
		Content:    body.Content,
		Summary:    body.Summary,
		CategoryID: body.CategoryID,
		TagIDs:     finalTagIDs,
	}

	switch body.Status {
	case "published":
		input.Status = model.PostPublished
	default:
		input.Status = model.PostDraft
	}

	id, err := ctrl.postUC.Create(c.Request.Context(), &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// UpdatePost 更新文章
func (ctrl *adminController) UpdatePost(c *gin.Context) {
	var body struct {
		Title       *string  `json:"title"`
		Slug        *string  `json:"slug"`
		Content     *string  `json:"content"`
		Summary     *string  `json:"summary"`
		CategoryID  *string  `json:"category_id"`
		Status      *string  `json:"status"`
		TagIDs      []string `json:"tag_ids"`
		NewTagNames []string `json:"new_tag_names"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var finalTagIDs []string
	if len(body.NewTagNames) > 0 {
		ids, err := ctrl.tagUC.GetOrCreates(c.Request.Context(), body.NewTagNames)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process tags"})
			return
		}
		finalTagIDs = ids
	}

	if len(body.TagIDs) > 0 {
		tagMap := make(map[string]struct{})
		for _, id := range finalTagIDs {
			tagMap[id] = struct{}{}
		}
		for _, id := range body.TagIDs {
			tagMap[id] = struct{}{}
		}
		finalTagIDs = make([]string, 0, len(tagMap))
		for id := range tagMap {
			finalTagIDs = append(finalTagIDs, id)
		}
	}

	input := usecase.PostUpdateInput{
		ID:         c.Param("id"),
		Title:      body.Title,
		Slug:       body.Slug,
		Content:    body.Content,
		Summary:    body.Summary,
		CategoryID: body.CategoryID,
	}

	if len(finalTagIDs) > 0 || len(body.TagIDs) > 0 {
		input.TagIDs = &finalTagIDs
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

// DeletePost 删除文章
func (ctrl *adminController) DeletePost(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.postUC.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// ListCategories 获取全量分类
func (ctrl *adminController) ListCategories(c *gin.Context) {
	cats, err := ctrl.categoryUC.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cats)
}

// ListTags 获取全量标签
func (ctrl *adminController) ListTags(c *gin.Context) {
	tags, err := ctrl.tagUC.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tags)
}
