package model

type ClientStatus int

const (
	ClientActive    ClientStatus = iota // 0: 活跃
	ClientInactive                      // 1: 不活跃
	ClientSuspended                     // 2: 暂停
)

// 客户元数据 ClientMetadata
// ClientID : 唯一标识租户的ID
// APIKey : 用于认证和授权的API密钥
// QuotaLimit : 配额限制，单位为字节
// CurrentUsage : 当前使用量，单位为字节

type ClientMetadata struct {
	ClientID     string       `json:"client_id"`
	APIKey       string       `json:"api_key"`
	QuotaLimit   int64        `json:"quota_limit"`
	CurrentUsage int64        `json:"current_usage"`
	Status       ClientStatus `json:"status"`
}
