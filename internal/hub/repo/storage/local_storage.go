package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type LocalFileStore struct {
	basePath string
}

// NewLocalFileStore 创建一个本地文件存储实例
// basePath 是文件存储的根目录，例如 "uploads"
func NewLocalFileStore(basePath string) (*LocalFileStore, error) {
	// 确保根目录存在
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create base path: %w", err)
	}
	return &LocalFileStore{basePath: basePath}, nil
}

// Upload 将文件保存到本地磁盘
func (s *LocalFileStore) Upload(ctx context.Context, path string, expires time.Duration, reader io.Reader) error {
	fullPath := filepath.Join(s.basePath, path)

	// 确保目标目录存在
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// 创建目标文件
	out, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", fullPath, err)
	}
	defer out.Close()

	// 写入内容
	if _, err = io.Copy(out, reader); err != nil {
		return fmt.Errorf("failed to write file %s: %w", fullPath, err)
	}

	return nil
}

// Download 从本地磁盘读取文件
func (s *LocalFileStore) Download(ctx context.Context, path string, expires time.Duration) (io.ReadCloser, error) {
	fullPath := filepath.Join(s.basePath, path)
	file, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %s", path)
		}
		return nil, fmt.Errorf("failed to open file %s: %w", fullPath, err)
	}
	return file, nil
}

// Delete 从本地磁盘删除文件
func (s *LocalFileStore) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(s.basePath, path)
	if err := os.Remove(fullPath); err != nil {
		if os.IsNotExist(err) {
			return nil // 文件不存在视为删除成功
		}
		return fmt.Errorf("failed to delete file %s: %w", fullPath, err)
	}
	return nil
}
