package usercase

import (
	"context"
	"errors"
)

var (
	// ErrStorage Storage 错误哨兵。
	ErrStorage = errors.New("storage")
	// ErrDatabase Database 错误哨兵。
	ErrDatabase = errors.New("database")
)

type File interface {
	// 文件上传：涉及存物理文件和存数据库记录两步
	Upload(ctx context.Context, in FileUploadInput) (*FileUploadOutput, error)

	// 删除：同时删除物理文件和数据库记录
	Delete(ctx context.Context, id string) error

	// 根据文件ID 获取文件信息
	GetByID(ctx context.Context, id string) (*FileDetailOutput, error)

	// 获取用户的所有文件
	ListAll(ctx context.Context, in FileListInput) ([]*FileDetailOutput, error)
}

// User 相关用例接口定义
type User interface {
	// 创建用户
	Create(ctx context.Context, in UserCreateInput) (*UserCreateOutput, error)

	// 根据 API Key 获取用户信息
	GetByAPIKey(ctx context.Context, apiKey string) (*UserDetailOutput, error)

	// 删除用户
	Delete(ctx context.Context, apiKey string) error

	// 更新用户信息
	Update(ctx context.Context, apiKey string, in UserUpdateInput) error
}

type Admin interface {
	GetSystemStats(ctx context.Context) (*SystemStatsOutput, error)
}
