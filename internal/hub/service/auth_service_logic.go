package service

import (
	"context"

	"github.com/ypb/your-personal-blog/internal/hub/model"
	"github.com/ypb/your-personal-blog/internal/hub/repository"
)

// authService 实现了 AuthService 接口
// 负责处理身份验证和授权相关的业务逻辑
type authService struct {
	clientRepo repository.ClientMetadataRepository
}

func NewAuthService(repo repository.ClientMetadataRepository) AuthService {
	return &authService{clientRepo: repo}
}

// Authenticate 校验 API Key 是否有效。
func (s *authService) Authenticate(ctx context.Context, apiKey string) (*model.ClientMetadata, error) {
	return s.clientRepo.GetByAPIKey(ctx, apiKey)
}

// Authorize TODO: 校验特定操作的权限，如 RBAC 或简单动作校验 (业务/Auth)
func (s *authService) Authorize(ctx context.Context, client *model.ClientMetadata, action string) (bool, error) {
	return true, nil
}

// CheckQuota TODO: 检查客户端是否有足够的配额进行当前操作，需聚合当前使用量 vs 额度 (业务/Auth)
func (s *authService) CheckQuota(ctx context.Context, clientID string, additionalSize int64) (bool, error) {
	return true, nil
}
