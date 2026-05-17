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

func (u *categoryUsecase) Create(ctx context.Context, name, slug string) (string, error) {
	cat := &model.Category{
		ID:   fmt.Sprintf("cat_%d", time.Now().UnixNano()),
		Name: name,
		Slug: slug,
	}
	err := u.categoryRepo.Create(ctx, cat)
	return cat.ID, err
}

func (u *categoryUsecase) Get(ctx context.Context, id string) (*model.Category, error) {
	return u.categoryRepo.GetByID(ctx, id)
}

func (u *categoryUsecase) Update(ctx context.Context, id, name, slug string) error {
	cat, err := u.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	cat.Name = name
	cat.Slug = slug
	return u.categoryRepo.Update(ctx, cat)
}

func (u *categoryUsecase) Delete(ctx context.Context, id string) error {
	return u.categoryRepo.Delete(ctx, id)
}

func (u *categoryUsecase) List(ctx context.Context) ([]*model.Category, error) {
	return u.categoryRepo.List(ctx)
}
