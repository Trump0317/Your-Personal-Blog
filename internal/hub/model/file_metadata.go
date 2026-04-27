package model

type FileStatus int

const (
	FileUploading FileStatus = iota // 0: 上传中
	FileActive                      // 1: 正常
	FileDeleted                     // 2: 已逻辑删除
)

// 文件元数据 FileMetadata
// ID : 文件ID
// ClientID : 客户端ID
// BucketID : 存储桶ID
// OriginalName : 原始文件名
// MimeType : 文件类型
// FileSize : 文件大小
// StoragePath : 存储路径
// CreatedAt : 创建时间戳
// UpdatedAt : 更新时间戳
// Status : 文件状态
type FileMetadata struct {
	ID           string     `json:"id"`
	ClientID     string     `json:"client_id"`
	BucketID     string     `json:"bucket_id"`
	OriginalName string     `json:"original_name"`
	MimeType     string     `json:"mime_type"`
	FileSize     int64      `json:"file_size"`
	StoragePath  string     `json:"storage_path"`
	CreatedAt    int64      `json:"created_at"`
	UpdatedAt    int64      `json:"updated_at"`
	Status       FileStatus `json:"status"`
}
