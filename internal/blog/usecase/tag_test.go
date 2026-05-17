package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ypb/your-personal-blog/internal/blog/model"
)

// MockTagRepo 是 TagRepo 的 Mock 实现
type MockTagRepo struct {
	mock.Mock
}

func (m *MockTagRepo) Save(ctx context.Context, tag *model.Tag) error {
	args := m.Called(ctx, tag)
	if tag.ID == "" {
		tag.ID = "new_generated_id" // 模拟数据库生成 ID
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
	return args.Get(0).([]*model.Tag), args.Error(1)
}

func (m *MockTagRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTagRepo) BatchGetByIDs(ctx context.Context, ids []string) ([]*model.Tag, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]*model.Tag), args.Error(1)
}

func (m *MockTagRepo) GetBySlug(ctx context.Context, slug string) (*model.Tag, error) {
	args := m.Called(ctx, slug)
	return args.Get(0).(*model.Tag), args.Error(1)
}

func TestTagUsecase_GetOrCreates(t *testing.T) {
	mockRepo := new(MockTagRepo)
	uc := NewTagUsecase(mockRepo)

	t.Run("处理混合标签：已有的匹配，没有的新建", func(t *testing.T) {
		ctx := context.Background()
		names := []string{"Go ", " golang", "NewTag"}

		// 预期 1: "go" (来自 "Go ") 已存在
		mockRepo.On("GetByName", ctx, "go").Return(&model.Tag{ID: "id_go", Name: "Go"}, nil)

		// 预期 2: "golang" 不存在 -> 创建
		mockRepo.On("GetByName", ctx, "golang").Return(nil, nil)
		mockRepo.On("Save", ctx, mock.MatchedBy(func(tag *model.Tag) bool {
			return tag.Name == "golang" && tag.Slug == "golang"
		})).Return(nil)

		// 预期 3: "newtag" 不存在 -> 创建
		mockRepo.On("GetByName", ctx, "newtag").Return(nil, nil)
		mockRepo.On("Save", ctx, mock.MatchedBy(func(tag *model.Tag) bool {
			return tag.Name == "NewTag"
		})).Return(nil)

		ids, err := uc.GetOrCreates(ctx, names)

		assert.NoError(t, err)
		assert.Len(t, ids, 3)
		assert.Contains(t, ids, "id_go")
		mockRepo.AssertExpectations(t)
	})

	t.Run("大小写去重测试", func(t *testing.T) {
		ctx := context.Background()
		names := []string{"Go", "go", "GO"}

		// 尽管输入了三个，但逻辑上应该只查询一次或处理一次核心名称
		mockRepo.On("GetByName", ctx, "go").Return(&model.Tag{ID: "id_go", Name: "Go"}, nil)

		ids, err := uc.GetOrCreates(ctx, names)

		assert.NoError(t, err)
		assert.Len(t, ids, 1)
		assert.Equal(t, "id_go", ids[0])
	})
}
