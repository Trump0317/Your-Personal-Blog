package usercase

import "io"

// FileUploadInput 文件上传输入参数。
type FileUploadInput struct {
	APIKey   string    // 用户唯一标识
	Content  io.Reader // 文件流
	FileName string    // 原始文件名
	Size     int64     // 文件大小
	Usage    string    // 用途（如：post_cover）
}

type FileListInput struct {
	APIKey string `form:"-"`        // 用户唯一标识
	Page   int    `form:"page"`     // 页码
	Size   int    `form:"pageSize"` // 每页大小
}

type UserCreateInput struct {
	APIKey       string // 用户唯一标识（可选，不传则自动生成）
	InitialQuota int64  // 初始配额（可选，单位字节）
}

type UserUpdateInput struct {
	QuotaLimit *int64 // 可选的配额限制更新
	Status     *int   // 可选的状态更新（0: 活跃, 1: 禁用, 2: 删除）
}
