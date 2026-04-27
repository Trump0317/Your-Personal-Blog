package repository

import (
	"context"

	"github.com/ypb/your-personal-blog/internal/hub/model"
)

// FileMetadataRepository 定义文件元数据的存储契约
// 负责所有与文件属性、位置以及逻辑状态相关的持久化操作，支持多租户隔离
type FileMetadataRepository interface {
	// Create 持久化一个新的文件元数据记录
	// 功能：在数据库中创建一条新的文件元数据记录，通常在文件成功上传到物理存储后调用
	// 参数：
	//   - ctx: 传递上下文信号。实现者应监听 ctx.Done() 以处理调用方主动取消上传或数据库写入操作。
	//   - meta: 包含文件详细信息的模型指针
	// 返回值：
	//   - error: 操作失败时返回错误描述，成功返回 nil
	Create(ctx context.Context, meta *model.FileMetadata) error

	// GetByID 根据唯一标识符检索文件元数据
	// 功能：通过文件的唯一标识符从存储中获取完整的元数据结构
	// 参数：
	//   - ctx: 传递上下文信号。实现者应根据 ctx 的 Deadline 决定查询的最长阻塞时间，超时应立即返回
	//   - id: 文件的唯一标识符 (UUID)
	// 返回值：
	//   - *model.FileMetadata: 匹配到的文件元数据指针
	//   - error: 当记录不存在时应返回 repository.ErrNotFound；数据库异常返回 repository.ErrInternal
	GetByID(ctx context.Context, id string) (*model.FileMetadata, error)

	// UpdateStatus 修改文件的逻辑状态
	// 功能：更新指定文件的状态字段（如：从“上传中”变更为“活跃”）
	// 参数：
	//   - ctx: 传递上下文信号。实现者应根据 ctx 的 Deadline 决定查询的最长阻塞时间，超时应立即返回
	//   - id: 文件的唯一标识符
	//   - status: 目标状态 (FileStatus)
	// 返回值：
	//   - error: 操作失败时返回错误
	UpdateStatus(ctx context.Context, id string, status model.FileStatus) error

	// Delete 物理删除文件元数据记录
	// 功能：从数据库中永久删除指定 ID 的元数据记录
	// 参数：
	//   - ctx: 传递上下文信号。实现者应监听 ctx.Done() 以处理调用方主动取消删除操作
	//   - id: 文件的唯一标识符
	// 返回值：
	//   - error: 删除失败时返回错误
	Delete(ctx context.Context, id string) error

	// ListByClient 分页获取指定租户的文件清单
	// 功能：拉取属于某个 Blog Client 的文件元数据清单，并返回符合条件的总记录数
	// 参数：
	//   - ctx: 传递上下文信号。考虑到分页查询可能涉及大数据量 COUNT，实现者应尊重 ctx 的超时设置，严禁导致长时间锁表。、
	//   - clientID: 租户的唯一识别 ID
	//   - offset: 跳过的记录条数，计算公式为 (page-1) * limit
	//   - limit: 本次查询获取的最大条数
	// 返回值：
	//   - []*model.FileMetadata: 文件元数据切片
	//   - int64: 符合查询条件的总记录数
	//   - error: 查询过程中发生数据库异常时返回错误
	ListByClient(ctx context.Context, clientID string, offset, limit int) ([]*model.FileMetadata, int64, error)

	// GetTotalSizeByClient 统计指定租户占用的总存储空间（字节）
	// 功能：累加指定租户下所有“活跃”状态文件的字节数，用于配额校验
	// 参数：
	//   - ctx: 传递上下文信号。鉴于该操作可能触发表全扫描，实现者必须支持超时控制以防阻塞
	//   - clientID: 租户的唯一识别 ID
	// 返回值：
	//   - int64: 已使用的总字节数
	//   - error: 统计计算出错时返回错误
	GetTotalSizeByClient(ctx context.Context, clientID string) (int64, error)
}
