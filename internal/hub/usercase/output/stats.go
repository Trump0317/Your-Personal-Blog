package output

// SystemStats 系统统计输出
type SystemStats struct {
TotalFiles       int64 `json:"total_files"`
TotalUsers       int64 `json:"total_users"`
TotalUsedStorage int64 `json:"total_used_storage"`
}
