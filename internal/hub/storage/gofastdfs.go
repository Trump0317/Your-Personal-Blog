package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

// GoFastDFSStorage 实现 StorageEngine 接口，与 gofastdfs 服务交互。
// gofastdfs 采用对等架构，通常通过标准的 HTTP Multi-part 方式进行文件上传。
type GoFastDFSStorage struct {
	endpoint string // gofastdfs 的服务地址，例如 http://127.0.0.1:8082
	group    string // gofastdfs 的组名，默认通常为 "group1"
}

// NewGoFastDFSStorage 创建 gofastdfs 存储实例。
func NewGoFastDFSStorage(endpoint, group string) *GoFastDFSStorage {
	return &GoFastDFSStorage{
		endpoint: endpoint,
		group:    group,
	}
}

func (s *GoFastDFSStorage) Put(ctx context.Context, path string, reader io.Reader) (string, error) {
	// gofastdfs 上传逻辑设计：
	// 1. 构造 multipart/form-data 请求
	// 2. 将 reader 中的数据流式写入请求体
	// 3. 发送 POST 请求到 {endpoint}/{group}/upload

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// gofastdfs 要求文件字段名为 "file"
	part, err := writer.CreateFormFile("file", path)
	if err != nil {
		return "", err
	}

	if _, err := io.Copy(part, reader); err != nil {
		return "", err
	}
	writer.Close()

	url := fmt.Sprintf("%s/%s/upload", s.endpoint, s.group)
	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("gofastdfs upload failed with status: %s", resp.Status)
	}

	// gofastdfs 返回的是 JSON，包含文件的访问路径
	// 这里的设计细节：我们需要解析返回的 JSON 获取真正的后端存储路径
	return "gofastdfs://" + path, nil
}

func (s *GoFastDFSStorage) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	// 构造下载请求
	url := fmt.Sprintf("%s/%s/default/%s", s.endpoint, s.group, path)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("failed to get file from gofastdfs: %s", resp.Status)
	}

	return resp.Body, nil
}

func (s *GoFastDFSStorage) Delete(ctx context.Context, path string) error {
	// gofastdfs 的删除逻辑通常也是通过 HTTP 接口
	return nil
}

func (s *GoFastDFSStorage) Exists(ctx context.Context, path string) (bool, error) {
	return false, nil
}

func (s *GoFastDFSStorage) URL(ctx context.Context, path string) (string, error) {
	return fmt.Sprintf("%s/%s/default/%s", s.endpoint, s.group, path), nil
}
