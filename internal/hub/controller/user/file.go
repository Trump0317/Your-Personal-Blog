package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ypb/your-personal-blog/internal/hub/controller/middleware"
	"github.com/ypb/your-personal-blog/internal/hub/usercase"
	"github.com/ypb/your-personal-blog/internal/hub/usercase/input"
	"github.com/ypb/your-personal-blog/internal/hub/usercase/output"
)

type FileController struct {
	fileUC usercase.File
	userUC usercase.User
}

func NewFileController(fileUC usercase.File, userUC usercase.User) *FileController {
	return &FileController{
		fileUC: fileUC,
		userUC: userUC,
	}
}

// Upload 上传文件处理函数
func (f *FileController) Upload(c *gin.Context) {
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

	in := input.FileUpload{
		Content:  file,
		FileName: header.Filename,
		Size:     header.Size,
		Usage:    c.DefaultPostForm("usage", "general"),
	}

	detail, err := f.fileUC.Upload(c.Request.Context(), userID, in)
	if err != nil {
		f.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, detail)
}

// GetByID 获取文件详细信息
func (f *FileController) GetByID(c *gin.Context) {
	id := c.Param("id")
	detail, err := f.fileUC.GetByID(c.Request.Context(), 0, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	c.JSON(http.StatusOK, detail)
}

// ListAll 获取用户拥有的全部文件
func (f *FileController) ListAll(c *gin.Context) {
	userID := f.getUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	files, err := f.fileUC.ListAll(c.Request.Context(), userID)
	if err != nil {
		f.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, files)
}

// Delete 删除文件接口
func (f *FileController) Delete(c *gin.Context) {
	id := c.Param("id")
	err := f.fileUC.Delete(c.Request.Context(), id)
	if err != nil {
		f.handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (f *FileController) handleError(c *gin.Context, err error) {
	if errors.Is(err, usercase.ErrStorage) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "storage error"})
	} else if errors.Is(err, usercase.ErrDatabase) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func (f *FileController) getUserID(c *gin.Context) string {
	if val, ok := c.Get(middleware.ContextUserKey); ok {
		if u, ok := val.(*output.UserDetail); ok {
			return u.ID
		}
	}
	return ""
}
