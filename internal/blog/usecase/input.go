package usecase

import "github.com/ypb/your-personal-blog/internal/blog/model"

type PostCreateInput struct {
	Title      string           `json:"title" validate:"required"`
	Slug       string           `json:"slug" validate:"required,slug"`
	Content    string           `json:"content" validate:"required"`
	Summary    string           `json:"summary"`
	Status     model.PostStatus `json:"status"`
	CategoryID string           `json:"category_id" validate:"required"`
	TagIDs     []string         `json:"tag_ids"` // 标签 ID 列表
}

type PostUpdateInput struct {
	ID         string            `json:"id" validate:"required"`
	Title      *string           `json:"title,omitempty"`
	Slug       *string           `json:"slug,omitempty" validate:"omitempty,slug"`
	Content    *string           `json:"content,omitempty"`
	Summary    *string           `json:"summary,omitempty"`
	Status     *model.PostStatus `json:"status,omitempty"`
	CategoryID *string           `json:"category_id,omitempty"`
	TagIDs     *[]string         `json:"tag_ids,omitempty"` // 更新后的全量标签 ID
}

type PostListInput struct {
	Page       int               `json:"page" validate:"gte=1"`
	PageSize   int               `json:"page_size" validate:"gte=1,lte=100"`
	Status     *model.PostStatus `json:"status,omitempty"`
	CategoryID *string           `json:"category_id,omitempty"`
	TagID      *string           `json:"tag_id,omitempty"`
	Query      string            `json:"query,omitempty"` // 搜索关键词
}
