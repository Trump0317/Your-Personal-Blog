package model

import "time"

type UserStatus int

const (
	UserActive   UserStatus = iota // 0: 活跃
	UserDisabled                   // 1: 禁用
	UserDeleted                    // 2: 删除
)

// 用户模型 User
// ID : 唯一标识用户的ID
// APIKey : 用于认证和授权的API密钥
// QuotaLimit : 配额限制，单位为字节
// CurrentUsage : 当前使用量，单位为字节

type User struct {
	APIKey       string     `json:"api_key"`
	QuotaLimit   int64      `json:"quota_limit"`
	CurrentUsage int64      `json:"current_usage"`
	Status       UserStatus `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
