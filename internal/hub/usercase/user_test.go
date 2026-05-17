package usercase_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ypb/your-personal-blog/internal/hub/model"
	"github.com/ypb/your-personal-blog/internal/hub/usercase"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Create(ctx context.Context, user model.User) (string, error) {
	args := m.Called(ctx, user)
	return args.String(0), args.Error(1)
}

func (m *MockUserRepo) GetByAPIKey(ctx context.Context, apiKey string) (*model.User, error) {
	args := m.Called(ctx, apiKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepo) Delete(ctx context.Context, apiKey string) error {
	args := m.Called(ctx, apiKey)
	return args.Error(0)
}

func (m *MockUserRepo) Update(ctx context.Context, apiKey string, user model.User) error {
	args := m.Called(ctx, apiKey, user)
	return args.Error(0)
}

func (m *MockUserRepo) UpdateStatus(ctx context.Context, apiKey string, status model.UserStatus) error {
	args := m.Called(ctx, apiKey, status)
	return args.Error(0)
}

func (m *MockUserRepo) UpdateUsage(ctx context.Context, apiKey string, delta int64) error {
	args := m.Called(ctx, apiKey, delta)
	return args.Error(0)
}

func TestUserUseCase_Create(t *testing.T) {
	mockRepo := new(MockUserRepo)
	uc := usercase.NewUserUseCase(mockRepo)
	ctx := context.Background()

	t.Run("success with specified apiKey", func(t *testing.T) {
		in := usercase.UserCreateInput{APIKey: "test-api-key"}
		mockRepo.On("Create", ctx, mock.MatchedBy(func(u model.User) bool {
			return u.APIKey == "test-api-key"
		})).Return("test-api-key", nil).Once()
		out, err := uc.Create(ctx, in)
		assert.NoError(t, err)
		assert.Equal(t, "test-api-key", out.APIKey)
	})
}

func TestUserUseCase_GetByAPIKey(t *testing.T) {
	mockRepo := new(MockUserRepo)
	uc := usercase.NewUserUseCase(mockRepo)
	ctx := context.Background()
	apiKey := "test-key"

	t.Run("success", func(t *testing.T) {
		now := time.Now()
		mockUser := &model.User{
			APIKey:    apiKey,
			Status:    model.UserActive,
			CreatedAt: now,
			UpdatedAt: now,
		}
		mockRepo.On("GetByAPIKey", ctx, apiKey).Return(mockUser, nil).Once()
		out, err := uc.GetByAPIKey(ctx, apiKey)
		assert.NoError(t, err)
		assert.Equal(t, apiKey, out.APIKey)
	})
}

func TestUserUseCase_Delete(t *testing.T) {
	mockRepo := new(MockUserRepo)
	uc := usercase.NewUserUseCase(mockRepo)
	ctx := context.Background()
	apiKey := "test-key"

	mockRepo.On("Delete", ctx, apiKey).Return(nil).Once()
	err := uc.Delete(ctx, apiKey)
	assert.NoError(t, err)
}
