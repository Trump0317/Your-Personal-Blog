package post

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// ── 公开路由 ──

func (h *Handler) RegisterPublic(r *gin.RouterGroup) {
	r.GET("/stats", h.Stats)
	r.GET("/archive", h.Archive)
	r.GET("/posts", h.ListPublic)
	r.GET("/posts/:slug", h.GetPublic)
}

// ── 管理端路由 ──

func (h *Handler) RegisterAdmin(r *gin.RouterGroup) {
	r.GET("/posts", h.ListAdmin)
	r.GET("/posts/:id", h.Get)
	r.POST("/posts", h.Create)
	r.PUT("/posts/:id", h.Update)
	r.DELETE("/posts/:id", h.Delete)
	r.PUT("/posts/:id/publish", h.Publish)
	r.PUT("/posts/:id/unpublish", h.Unpublish)
}

// ── 公开接口 ──

func (h *Handler) ListPublic(c *gin.Context) {
	limit := 12
	if v := c.Query("size"); v != "" {
		if n, _ := strconv.Atoi(v); n > 0 && n <= 50 {
			limit = n
		}
	}
	if v := c.Query("page_size"); v != "" {
		if n, _ := strconv.Atoi(v); n > 0 && n <= 50 {
			limit = n
		}
	}
	f := Filter{Limit: limit}
	if v := c.Query("page"); v != "" {
		page, _ := strconv.Atoi(v)
		if page > 0 {
			f.Offset = (page - 1) * limit
		}
	}
	if v := c.Query("category_id"); v != "" {
		f.CategoryID = v
	}
	if v := c.Query("category"); v != "" {
		if f.CategoryID == "" {
			f.CategoryID = v
		}
	}
	if v := c.Query("tag_id"); v != "" {
		f.TagID = v
	}
	if v := c.Query("tag"); v != "" {
		if f.TagID == "" {
			f.TagID = v
		}
	}
	if v := c.Query("year"); v != "" {
		f.Year, _ = strconv.Atoi(v)
	}
	if v := c.Query("month"); v != "" {
		f.Month, _ = strconv.Atoi(v)
	}
	if v := c.Query("q"); v != "" {
		f.Query = v
	}
	pub := true
	f.Published = &pub

	posts, total, err := h.svc.List(c.Request.Context(), f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"total": total,
		"page":  f.Offset/f.Limit + 1,
	})
}

func (h *Handler) GetPublic(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "slug required"})
		return
	}
	post, err := h.svc.Get(c.Request.Context(), slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}
	if !post.Published {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}
	c.JSON(http.StatusOK, post)
}

func (h *Handler) Get(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "slug required"})
		return
	}
	post, err := h.svc.Get(c.Request.Context(), slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}
	c.JSON(http.StatusOK, post)
}

// ── 管理端接口 ──

func (h *Handler) ListAdmin(c *gin.Context) {
	f := Filter{Limit: 50}
	if v := c.Query("page"); v != "" {
		page, _ := strconv.Atoi(v)
		if page > 0 {
			f.Offset = (page - 1) * f.Limit
		}
	}
	// 管理员可以看到所有状态
	posts, total, err := h.svc.List(c.Request.Context(), f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"total": total,
	})
}

type createPostBody struct {
	Title      string   `json:"title" binding:"required"`
	Slug       string   `json:"slug"`
	Content    string   `json:"content"`
	Summary    string   `json:"summary"`
	CategoryID string   `json:"category_id"`
	TagIDs     []string `json:"tag_ids"`
	Published  bool     `json:"published"`
}

func (h *Handler) Create(c *gin.Context) {
	var b createPostBody
	if err := c.ShouldBindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.svc.Create(c.Request.Context(), CreateInput{
		Title:      b.Title,
		Slug:       b.Slug,
		Content:    b.Content,
		Summary:    b.Summary,
		CategoryID: b.CategoryID,
		TagIDs:     b.TagIDs,
		Published:  b.Published,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

type updatePostBody struct {
	Title      *string  `json:"title"`
	Slug       *string  `json:"slug"`
	Content    *string  `json:"content"`
	Summary    *string  `json:"summary"`
	CategoryID *string  `json:"category_id"`
	TagIDs     []string `json:"tag_ids"`
	Published  *bool    `json:"published"`
}

func (h *Handler) Update(c *gin.Context) {
	var b updatePostBody
	if err := c.ShouldBindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.Update(c.Request.Context(), UpdateInput{
		ID:         c.Param("id"),
		Title:      b.Title,
		Slug:       b.Slug,
		Content:    b.Content,
		Summary:    b.Summary,
		CategoryID: b.CategoryID,
		TagIDs:     b.TagIDs,
		Published:  b.Published,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) Publish(c *gin.Context) {
	if err := h.svc.Publish(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) Unpublish(c *gin.Context) {
	if err := h.svc.Unpublish(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) Stats(c *gin.Context) {
	stats, err := h.svc.Stats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *Handler) Archive(c *gin.Context) {
	years, err := h.svc.Archive(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, years)
}
