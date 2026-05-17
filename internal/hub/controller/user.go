package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ypb/your-personal-blog/internal/hub/usercase"
	"github.com/ypb/your-personal-blog/internal/pkg/logger"
)

func NewUserRouter(group *gin.RouterGroup, l logger.Interface, file usercase.File, user usercase.User) {
	userHandler := NewUserHandler(file, user)

	userGroup := group.Group("/user")
	{
		// 注册接口 - 公开
		userGroup.POST("/register", userHandler.Register)
	}
	// V1 API 规范
	fileGroup := group.Group("/file")
	{
		// 1. 上传文件 - API Key 鉴权
		fileGroup.POST("/upload", AuthMiddleware(user), userHandler.Upload)

		// 2. 获取单个文件 - 公开/Token
		fileGroup.GET("/:id", userHandler.GetByID)

		// 3. 获取用户全部文件 - API Key 鉴权
		fileGroup.GET("/all", AuthMiddleware(user), userHandler.ListAll)

		// 4. 删除文件 - API Key 鉴权
		fileGroup.DELETE("/:id", AuthMiddleware(user), userHandler.Delete)
	}
}

type UserHandler struct {
	fileUC usercase.File
	userUC usercase.User
}

func NewUserHandler(fileUC usercase.File, userUC usercase.User) *UserHandler {
	return &UserHandler{
		fileUC: fileUC, // 依赖注入 File 用于处理文件相关的业务逻辑
		userUC: userUC, // 依赖注入 User 用于处理用户注册相关的业务逻辑
	}
}

// Register 用户注册接口
func (u *UserHandler) Register(c *gin.Context) {
	// 前端传回的数据是一个字符串（APIKey），配额暂时用硬编码
	var apiKey string
	if err := c.ShouldBindJSON(&apiKey); err != nil && c.Request.Context().Value(gin.BodyBytesKey) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body, expected string"})
		return
	}

	in := usercase.UserCreateInput{
		APIKey:       apiKey,
		InitialQuota: 1024 * 1024 * 100, // 硬编码 100MB
	}

	detail, err := u.userUC.Create(c.Request.Context(), in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, detail)
}

// Upload 上传文件处理函数
func (f *UserHandler) Upload(c *gin.Context) {
	userID := f.getUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file is uploaded"})
		return
	}
	defer file.Close()

	in := usercase.FileUploadInput{
		APIKey:   userID,
		Content:  file,
		FileName: header.Filename,
		Size:     header.Size,
		Usage:    c.DefaultPostForm("usage", "general"),
	}

	detail, err := f.fileUC.Upload(c.Request.Context(), in)
	if err != nil {
		f.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, detail)
}

// GetByID 获取文件详细信息
func (f *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	detail, err := f.fileUC.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	c.JSON(http.StatusOK, detail)
}

// ListAll 获取用户拥有的全部文件
func (f *UserHandler) ListAll(c *gin.Context) {
	userID := f.getUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var in usercase.FileListInput
	if err := c.ShouldBindQuery(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters"})
		return
	}
	in.APIKey = userID

	files, err := f.fileUC.ListAll(c.Request.Context(), in)
	if err != nil {
		f.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, files)
}

// Delete 删除文件接口
func (f *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := f.fileUC.Delete(c.Request.Context(), id)
	if err != nil {
		f.handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (f *UserHandler) handleError(c *gin.Context, err error) {
	if errors.Is(err, usercase.ErrStorage) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "storage error"})
	} else if errors.Is(err, usercase.ErrDatabase) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func (f *UserHandler) getUserID(c *gin.Context) string {
	if val, ok := c.Get(ContextUserKey); ok {
		if u, ok := val.(*usercase.UserDetailOutput); ok {
			return u.APIKey
		}
	}
	return ""
}
