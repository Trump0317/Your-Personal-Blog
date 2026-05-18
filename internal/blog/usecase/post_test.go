package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ypb/your-personal-blog/internal/blog/model"
	"github.com/ypb/your-personal-blog/internal/blog/repo"
)

// MockPostRepo 是 PostRepo 的 Mock 实现
type MockPostRepo struct {
	mock.Mock
}

func (m *MockPostRepo) Create(ctx context.Context, post *model.Post) error {
	args := m.Called(ctx, post)
	return args.Error(0)
}

func (m *MockPostRepo) GetByID(ctx context.Context, id string) (*model.Post, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Post), args.Error(1)
}

func (m *MockPostRepo) GetBySlug(ctx context.Context, slug string) (*model.Post, error) {
	args := m.Called(ctx, slug)
	return args.Get(0).(*model.Post), args.Error(1)
}

func (m *MockPostRepo) Update(ctx context.Context, post *model.Post) error {
	args := m.Called(ctx, post)
	return args.Error(0)
}

func (m *MockPostRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPostRepo) UpdateStatus(ctx context.Context, id string, status model.PostStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockPostRepo) SetTags(ctx context.Context, postID string, tagIDs []string) error {
	args := m.Called(ctx, postID, tagIDs)
	return args.Error(0)
}

func (m *MockPostRepo) List(ctx context.Context, filter *repo.PostFilter) ([]*model.Post, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*model.Post), int64(args.Int(1)), args.Error(2)
}

// MockTagUsecase 是 Tag Usecase 的 Mock 实现
type MockTagUsecase struct {
	mock.Mock
}

func (m *MockTagUsecase) List(ctx context.Context) ([]*model.Tag, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*model.Tag), args.Error(1)
}

func (m *MockTagUsecase) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTagUsecase) GetOrCreates(ctx context.Context, names []string) ([]string, error) {
	args := m.Called(ctx, names)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockTagUsecase) ListByIDs(ctx context.Context, ids []string) ([]*model.Tag, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]*model.Tag), args.Error(1)
}

func (m *MockTagUsecase) GetTagsMapByIDs(ctx context.Context, ids []string) (map[string]*model.Tag, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).(map[string]*model.Tag), args.Error(1)
}

// MockCategoryRepo 是 CategoryRepo 的 Mock 实现
type MockCategoryRepo struct {
	mock.Mock
}

func (m *MockCategoryRepo) Create(ctx context.Context, category *model.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryRepo) GetByID(ctx context.Context, id string) (*model.Category, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Category), args.Error(1)
}

func (m *MockCategoryRepo) GetBySlug(ctx context.Context, slug string) (*model.Category, error) {
	args := m.Called(ctx, slug)
	return args.Get(0).(*model.Category), args.Error(1)
}

func (m *MockCategoryRepo) List(ctx context.Context) ([]*model.Category, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*model.Category), args.Error(1)
}

func (m *MockCategoryRepo) Update(ctx context.Context, category *model.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestPostUsecase_Create(t *testing.T) {
	mockRepo := new(MockPostRepo)
	mockTagUC := new(MockTagUsecase)
	mockCatRepo := new(MockCategoryRepo)

	tuc := NewPostUsecase(mockRepo, mockCatRepo, mockTagUC)

	t.Run("成功创建文章并关联标签", func(t *testing.T) {
		ctx := context.Background()
		input := &PostCreateInput{
			Title:      "单元测试",
			Slug:       "unit-test",
			Content:    "内容",
			CategoryID: "cat_1",
			TagIDs:     []string{"tag_1", "tag_2"},
		}

		// 预期 1: 调用 PostRepo 保存文章本体
		mockRepo.On("Create", ctx, mock.MatchedBy(func(p *model.Post) bool {
			return p.Title == input.Title && p.CategoryID == input.CategoryID
		})).Return(nil)

		// 预期 2: 调用 PostRepo.SetTags 保存所有标签关系
		mockRepo.On("SetTags", ctx, mock.Anything, mock.MatchedBy(func(ids []string) bool {
			return len(ids) == 2 && contains(ids, "tag_1") && contains(ids, "tag_2")
		})).Return(nil)

		id, err := tuc.Create(ctx, input)
		assert.NoError(t, err)
		assert.NotEmpty(t, id)

		mockRepo.AssertExpectations(t)
	})
}

func TestPostUsecase_Get(t *testing.T) {
	mockRepo := new(MockPostRepo)
	mockTagUC := new(MockTagUsecase)
	mockCatRepo := new(MockCategoryRepo)
	tuc := NewPostUsecase(mockRepo, mockCatRepo, mockTagUC)

	t.Run("成功获取文章详情并填充关联项", func(t *testing.T) {
		ctx := context.Background()
		postID := "post_1"
		mockPost := &model.Post{
			ID:         postID,
			Title:      "测试文章",
			CategoryID: "cat_1",
			TagIDs:     []string{"tag_1", "tag_2"},
		}

		mockRepo.On("GetBySlug", ctx, postID).Return((*model.Post)(nil), assert.AnError)
		mockRepo.On("GetByID", ctx, postID).Return(mockPost, nil)
		mockCatRepo.On("GetByID", ctx, "cat_1").Return(&model.Category{ID: "cat_1", Name: "分类1"}, nil)
		mockTagUC.On("ListByIDs", ctx, mockPost.TagIDs).Return([]*model.Tag{
			{ID: "tag_1", Name: "标签1"},
			{ID: "tag_2", Name: "标签2"},
		}, nil)

		out, err := tuc.Get(ctx, postID)
		assert.NoError(t, err)
		assert.Equal(t, "测试文章", out.Title)
		assert.Equal(t, "分类1", out.Category.Name)
		assert.Len(t, out.Tags, 2)
	})
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
