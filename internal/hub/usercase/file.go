package usercase

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/ypb/your-personal-blog/internal/hub/model"
	"github.com/ypb/your-personal-blog/internal/hub/repo"
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

func (u *fileUseCase) Upload(ctx context.Context, in FileUploadInput) (*FileUploadOutput, error) {
	// 1. 生成存储路径
	ext := filepath.Ext(in.FileName)
	saveName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	subDir := time.Now().Format("2006-01-02")
	storagePath := filepath.Join(subDir, saveName)

	// 2. 准备数据库记录
	fileRecord := model.File{
		Uploader:    in.APIKey,
		FileName:    in.FileName,
		FileSize:    in.Size,
		Usage:       in.Usage,
		StoragePath: storagePath,
		Status:      model.FileActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
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

	return &FileUploadOutput{
		ID: id,
	}, nil
}

func (u *fileUseCase) GetByID(ctx context.Context, id string) (*FileDetailOutput, error) {
	file, err := u.dbRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabase, err)
	}

	return &FileDetailOutput{
		ID:        file.ID,
		FileName:  file.FileName,
		FileSize:  file.FileSize,
		MimeType:  file.MimeType,
		Usage:     file.Usage,
		CreatedAt: file.CreatedAt,
	}, nil
}

func (u *fileUseCase) ListAll(ctx context.Context, in FileListInput) ([]*FileDetailOutput, error) {
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.Size <= 0 {
		in.Size = 10
	}
	offset := (in.Page - 1) * in.Size

	files, err := u.dbRepo.ListByUser(ctx, in.APIKey, in.Size, offset)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabase, err)
	}

	out := make([]*FileDetailOutput, 0, len(files))
	for _, f := range files {
		out = append(out, &FileDetailOutput{
			ID:        f.ID,
			FileName:  f.FileName,
			FileSize:  f.FileSize,
			MimeType:  f.MimeType,
			Usage:     f.Usage,
			CreatedAt: f.CreatedAt,
		})
	}
	return out, nil
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
