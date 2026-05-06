package input

import "io"

// FileUpload 文件上传输入参数。
type FileUpload struct {
	Content  io.Reader // 文件流
	FileName string    // 原始文件名
	Size     int64     // 文件大小
	Usage    string    // 用途（如：post_cover）
}
