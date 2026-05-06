package model

// SystemStats 系统概览统计模型
type SystemStats struct {
	TotalFiles       int64 `json:"total_files"`        // 总文件数量
	TotalUsers       int64 `json:"total_users"`        // 总用户数量
	TotalUsedStorage int64 `json:"total_used_storage"` // 总存储空间占用 (Bytes)
}
