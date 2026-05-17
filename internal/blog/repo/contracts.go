package repo

import (
	"context"

	"github.com/ypb/your-personal-blog/internal/blog/model"
)

// PostRepo 文章存储接口：只负责文章本体及与其强关系的分类/标签管理
type PostRepo interface {
	Create(ctx context.Context, post *model.Post) error
	GetByID(ctx context.Context, id string) (*model.Post, error)
	GetBySlug(ctx context.Context, slug string) (*model.Post, error)
	Update(ctx context.Context, post *model.Post) error
	Delete(ctx context.Context, id string) error

	// 状态与关联维护
	UpdateStatus(ctx context.Context, id string, status model.PostStatus) error
	// SetTags 重置文章与标签的关联（多对多中间表操作）
	SetTags(ctx context.Context, postID string, tagIDs []string) error

	// 查询接口
	List(ctx context.Context, filter *PostFilter) ([]*model.Post, int64, error)
}

// PostFilter 文章列表查询过滤器
type PostFilter struct {
	Offset     int
	Limit      int
	Status     *model.PostStatus
	CategoryID string
	TagID      string
	Query      string // 关键词搜索
}

// CategoryRepo 分类字典接口：负责分类元数据的维护
type CategoryRepo interface {
	Create(ctx context.Context, category *model.Category) error
	GetByID(ctx context.Context, id string) (*model.Category, error)
	GetBySlug(ctx context.Context, slug string) (*model.Category, error)
	List(ctx context.Context) ([]*model.Category, error)
	Update(ctx context.Context, category *model.Category) error
	Delete(ctx context.Context, id string) error
}

// TagRepo 标签字典接口：负责标签池的维护
type TagRepo interface {
	Save(ctx context.Context, tag *model.Tag) error
	GetByID(ctx context.Context, id string) (*model.Tag, error)
	GetByName(ctx context.Context, name string) (*model.Tag, error)
	List(ctx context.Context) ([]*model.Tag, error)
	// BatchGetByIDs 用于文章列表展示时，一次性补全标签对象
	BatchGetByIDs(ctx context.Context, ids []string) ([]*model.Tag, error)
	Delete(ctx context.Context, id string) error
}
