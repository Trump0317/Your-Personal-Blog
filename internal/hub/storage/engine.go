package storage

import (
	"context"
	"io"
)

// StorageEngine 定义了物理文件存储的抽象接口。
// 它不关心业务元数据，只负责字节流的读写和物理管理。
type StorageEngine interface {
	// Put 将数据流保存到指定路径。
	// path 通常是生成的唯一键或相对路径。
	// 返回保存后的物理地址或标识。
	Put(ctx context.Context, path string, reader io.Reader) (physicalPath string, err error)

	// Get 读取指定路径的文件内容。
	// 调用者负责关闭返回的 ReadCloser。
	Get(ctx context.Context, path string) (io.ReadCloser, error)

	// Delete 删除物理文件。
	Delete(ctx context.Context, path string) error

	// Exists 检查物理文件是否存在。
	Exists(ctx context.Context, path string) (bool, error)

	// URL 获取文件的访问链接（如果存储引擎支持公网访问或签名 URL）。
	URL(ctx context.Context, path string) (string, error)
}
