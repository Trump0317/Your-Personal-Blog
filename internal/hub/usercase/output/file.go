package output

import "time"

// FileDetail 文件详情输出参数。
type FileDetail struct {
	ID        string
	FileName  string
	FileSize  int64
	MimeType  string
	Usage     string
	CreatedAt time.Time
}
