package model

type UserStatus int

const (
	UserActive   UserStatus = iota // 0: 活跃
	UserDisabled                   // 1: 禁用
	UserDeleted                    // 2: 已删除
)

// 用户模型 User
// ID : 唯一标识用户的ID
// APIKey : 用于认证和授权的API密钥
// QuotaLimit : 配额限制，单位为字节
// CurrentUsage : 当前使用量，单位为字节

type User struct {
	ID           string     `json:"id"`
	APIKey       string     `json:"api_key"`
	QuotaLimit   int64      `json:"quota_limit"`
	CurrentUsage int64      `json:"current_usage"`
	Status       UserStatus `json:"status"`
	CreatedAt    int64      `json:"created_at"`
	UpdatedAt    int64      `json:"updated_at"`
}
