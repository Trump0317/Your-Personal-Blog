package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ypb/your-personal-blog/internal/hub/repository"
)

// Response 是统一的 API 响应格式
type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// HandleError 负责将业务逻辑错误或 Repository 错误转换为适当的 HTTP 状态码
func HandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	var statusCode int
	var message string

	// 判断错误类型并进行状态码映射
	switch {
	case errors.Is(err, repository.ErrNotFound):
		statusCode = http.StatusNotFound
		message = "resource not found"
	case errors.Is(err, repository.ErrAlreadyExists):
		statusCode = http.StatusConflict
		message = "resource already exists"
	case errors.Is(err, repository.ErrInternal):
		statusCode = http.StatusInternalServerError
		message = "internal server error"
	default:
		// 未知错误统一处理为 500
		statusCode = http.StatusInternalServerError
		message = "an unexpected error occurred"
	}

	c.AbortWithStatusJSON(statusCode, Response{
		Message: message,
	})
}
