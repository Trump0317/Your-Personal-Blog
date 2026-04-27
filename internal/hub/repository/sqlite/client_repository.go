package sqlite

import (
	"context"
	"database/sql"

	"github.com/ypb/your-personal-blog/internal/hub/model"
	"github.com/ypb/your-personal-blog/internal/hub/repository"
)

// clientRepository 实现了 repository.ClientMetadataRepository 接口
// 与SQLite数据库交互，负责租户信息的持久化操作
type clientRepository struct {
	db *sql.DB
}

// NewClientRepository 创建一个 SQLite 实现的 Client 仓库
func NewClientRepository(db *sql.DB) repository.ClientMetadataRepository {
	return &clientRepository{db: db}
}

// Create TODO: 注册新的租户及其初始配额 (基础设施/SQLite)
func (r *clientRepository) Create(ctx context.Context, meta *model.ClientMetadata) error {
	return nil
}

// GetByID TODO: 根据租户 ID 查询详细信息 (基础设施/SQLite)
func (r *clientRepository) GetByID(ctx context.Context, id string) (*model.ClientMetadata, error) {
	return nil, nil
}

// GetByAPIKey TODO: 通过 API Key 检索租户信息，用于鉴权校验 (基础设施/SQLite)
func (r *clientRepository) GetByAPIKey(ctx context.Context, apiKey string) (*model.ClientMetadata, error) {
	return nil, nil
}

// Delete TODO: 移除租户信息 (基础设施/SQLite)
func (r *clientRepository) Delete(ctx context.Context, id string) error {
	return nil
}

// UpdateStatus TODO: 启用或禁用租户账号 (基础设施/SQLite)
func (r *clientRepository) UpdateStatus(ctx context.Context, id string, status model.ClientStatus) error {
	return nil
}

// GetTotalUsage TODO: 实时计算租户的所有文件大小总和 (基础设施/SQLite)
func (r *clientRepository) GetTotalUsage(ctx context.Context, clientID string) (int64, error) {
	return 0, nil
}
