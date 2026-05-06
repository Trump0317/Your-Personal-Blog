package input

import "github.com/ypb/your-personal-blog/internal/blog/model"

type PostCreate struct {
	Title      string  `json:"title" validate:"required"`
	Slug       string  `json:"slug" validate:"required,slug"`
	Content    string  `json:"content" validate:"required"`
	Summary    string  `json:"summary"`
	IsTop      bool    `json:"is_top"`
	CategoryID int64   `json:"category_id" validate:"required"`
	TagIDs     []int64 `json:"tag_ids"`
}

type PostUpdate struct {
	ID         int64             `json:"id" validate:"required"`
	Title      *string           `json:"title,omitempty"`
	Slug       *string           `json:"slug,omitempty" validate:"omitempty,slug"`
	Content    *string           `json:"content,omitempty"`
	Summary    *string           `json:"summary,omitempty"`
	IsTop      *bool             `json:"is_top,omitempty"`
	Status     *model.PostStatus `json:"status,omitempty"`
	CategoryID *int64            `json:"category_id,omitempty"`
	TagIDs     *[]int64          `json:"tag_ids,omitempty"`
}

type PostList struct {
	Page       int               `json:"page" validate:"gte=1"`
	PageSize   int               `json:"page_size" validate:"gte=1,lte=100"`
	Status     *model.PostStatus `json:"status,omitempty"`
	CategoryID *int64            `json:"category_id,omitempty"`
	TagID      *int64            `json:"tag_id,omitempty"`
	Query      string            `json:"query,omitempty"` // 搜索关键词
}
