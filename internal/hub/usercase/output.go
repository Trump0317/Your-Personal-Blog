package usercase

import (
	"time"
)

// FileDetailOutput 文件详情输出参数。
type FileDetailOutput struct {
	ID        string
	FileName  string
	FileSize  int64
	MimeType  string
	Usage     string
	CreatedAt time.Time
}

// SystemStatsOutput 系统统计输出
type SystemStatsOutput struct {
	TotalFiles       int64 `json:"total_files"`
	TotalUsers       int64 `json:"total_users"`
	TotalUsedStorage int64 `json:"total_used_storage"`
}

type UserCreateOutput struct {
	APIKey string
}

type UserDetailOutput struct {
	APIKey       string
	QuotaLimit   int64
	CurrentUsage int64
	Status       int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type FileUploadOutput struct {
	ID string
}
