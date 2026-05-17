package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ypb/your-personal-blog/internal/hub/controller"
	"github.com/ypb/your-personal-blog/internal/hub/usercase"
)

// MockFileUseCase 是对 usercase.File 的 mock 实现
type MockFileUseCase struct {
	mock.Mock
}

func (m *MockFileUseCase) Upload(ctx context.Context, in usercase.FileUploadInput) (*usercase.FileUploadOutput, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usercase.FileUploadOutput), args.Error(1)
}

func (m *MockFileUseCase) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockFileUseCase) GetByID(ctx context.Context, id string) (*usercase.FileDetailOutput, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usercase.FileDetailOutput), args.Error(1)
}

func (m *MockFileUseCase) ListAll(ctx context.Context, in usercase.FileListInput) ([]*usercase.FileDetailOutput, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*usercase.FileDetailOutput), args.Error(1)
}

// MockUserUseCase 是对 usercase.User 的 mock 实现
type MockUserUseCase struct {
	mock.Mock
}

func (m *MockUserUseCase) Create(ctx context.Context, in usercase.UserCreateInput) (*usercase.UserCreateOutput, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usercase.UserCreateOutput), args.Error(1)
}

func (m *MockUserUseCase) GetByAPIKey(ctx context.Context, apiKey string) (*usercase.UserDetailOutput, error) {
	args := m.Called(ctx, apiKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usercase.UserDetailOutput), args.Error(1)
}

func (m *MockUserUseCase) Delete(ctx context.Context, apiKey string) error {
	args := m.Called(ctx, apiKey)
	return args.Error(0)
}

func (m *MockUserUseCase) Update(ctx context.Context, apiKey string, in usercase.UserUpdateInput) error {
	args := m.Called(ctx, apiKey, in)
	return args.Error(0)
}

func setupTestRouter(fileUC usercase.File, userUC usercase.User) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	v1 := r.Group("/api/v1")
	controller.NewUserRouter(v1, nil, fileUC, userUC)
	return r
}

func TestUserHandler_Register(t *testing.T) {
	mockFile := new(MockFileUseCase)
	mockUser := new(MockUserUseCase)
	r := setupTestRouter(mockFile, mockUser)

	apiKey := "test-api-key"
	in := usercase.UserCreateInput{
		APIKey:       apiKey,
		InitialQuota: 1024 * 1024 * 100,
	}
	mockUser.On("Create", mock.Anything, in).Return(&usercase.UserCreateOutput{APIKey: apiKey}, nil)

	body, _ := json.Marshal(apiKey)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/user/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp usercase.UserCreateOutput
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, apiKey, resp.APIKey)
}

func TestUserHandler_Upload(t *testing.T) {
	mockFile := new(MockFileUseCase)
	mockUser := new(MockUserUseCase)
	r := setupTestRouter(mockFile, mockUser)

	apiKey := "test-api-key"

	// 构造 multipart/form-data 请求
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.txt")
	part.Write([]byte("hello world"))
	writer.WriteField("usage", "avatar")
	writer.Close()

	mockUser.On("GetByAPIKey", mock.Anything, apiKey).Return(&usercase.UserDetailOutput{APIKey: apiKey, Status: 0}, nil)
	mockFile.On("Upload", mock.Anything, mock.MatchedBy(func(in usercase.FileUploadInput) bool {
		return in.FileName == "test.txt" && in.Usage == "avatar"
	})).Return(&usercase.FileUploadOutput{ID: "file-123"}, nil)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/file/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-API-Key", apiKey)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp usercase.FileUploadOutput
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "file-123", resp.ID)
}

func TestUserHandler_ListAll(t *testing.T) {
	mockFile := new(MockFileUseCase)
	mockUser := new(MockUserUseCase)
	r := setupTestRouter(mockFile, mockUser)

	apiKey := "test-api-key"
	expectedFiles := []*usercase.FileDetailOutput{
		{ID: "1", FileName: "f1.txt", CreatedAt: time.Now()},
	}

	mockUser.On("GetByAPIKey", mock.Anything, apiKey).Return(&usercase.UserDetailOutput{APIKey: apiKey, Status: 0}, nil)
	mockFile.On("ListAll", mock.Anything, mock.MatchedBy(func(in usercase.FileListInput) bool {
		return in.APIKey == apiKey && in.Page == 1 && in.Size == 10
	})).Return(expectedFiles, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/file/all?page=1&pageSize=10", nil)
	req.Header.Set("X-API-Key", apiKey)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp []*usercase.FileDetailOutput
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Len(t, resp, 1)
	assert.Equal(t, "1", resp[0].ID)
}

func TestUserHandler_GetByID(t *testing.T) {
	mockFile := new(MockFileUseCase)
	mockUser := new(MockUserUseCase)
	r := setupTestRouter(mockFile, mockUser)

	fileID := "file-123"
	mockFile.On("GetByID", mock.Anything, fileID).Return(&usercase.FileDetailOutput{
		ID: fileID, FileName: "test.png",
	}, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/file/"+fileID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp usercase.FileDetailOutput
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, fileID, resp.ID)
}
