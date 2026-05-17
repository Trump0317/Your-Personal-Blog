package model

import "time"

type FileStatus int

const (
	FileUploading FileStatus = iota // 0: 上传中
	FileActive                      // 1: 正常
	FileDeleted                     // 2: 已逻辑删除
)

// 文件模型 File
// ID : 文件ID
// Uploader : 上传者APIKey
// FileName : 文件名
// MimeType : 文件类型
// FileSize : 文件大小
// StorageType : 存储类型
// StoragePath : 存储路径,文件二进制本体位置
// CreatedAt : 创建时间戳
// UpdatedAt : 更新时间戳
// Status : 文件状态
type File struct {
	ID          string     `json:"id"`
	Uploader    string     `json:"uploader"`
	FileName    string     `json:"file_name"`
	MimeType    string     `json:"mime_type"`
	FileSize    int64      `json:"file_size"`
	Usage       string     `json:"usage"`
	StoragePath string     `json:"storage_path"`
	Status      FileStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
