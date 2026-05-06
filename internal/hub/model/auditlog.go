package model

// 审计日志模型 AuditLog
// ID : 唯一标识日志的ID
// UserID : 相关用户的ID
// Action : 执行的操作类型（如：文件上传、文件下载、用户登录等）
// Resource : 相关资源的描述（如：文件名、IP地址等）
// Timestamp : 操作发生的时间戳

const (
	ActionFileUpload   = "file_upload"
	ActionFileDownload = "file_download"
	ActionUserLogin    = "user_login"
	ActionUserLogout   = "user_logout"
	ActionQuotaChange  = "quota_change"
)

type AuditLog struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Action    string `json:"action"`
	Resource  string `json:"resource"`
	Timestamp int64  `json:"timestamp"`
}
