package usecase

import (
	"context"

	"github.com/ypb/your-personal-blog/internal/blog/model"
)

type Post interface {
	// 管理端 (Admin & CLI) 常用接口
	Create(ctx context.Context, in *PostCreateInput) (string, error)
	Update(ctx context.Context, in *PostUpdateInput) error
	Delete(ctx context.Context, id string) error
	Publish(ctx context.Context, id string) error
	Unpublish(ctx context.Context, id string) error

	// 展示端 & 列表分页接口
	Get(ctx context.Context, idOrSlug string) (*model.Post, error)
	List(ctx context.Context, in *PostListInput) ([]*model.Post, int64, error)
}

// Category 分类业务接口
type Category interface {
	Create(ctx context.Context, name, slug string) (string, error)
	Get(ctx context.Context, id string) (*model.Category, error)
	Update(ctx context.Context, id, name, slug string) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*model.Category, error)
}

// Tag 标签业务接口
type Tag interface {
	// List 获取所有标签列表
	List(ctx context.Context) ([]*model.Tag, error)
	// Delete 删除标签
	Delete(ctx context.Context, id string) error
	// GetOrCreates 处理“选择或新建”逻辑：输入混合名称，返回归一化后的 ID 列表
	GetOrCreates(ctx context.Context, names []string) ([]string, error)
	// ListByIDs 批量获取标签详情，用于文章详情展示时的对象填充
	ListByIDs(ctx context.Context, ids []string) ([]*model.Tag, error)
}
