package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ypb/your-personal-blog/internal/blog/model"
	"github.com/ypb/your-personal-blog/internal/blog/repo"
)

type tagUsecase struct {
	tagRepo repo.TagRepo
}

// NewTagUsecase 创建标签业务逻辑实例
func NewTagUsecase(tr repo.TagRepo) Tag {
	return &tagUsecase{
		tagRepo: tr,
	}
}

func (u *tagUsecase) List(ctx context.Context) ([]*model.Tag, error) {
	return u.tagRepo.List(ctx)
}

func (u *tagUsecase) Delete(ctx context.Context, id string) error {
	return u.tagRepo.Delete(ctx, id)
}

func (u *tagUsecase) GetOrCreates(ctx context.Context, names []string) ([]string, error) {
	if len(names) == 0 {
		return []string{}, nil
	}

	// 1. 本地去重并整理
	uniqueNames := make(map[string]string) // key: lowerCase, value: originalCase
	for _, name := range names {
		n := strings.TrimSpace(name)
		if n != "" {
			uniqueNames[strings.ToLower(n)] = n
		}
	}

	var finalIDs []string

	// 2. 遍历处理每一个名称 (此处可进一步优化为批量查询，目前先实现核心逻辑)
	for lowerName, originalName := range uniqueNames {
		// 检查库中是否已存在（忽略大小写，Repo 层应支持 GetByName）
		existing, err := u.tagRepo.GetByName(ctx, lowerName)
		if err != nil {
			// 如果是“未找到”错误，则继续创建逻辑
			if err.Error() == "tag not found" {
				goto create
			}
			return nil, fmt.Errorf("failed to check tag existence: %w", err)
		}

		if existing != nil {
			finalIDs = append(finalIDs, existing.ID)
			continue
		}

	create:
		// 3. 库中没有，创建新标签
		newTag := &model.Tag{
			ID:        fmt.Sprintf("tag_%d", time.Now().UnixNano()), // 为内存 Repo 显式生成 ID
			Name:      originalName,
			Slug:      lowerName, // 简单处理，直接用小写名称作为 Slug
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := u.tagRepo.Save(ctx, newTag); err != nil {
			return nil, fmt.Errorf("failed to auto-create tag [%s]: %w", originalName, err)
		}
		finalIDs = append(finalIDs, newTag.ID)
	}

	return finalIDs, nil
}

func (u *tagUsecase) ListByIDs(ctx context.Context, ids []string) ([]*model.Tag, error) {
	if len(ids) == 0 {
		return []*model.Tag{}, nil
	}
	return u.tagRepo.BatchGetByIDs(ctx, ids)
}
