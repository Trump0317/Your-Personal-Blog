package usercase

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/ypb/your-personal-blog/internal/hub/model"
	"github.com/ypb/your-personal-blog/internal/hub/repo"
	"github.com/ypb/your-personal-blog/internal/hub/usercase/input"
	"github.com/ypb/your-personal-blog/internal/hub/usercase/output"
)

type fileUseCase struct {
	storageRepo repo.FileStore // 负责物理存储
	dbRepo      repo.FileRepo  // 负责数据库 SQL 操作
}

// NewFileUseCase 创建 FileUseCase 实例。
func NewFileUseCase(storageRepo repo.FileStore, dbRepo repo.FileRepo) File {
	return &fileUseCase{
		storageRepo: storageRepo,
		dbRepo:      dbRepo,
	}
}

func (u *fileUseCase) Upload(ctx context.Context, userID string, in input.FileUpload) (*output.FileDetail, error) {
	// 1. 生成存储路径
	ext := filepath.Ext(in.FileName)
	saveName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	subDir := time.Now().Format("2006-01-02")
	storagePath := filepath.Join(subDir, saveName)

	// 2. 准备数据库记录
	fileRecord := model.File{
		UploaderID:   userID,
		OriginalName: in.FileName,
		FileSize:     in.Size,
		StoragePath:  storagePath,
		StorageType:  "local",
		Status:       model.FileActive,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}

	// 3. 执行物理存储上传 (设置默认过期时间 0 表示永久)
	if err := u.storageRepo.Upload(ctx, storagePath, 0, in.Content); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrStorage, err)
	}

	// 4. 保存到数据库
	id, err := u.dbRepo.Create(ctx, fileRecord)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabase, err)
	}

	return &output.FileDetail{
		ID:        id,
		FileName:  in.FileName,
		FileSize:  in.Size,
		Usage:     in.Usage,
		CreatedAt: time.Now(),
	}, nil
}

func (u *fileUseCase) GetByID(ctx context.Context, expires time.Duration, id string) (*output.FileDetail, error) {
	file, err := u.dbRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabase, err)
	}

	return &output.FileDetail{
		ID:        file.ID,
		FileName:  file.OriginalName,
		FileSize:  file.FileSize,
		Usage:     file.StorageType,
		CreatedAt: time.Unix(file.CreatedAt, 0),
	}, nil
}

func (u *fileUseCase) ListAll(ctx context.Context, userID string) ([]*output.FileDetail, error) {
	files, err := u.dbRepo.ListByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabase, err)
	}

	var result []*output.FileDetail
	for _, f := range files {
		result = append(result, &output.FileDetail{
			ID:        f.ID,
			FileName:  f.OriginalName,
			FileSize:  f.FileSize,
			Usage:     f.StorageType,
			CreatedAt: time.Unix(f.CreatedAt, 0),
		})
	}
	return result, nil
}

func (u *fileUseCase) Delete(ctx context.Context, id string) error {
	// 1. 获取文件元数据
	file, err := u.dbRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabase, err)
	}

	// 2. 删除物理文件
	if err := u.storageRepo.Delete(ctx, file.StoragePath); err != nil {
		return fmt.Errorf("%w: %v", ErrStorage, err)
	}

	// 3. 删除数据库记录
	if err := u.dbRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%w: %v", ErrDatabase, err)
	}

	return nil
}
