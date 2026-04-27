package sqlite

import (
	"context"
	"database/sql"

	"github.com/ypb/your-personal-blog/internal/hub/model"
	"github.com/ypb/your-personal-blog/internal/hub/repository"
)

// fileRepository 实现了 repository.FileMetadataRepository 接口
// 与SQLite数据库交互，负责文件元数据的持久化操作
type fileRepository struct {
	db *sql.DB
}

func NewFileRepository(db *sql.DB) repository.FileMetadataRepository {
	return &fileRepository{db: db}
}

// Create TODO: 实现文件的元数据入库操作 (基础设施/SQLite)
func (r *fileRepository) Create(ctx context.Context, meta *model.FileMetadata) error {
	return nil
}

// GetByID TODO: 根据文件 ID 查询完整的元数据信息 (基础设施/SQLite)
func (r *fileRepository) GetByID(ctx context.Context, id string) (*model.FileMetadata, error) {
	return nil, nil
}

// UpdateStatus TODO: 更新文件的生命周期状态 (如从 Uploading 变为 Active) (基础设施/SQLite)
func (r *fileRepository) UpdateStatus(ctx context.Context, id string, status model.FileStatus) error {
	return nil
}

// Delete TODO: 从数据库中物理删除或标记删除文件记录 (基础设施/SQLite)
func (r *fileRepository) Delete(ctx context.Context, id string) error {
	return nil
}

// ListByClient TODO: 分页查询特定租户下的文件列表 (基础设施/SQLite)
func (r *fileRepository) ListByClient(ctx context.Context, clientID string, offset, limit int) ([]*model.FileMetadata, int64, error) {
	return nil, 0, nil
}

// GetTotalSizeByClient TODO: 统计特定租户已占用的存储总空间 (基础设施/SQLite)
func (r *fileRepository) GetTotalSizeByClient(ctx context.Context, clientID string) (int64, error) {
	return 0, nil
}
