package tag

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterPublic 注册公开路由
func (h *Handler) RegisterPublic(r *gin.RouterGroup) {
	r.GET("/tags", h.List)
}

// RegisterAdmin 注册管理端路由
func (h *Handler) RegisterAdmin(r *gin.RouterGroup) {
	r.GET("/tags", h.List)
	r.POST("/tags", h.Create)
	r.DELETE("/tags/:id", h.Delete)
}

func (h *Handler) List(c *gin.Context) {
	tags, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if tags == nil {
		tags = []*Tag{}
	}
	c.JSON(http.StatusOK, tags)
}

type createTagInput struct {
	Name string `json:"name" binding:"required"`
}

func (h *Handler) Create(c *gin.Context) {
	var in createTagInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name required"})
		return
	}
	id, err := h.svc.CreateOrGet(c.Request.Context(), in.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
