package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/ypb/your-personal-blog/internal/blog/model"
	"github.com/ypb/your-personal-blog/internal/blog/repo"
)

type categoryUsecase struct {
	categoryRepo repo.CategoryRepo
}

func NewCategoryUsecase(cr repo.CategoryRepo) Category {
	return &categoryUsecase{
		categoryRepo: cr,
	}
}

func (u *categoryUsecase) Create(ctx context.Context, in *CategoryCreateInput) (string, error) {
	cat := &model.Category{
		ID:   fmt.Sprintf("cat_%d", time.Now().UnixNano()),
		Name: in.Name,
		Slug: in.Slug,
	}
	err := u.categoryRepo.Create(ctx, cat)
	return cat.ID, err
}

func (u *categoryUsecase) Get(ctx context.Context, id string) (*model.Category, error) {
	return u.categoryRepo.GetByID(ctx, id)
}

func (u *categoryUsecase) Update(ctx context.Context, in *CategoryUpdateInput) error {
	cat, err := u.categoryRepo.GetByID(ctx, in.ID)
	if err != nil {
		return err
	}
	cat.Name = in.Name
	cat.Slug = in.Slug
	return u.categoryRepo.Update(ctx, cat)
}

func (u *categoryUsecase) Delete(ctx context.Context, id string) error {
	return u.categoryRepo.Delete(ctx, id)
}

func (u *categoryUsecase) List(ctx context.Context) ([]*CategoryDetailOutput, error) {
	cats, err := u.categoryRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	outputs := make([]*CategoryDetailOutput, 0, len(cats))
	for _, c := range cats {
		outputs = append(outputs, &CategoryDetailOutput{
			ID:   c.ID,
			Name: c.Name,
			Slug: c.Slug,
			// PostCount 逻辑后续根据业务需求在 Repo 层增加统计支持
		})
	}
	return outputs, nil
}
