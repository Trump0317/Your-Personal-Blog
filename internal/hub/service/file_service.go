package service

import (
	"context"
	"io"

	"github.com/ypb/your-personal-blog/internal/hub/model"
)

// FileService 定义了文件管理的核心业务逻辑接口。
// 它负责协调 Repository 和物理存储引擎，处理复杂的业务规则（如配额检查、重试逻辑等）。
type FileService interface {
	// Upload 处理文件上传流程。
	// 按照“状态标记+异步对齐”策略：
	// 1. 在数据库创建状态为 Uploading 的元数据。
	// 2. 将流写入物理存储。
	// 3. 更新元数据状态为 Active。
	// 参数:
	//   ctx: 上下文，支持超时和取消。
	//   clientID: 上传者的客户端唯一 ID。
	//   reader: 文件内容流。
	//   fileName: 原始文件名。
	//   size: 文件大小（字节）。
	//   bucket: 存储桶名称。
	// 返回:
	//   *model.FileMetadata: 创建成功后的元数据对象。
	//   error: 配额不足、存储失败或数据库异常
	Upload(ctx context.Context, clientID string, reader io.Reader, fileName string, size int64, bucket string) (*model.FileMetadata, error)

	// Download 获取文件的读取流及元数据。
	// 参数:
	//   ctx: 上下文。
	//   fileID: 文件的唯一标识 UUID。
	// 返回:
	//   io.ReadCloser: 文件流，调用者负责关闭。
	//   *model.FileMetadata: 文件的元数据。
	//   error: 如果文件不存在或状态非 Active，返回 ErrNotFound。
	Download(ctx context.Context, fileID string) (io.ReadCloser, *model.FileMetadata, error)

	// Delete 执行逻辑删除。
	// 为了数据安全，Service 层通常 solo 做逻辑删除，由清理任务执行物理删除。
	// 参数:
	//   ctx: 上下文。
	//   fileID: 文件的唯一标识 UUID。
	// 返回:
	//   error: 如果文件不存在或已删除。
	Delete(ctx context.Context, fileID string) error

	// GetMetadata 获取单个文件的详细信息。
	GetMetadata(ctx context.Context, fileID string) (*model.FileMetadata, error)

	// ListMetadata 分页查询文件列表。
	// 参数:
	//   clientID: 可选，按客户端过滤。
	//   page: 页码，从 1 开始。
	//   pageSize: 每页数量。
	// 返回:
	//   []*model.FileMetadata: 文件元数据列表。
	//   int64: 符合条件的总记录数。
	//   error: 查询异常。
	ListMetadata(ctx context.Context, clientID string, page, pageSize int) ([]*model.FileMetadata, int64, error)
}
