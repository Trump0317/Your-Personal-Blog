package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// LocalStorage 实现 StorageEngine 接口，使用本地文件系统存储。
type LocalStorage struct {
	rootDir string
}

// NewLocalStorage 创建一个新的本地存储实例。
func NewLocalStorage(rootDir string) (*LocalStorage, error) {
	// 确保根目录存在
	if err := os.MkdirAll(rootDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create root directory: %w", err)
	}
	return &LocalStorage{rootDir: rootDir}, nil
}

func (s *LocalStorage) Put(ctx context.Context, path string, reader io.Reader) (string, error) {
	fullPath := filepath.Join(s.rootDir, path)

	// 确保父目录存在
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return "", fmt.Errorf("failed to create directory for file: %w", err)
	}

	out, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, reader)
	if err != nil {
		return "", fmt.Errorf("failed to copy data: %w", err)
	}

	return fullPath, nil
}

func (s *LocalStorage) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	fullPath := filepath.Join(s.rootDir, path)
	f, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %s", path)
		}
		return nil, err
	}
	return f, nil
}

func (s *LocalStorage) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(s.rootDir, path)
	err := os.Remove(fullPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (s *LocalStorage) Exists(ctx context.Context, path string) (bool, error) {
	fullPath := filepath.Join(s.rootDir, path)
	_, err := os.Stat(fullPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (s *LocalStorage) URL(ctx context.Context, path string) (string, error) {
	// 对于本地文件系统，通常返回一个相对路径或特定的下载 API URL
	// 这里暂且返回路径本身
	return path, nil
}
