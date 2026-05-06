package usecase

import (
	"context"

	"github.com/ypb/your-personal-blog/internal/blog/model"
	"github.com/ypb/your-personal-blog/internal/blog/usecase/input"
	"github.com/ypb/your-personal-blog/internal/blog/usecase/output"
)

type Post interface {
	// 管理端 (Admin & CLI) 常用接口
	Create(ctx context.Context, in *input.PostCreate) (int64, error)
	Update(ctx context.Context, in *input.PostUpdate) error
	Delete(ctx context.Context, id int64) error
	Publish(ctx context.Context, id int64) error
	Unpublish(ctx context.Context, id int64) error

	// 展示端 & 列表分页接口
	Get(ctx context.Context, idOrSlug string) (*output.PostDetail, error)
	List(ctx context.Context, in *input.PostList) ([]*model.Post, int64, error)
}

type Taxonomy interface {
	ListCategories(ctx context.Context) ([]*model.Category, error)
	ListTags(ctx context.Context) ([]*model.Tag, error)
}
