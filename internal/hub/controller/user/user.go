package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ypb/your-personal-blog/internal/hub/usercase"
	"github.com/ypb/your-personal-blog/internal/hub/usercase/input"
)

type UserController struct {
	userUC usercase.User
}

func NewUserController(userUC usercase.User) *UserController {
	return &UserController{
		userUC: userUC,
	}
}

// Register 用户注册接口
func (u *UserController) Register(c *gin.Context) {
	// 目前简单处理，不需要传递参数，直接创建用户并分配 APIKey
	var in input.UserCreate
	if err := c.ShouldBindJSON(&in); err != nil && c.Request.ContentLength > 0 {
		// 如果有 Body 但格式不对，报错
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	detail, err := u.userUC.Create(c.Request.Context(), in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, detail)
}
