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

func (u *tagUsecase) List(ctx context.Context) ([]*TagDetailOutput, error) {
	tags, err := u.tagRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	outputs := make([]*TagDetailOutput, 0, len(tags))
	for _, t := range tags {
		outputs = append(outputs, &TagDetailOutput{
			ID:   t.ID,
			Name: t.Name,
			Slug: t.Slug,
		})
	}
	return outputs, nil
}

func (u *tagUsecase) Delete(ctx context.Context, id string) error {
	return u.tagRepo.Delete(ctx, id)
}

func (u *tagUsecase) Create(ctx context.Context, in *TagCreateInput) error {
	// 1. 查重
	existing, err := u.tagRepo.GetByName(ctx, in.Name)
	if err == nil && existing != nil {
		return fmt.Errorf("tag already exists")
	}

	// 2. 构造模型
	tag := &model.Tag{
		ID:        fmt.Sprintf("tag_%d", time.Now().UnixNano()),
		Name:      in.Name,
		Slug:      in.Slug,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if tag.Slug == "" {
		tag.Slug = strings.ToLower(tag.Name)
	}

	return u.tagRepo.Create(ctx, tag)
}

func (u *tagUsecase) ListByIDs(ctx context.Context, ids []string) ([]*TagDetailOutput, error) {
	if len(ids) == 0 {
		return []*TagDetailOutput{}, nil
	}
	tags, err := u.tagRepo.BatchGetByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	outputs := make([]*TagDetailOutput, 0, len(tags))
	for _, t := range tags {
		outputs = append(outputs, &TagDetailOutput{
			ID:   t.ID,
			Name: t.Name,
			Slug: t.Slug,
		})
	}
	return outputs, nil
}
