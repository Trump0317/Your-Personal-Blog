package usercase_test

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ypb/your-personal-blog/internal/hub/model"
	"github.com/ypb/your-personal-blog/internal/hub/usercase"
)

// MockFileStore 是对 repo.FileStore 的 mock 实现
type MockFileStore struct {
	mock.Mock
}

func (m *MockFileStore) Upload(ctx context.Context, id string, expires time.Duration, reader io.Reader) error {
	args := m.Called(ctx, id, expires, reader)
	return args.Error(0)
}

func (m *MockFileStore) Download(ctx context.Context, id string, expires time.Duration) (io.ReadCloser, error) {
	args := m.Called(ctx, id, expires)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *MockFileStore) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockFileRepo 是对 repo.FileRepo 的 mock 实现
type MockFileRepo struct {
	mock.Mock
}

func (m *MockFileRepo) Create(ctx context.Context, file model.File) (string, error) {
	args := m.Called(ctx, file)
	return args.String(0), args.Error(1)
}

func (m *MockFileRepo) GetByID(ctx context.Context, id string) (*model.File, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.File), args.Error(1)
}

func (m *MockFileRepo) ListByUser(ctx context.Context, apikey string, limit, offset int) ([]*model.File, error) {
	args := m.Called(ctx, apikey, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.File), args.Error(1)
}

func (m *MockFileRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockFileRepo) UpdateStatus(ctx context.Context, id string, status model.FileStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func TestFileUseCase_Upload(t *testing.T) {
	mockStore := new(MockFileStore)
	mockRepo := new(MockFileRepo)
	uc := usercase.NewFileUseCase(mockStore, mockRepo)

	ctx := context.Background()
	content := bytes.NewBufferString("hello world")
	in := usercase.FileUploadInput{
		APIKey:   "user-1",
		FileName: "test.txt",
		Size:     11,
		Usage:    "post_cover",
		Content:  content,
	}
	userID := "user-1"

	// 预期行为
	mockStore.On("Upload", ctx, mock.AnythingOfType("string"), mock.Anything, content).Return(nil)
	mockRepo.On("Create", ctx, mock.MatchedBy(func(f model.File) bool {
		return f.FileName == "test.txt" && f.Uploader == userID && f.Usage == "post_cover"
	})).Return("file-123", nil)

	// 执行
	out, err := uc.Upload(ctx, in)

	// 断言
	assert.NoError(t, err)
	assert.NotNil(t, out)
	assert.Equal(t, "file-123", out.ID)

	mockStore.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestFileUseCase_Delete(t *testing.T) {
	mockStore := new(MockFileStore)
	mockRepo := new(MockFileRepo)
	uc := usercase.NewFileUseCase(mockStore, mockRepo)

	ctx := context.Background()
	fileID := "file-123"
	storagePath := "2026-05-10/test.txt"

	// 预期行为
	// 1. 先查数据库获取路径
	mockRepo.On("GetByID", ctx, fileID).Return(&model.File{
		ID:          fileID,
		StoragePath: storagePath,
	}, nil)

	// 2. 删除物理存储
	mockStore.On("Delete", ctx, storagePath).Return(nil)

	// 3. 删除数据库记录
	mockRepo.On("Delete", ctx, fileID).Return(nil)

	// 执行
	err := uc.Delete(ctx, fileID)

	// 断言
	assert.NoError(t, err)

	mockStore.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestFileUseCase_ListAll(t *testing.T) {
	mockStore := new(MockFileStore)
	mockRepo := new(MockFileRepo)
	uc := usercase.NewFileUseCase(mockStore, mockRepo)

	ctx := context.Background()
	apiKey := "user-1"
	in := usercase.FileListInput{
		APIKey: apiKey,
		Page:   2,
		Size:   10,
	}

	mockFiles := []*model.File{
		{ID: "f1", FileName: "file1.png", Uploader: apiKey, Usage: "avatar"},
		{ID: "f2", FileName: "file2.jpg", Uploader: apiKey, Usage: "post"},
	}

	// 预期行为：page=2, size=10 -> limit=10, offset=10
	mockRepo.On("ListByUser", ctx, apiKey, 10, 10).Return(mockFiles, nil)

	// 执行
	out, err := uc.ListAll(ctx, in)

	// 断言
	assert.NoError(t, err)
	assert.Len(t, out, 2)
	assert.Equal(t, "f1", out[0].ID)
	assert.Equal(t, "avatar", out[0].Usage)

	mockRepo.AssertExpectations(t)
}
