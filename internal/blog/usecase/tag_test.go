package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ypb/your-personal-blog/internal/blog/model"
)

type MockTagRepo struct {
	mock.Mock
}

func (m *MockTagRepo) Create(ctx context.Context, tag *model.Tag) error {
	args := m.Called(ctx, tag)
	if tag.ID == "" {
		tag.ID = "new_generated_id"
	}
	return args.Error(0)
}

func (m *MockTagRepo) GetByID(ctx context.Context, id string) (*model.Tag, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Tag), args.Error(1)
}

func (m *MockTagRepo) GetByName(ctx context.Context, name string) (*model.Tag, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Tag), args.Error(1)
}

func (m *MockTagRepo) List(ctx context.Context) ([]*model.Tag, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Tag), args.Error(1)
}

func (m *MockTagRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTagRepo) BatchGetByIDs(ctx context.Context, ids []string) ([]*model.Tag, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Tag), args.Error(1)
}

func (m *MockTagRepo) GetBySlug(ctx context.Context, slug string) (*model.Tag, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Tag), args.Error(1)
}

func TestTagUsecase_Create(t *testing.T) {
	mockRepo := new(MockTagRepo)
	uc := NewTagUsecase(mockRepo)

	t.Run("Create success", func(t *testing.T) {
		ctx := context.Background()
		in := &TagCreateInput{Name: "Go", Slug: "go-lang"}
		mockRepo.On("GetByName", ctx, "Go").Return(nil, nil)
		mockRepo.On("Create", ctx, mock.Anything).Return(nil)
		err := uc.Create(ctx, in)
		assert.NoError(t, err)
	})
}

func TestTagUsecase_List(t *testing.T) {
	mockRepo := new(MockTagRepo)
	uc := NewTagUsecase(mockRepo)
	t.Run("List SUCCESS", func(t *testing.T) {
		ctx := context.Background()
		mockRepo.On("List", ctx).Return([]*model.Tag{{ID: "1", Name: "Go"}}, nil)
		res, err := uc.List(ctx)
		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})
}
