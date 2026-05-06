package usercase

import (
	"context"
	"errors"
	"time"

	"github.com/ypb/your-personal-blog/internal/hub/usercase/input"
	"github.com/ypb/your-personal-blog/internal/hub/usercase/output"
)

var (
	// ErrStorage Storage 错误哨兵。
	ErrStorage = errors.New("storage")
	// ErrDatabase Database 错误哨兵。
	ErrDatabase = errors.New("database")
)

type File interface {
	// 处理文件上传：涉及存物理文件和存数据库记录两步
	Upload(ctx context.Context, userID string, in input.FileUpload) (*output.FileDetail, error)

	// 根据 ID 获取文件信息
	GetByID(ctx context.Context, expires time.Duration, id string) (*output.FileDetail, error)

	// 获取用户下的所有文件
	ListAll(ctx context.Context, userID string) ([]*output.FileDetail, error)

	// 安全删除：同时删除物理文件和数据库记录
	Delete(ctx context.Context, id string) error
}

type User interface {
	// 创建用户
	Create(ctx context.Context, in input.UserCreate) (*output.UserDetail, error)

	// 获取用户信息
	GetByID(ctx context.Context, id string) (*output.UserDetail, error)

	GetByAPIKey(ctx context.Context, apiKey string) (*output.UserDetail, error)

	// 删除用户
	Delete(ctx context.Context, id string) error
}

type Admin interface {
	GetSystemStats(ctx context.Context) (*output.SystemStats, error)
	VerifyAPIKey(apiKey string) bool
}
