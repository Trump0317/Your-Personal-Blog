package repo

import (
	"context"

	"github.com/ypb/your-personal-blog/internal/blog/model"
)

// PostRepo 文章存储接口
type PostRepo interface {
	Create(ctx context.Context, post *model.Post) error
	GetByID(ctx context.Context, id int64) (*model.Post, error)
	GetBySlug(ctx context.Context, slug string) (*model.Post, error)
	Update(ctx context.Context, post *model.Post) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, offset, limit int, status *model.PostStatus) ([]*model.Post, int64, error)
}

// TaxonomyRepo 分类与标签存储接口
type TaxonomyRepo interface {
	CreateCategory(ctx context.Context, category *model.Category) error
	GetCategoryByID(ctx context.Context, id int64) (*model.Category, error)
	ListCategories(ctx context.Context) ([]*model.Category, error)

	CreateTag(ctx context.Context, tag *model.Tag) error
	ListTags(ctx context.Context) ([]*model.Tag, error)
}
