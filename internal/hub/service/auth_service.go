package service

import (
	"context"

	"github.com/ypb/your-personal-blog/internal/hub/model"
)

// AuthService 处理所有与身份验证和授权相关的业务逻辑。
// 它通过 ClientID/APIKey 识别租户，并校验其操作权限。
type AuthService interface {
	// Authenticate 校验 API Key 是否有效。
	// 参数:
	//   ctx: 上下文。
	//   apiKey: 客户端提供的 API 密钥。
	// 返回:
	//   *model.ClientMetadata: 关联的客户端信息。
	//   error: 如果 Key 无效或客户端被禁用，返回相应的业务错误。
	Authenticate(ctx context.Context, apiKey string) (*model.ClientMetadata, error)

	// Authorize 校验特定操作的权限（预留）。
	// 目前可能仅需校验 ClientStatus，未来可扩展细粒度权限。
	// 参数:
	//   client: 经过认证的客户端。
	//   action: 动作名称（如 "upload", "delete"）。
	// 返回:
	//   bool: 是否允许该操作。
	Authorize(ctx context.Context, client *model.ClientMetadata, action string) (bool, error)

	// CheckQuota 检查客户端是否有足够的配额进行当前操作。
	// 参数:
	//   clientID: 客户端 ID。
	//   additionalSize: 本次操作新增的存储字节数。
	// 返回:
	//   bool: 配额是否充足。
	//   error: 数据库查询异常。
	CheckQuota(ctx context.Context, clientID string, additionalSize int64) (bool, error)
}
