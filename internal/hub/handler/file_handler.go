package handler

import (
	"net/http"

	"github.com/ypb/your-personal-blog/internal/hub/service"

	"github.com/gin-gonic/gin"
)

// FileHandler
// 处理文件相关的 HTTP 请求
type FileHandler struct {
	fileSvc service.FileService
	authSvc service.AuthService
}

// 创建 FileHandler 实例
func NewFileHandler(fileSvc service.FileService, authSvc service.AuthService) *FileHandler {
	return &FileHandler{fileSvc: fileSvc, authSvc: authSvc}
}

// Upload 上传文件接口
func (h *FileHandler) Upload(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "upload endpoint linked"})
}

// Download 下载文件接口
func (h *FileHandler) Download(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "download endpoint linked"})
}

// Delete 删除文件接口
func (h *FileHandler) Delete(c *gin.Context) {
	c.Status(204)
}

// List 列出文件接口
func (h *FileHandler) List(c *gin.Context) {
	c.JSON(200, gin.H{"files": []interface{}{}})
}

// GetMetadata 获取文件元数据接口
func (h *FileHandler) GetMetadata(c *gin.Context) {
	c.Status(200)
}

// GetClientStats 获取客户端配额和使用统计接口
func (h *FileHandler) GetClientStats(c *gin.Context) {
	c.JSON(200, gin.H{"quota_limit": 0, "used_size": 0})
}

// RegisterRoutes 注册路由
func (h *FileHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/files")
	{
		// 基础操作
		api.POST("/upload", h.Upload)           // 上传
		api.GET("/:id", h.Download)             // 下载
		api.DELETE("/:id", h.Delete)            // 删除
		api.GET("", h.List)                     // 列出文件
		api.GET("/:id/metadata", h.GetMetadata) // 获取单文件元数据
		api.GET("/stats", h.GetClientStats)     // 获取配额与磁盘占用统计
	}
}
