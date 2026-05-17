package repo

import (
	"context"
	"io"
	"time"

	model "github.com/ypb/your-personal-blog/internal/hub/model"
)

// FileStore 负责物理存储的接口。
type FileStore interface {
	Upload(ctx context.Context, id string, expires time.Duration, reader io.Reader) error
	Download(ctx context.Context, id string, expires time.Duration) (io.ReadCloser, error)
	Delete(ctx context.Context, id string) error
}

// FileRepo 负责数据库 SQL 操作的接口。
type FileRepo interface {
	Create(ctx context.Context, file model.File) (string, error)                             // 保存文件记录，返回 ID
	GetByID(ctx context.Context, id string) (*model.File, error)                             // 根据 ID 获取文件记录
	ListByUser(ctx context.Context, apikey string, limit, offset int) ([]*model.File, error) // 获取用户下的所有文件
	Delete(ctx context.Context, id string) error                                             // 根据 ID 删除文件记录
	UpdateStatus(ctx context.Context, id string, status model.FileStatus) error              // 更新文件状态
}

type UserRepo interface {
	Create(ctx context.Context, user model.User) (string, error)
	GetByAPIKey(ctx context.Context, apiKey string) (*model.User, error)
	Delete(ctx context.Context, apiKey string) error
	UpdateStatus(ctx context.Context, apiKey string, status model.UserStatus) error
	UpdateUsage(ctx context.Context, apiKey string, delta int64) error
}
