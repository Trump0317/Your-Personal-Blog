package usercase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ypb/your-personal-blog/internal/hub/model"
	"github.com/ypb/your-personal-blog/internal/hub/repo"
)

type UserUseCase struct {
	dbRepo repo.UserRepo // 负责数据库 SQL 操作
}

func NewUserUseCase(dbRepo repo.UserRepo) User {
	return &UserUseCase{
		dbRepo: dbRepo,
	}
}

func (u *UserUseCase) Create(ctx context.Context, in UserCreateInput) (*UserCreateOutput, error) {
	apiKey := in.APIKey
	if apiKey == "" {
		apiKey = uuid.New().String()
	}

	// 如果未指定初始配额，默认 100MB
	quota := in.InitialQuota
	if quota <= 0 {
		quota = 1024 * 1024 * 100
	}

	userRecord := model.User{
		APIKey:       apiKey,
		QuotaLimit:   quota,
		CurrentUsage: 0,
		Status:       model.UserActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err := u.dbRepo.Create(ctx, userRecord)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabase, err)
	}

	return &UserCreateOutput{
		APIKey: apiKey,
	}, nil
}

func (u *UserUseCase) Update(ctx context.Context, apiKey string, in UserUpdateInput) error {
	// 目前仅返回 nil，未来可在此实现用户状态更新或配额调整
	return nil
}

func (u *UserUseCase) Delete(ctx context.Context, apiKey string) error {
	if err := u.dbRepo.Delete(ctx, apiKey); err != nil {
		return fmt.Errorf("%w: %v", ErrDatabase, err)
	}
	return nil
}

func (u *UserUseCase) GetByAPIKey(ctx context.Context, apiKey string) (*UserDetailOutput, error) {
	user, err := u.dbRepo.GetByAPIKey(ctx, apiKey)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabase, err)
	}

	return &UserDetailOutput{
		APIKey:       user.APIKey,
		QuotaLimit:   user.QuotaLimit,
		CurrentUsage: user.CurrentUsage,
		Status:       int(user.Status),
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}, nil
}
