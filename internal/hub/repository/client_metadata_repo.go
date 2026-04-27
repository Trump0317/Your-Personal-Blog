package repository

import (
	"context"

	"github.com/ypb/your-personal-blog/internal/hub/model"
)

// ClientMetadataRepository 定义租户（Blog Client）信息的存储契约
// 它是Storage Hub 多租户隔离和安全准入的第一道防线
type ClientMetadataRepository interface {
	// Create 注册一个新的租户信息
	// 功能：在数据库中创建租户记录，包含分配初始配额和状态
	// 参数：
	//   - ctx: 上下文
	//   - meta: 租户元数据模型指针
	// 返回值：
	//   - error: 操作失败时返回错误
	Create(ctx context.Context, meta *model.ClientMetadata) error

	// GetByAPIKey 根据 API Key 检索租户信息
	// 功能：通过 API Key 识别请求方身份，用于鉴权中间件。
	// 参数：
	//   - ctx: 上下文
	//   - apiKey: 请求头中携带的鉴权密钥
	// 返回值：
	//   - *model.ClientMetadata: 租户元数据模型
	//   - error: 密钥无效时返回 repository.ErrNotFound；数据库异常返回 repository.ErrInternal。
	GetByAPIKey(ctx context.Context, apiKey string) (*model.ClientMetadata, error)

	// UpdateStatus 变更租户的运行状态
	// 功能：更新租户状态（如因欠费或违规暂停服务）
	// 参数：
	//   - ctx: 上下文
	//   - clientID: 租户唯一标识
	//   - status: 目标状态 (ClientStatus)
	// 返回值：
	//   - error: 更新失败时返回错误
	UpdateStatus(ctx context.Context, clientID string, status model.ClientStatus) error

	// Delete 移除租户信息
	// 功能：从数据库中永久删除租户记录
	// 参数：
	//   - ctx: 传递上下文信号。
	//   - clientID: 租户唯一标识
	// 返回值：
	//   - error: 操作失败时返回错误
	Delete(ctx context.Context, clientID string) error

	// GetTotalUsage 获取租户当前已消耗的存储空间
	// 功能：返回租户已使用的总字节数，通常用于配额余量校验
	// 参数：
	//   - ctx: 传递上下文信号。如果是实时计算模式，应尊重 ctx 的超时设置。
	//   - clientID: 租户唯一标识
	// 返回值：
	//   - int64: 已使用的总字节数
	//   - error: 统计失败时返回错误
	GetTotalUsage(ctx context.Context, clientID string) (int64, error)
}
