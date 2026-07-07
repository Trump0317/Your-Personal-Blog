package category

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

func (h *Handler) RegisterPublic(r *gin.RouterGroup) {
	r.GET("/categories", h.List)
}

func (h *Handler) RegisterAdmin(r *gin.RouterGroup) {
	r.GET("/categories", h.List)
	r.POST("/categories", h.Create)
	r.PUT("/categories/:id", h.Update)
	r.DELETE("/categories/:id", h.Delete)
}

func (h *Handler) List(c *gin.Context) {
	cats, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if cats == nil {
		cats = []*Category{}
	}
	c.JSON(http.StatusOK, cats)
}

type createCategoryInput struct {
	Name string `json:"name" binding:"required"`
}

func (h *Handler) Create(c *gin.Context) {
	var in createCategoryInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name required"})
		return
	}
	id, err := h.svc.Create(c.Request.Context(), in.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

type updateCategoryInput struct {
	Name string `json:"name" binding:"required"`
}

func (h *Handler) Update(c *gin.Context) {
	var in updateCategoryInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name required"})
		return
	}
	if err := h.svc.Update(c.Request.Context(), c.Param("id"), in.Name); err != nil {
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
