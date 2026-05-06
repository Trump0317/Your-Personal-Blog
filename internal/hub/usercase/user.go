package usercase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ypb/your-personal-blog/internal/hub/model"
	"github.com/ypb/your-personal-blog/internal/hub/repo"
	"github.com/ypb/your-personal-blog/internal/hub/usercase/input"
	"github.com/ypb/your-personal-blog/internal/hub/usercase/output"
)

type UserUseCase struct {
	dbRepo repo.UserRepo // 负责数据库 SQL 操作
}

func NewUserUseCase(dbRepo repo.UserRepo) User {
	return &UserUseCase{
		dbRepo: dbRepo,
	}
}

func (u *UserUseCase) Create(ctx context.Context, in input.UserCreate) (*output.UserDetail, error) {
	apiKey := in.APIKey
	if apiKey == "" {
		apiKey = uuid.New().String()
	}

	userRecord := model.User{
		ID:           uuid.New().String(),
		APIKey:       apiKey,
		QuotaLimit:   1024 * 1024 * 100, // 默认 100MB
		CurrentUsage: 0,
		Status:       model.UserActive,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}

	id, err := u.dbRepo.Create(ctx, userRecord)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &output.UserDetail{
		ID:     id,
		APIKey: apiKey,
	}, nil
}

func (u *UserUseCase) GetByID(ctx context.Context, id string) (*output.UserDetail, error) {
	user, err := u.dbRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return &output.UserDetail{
		ID:     user.ID,
		APIKey: user.APIKey,
	}, nil
}

func (u *UserUseCase) Delete(ctx context.Context, id string) error {
	return u.dbRepo.Delete(ctx, id)
}

func (u *UserUseCase) GetByAPIKey(ctx context.Context, apiKey string) (*output.UserDetail, error) {
	user, err := u.dbRepo.GetByAPIKey(ctx, apiKey)
	if err != nil {
		return nil, fmt.Errorf("user not found by APIKey: %w", err)
	}

	return &output.UserDetail{
		ID:     user.ID,
		APIKey: user.APIKey,
	}, nil
}
