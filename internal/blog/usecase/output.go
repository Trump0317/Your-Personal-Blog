package usecase

import (
	"time"

	"github.com/ypb/your-personal-blog/internal/blog/model"
)

type PostDetailOutput struct {
	ID          string                `json:"id"`
	Title       string                `json:"title"`
	Slug        string                `json:"slug"`
	Content     string                `json:"content"`      // Markdown 原文 (用于编辑)
	HTMLContent string                `json:"html_content"` // 渲染后的 HTML (用于展示)
	Summary     string                `json:"summary"`
	Status      model.PostStatus      `json:"status"`
	ViewCount   int                   `json:"view_count"`
	CategoryID  string                `json:"category_id"`
	Category    *CategoryDetailOutput `json:"category,omitempty"`
	Tags        []TagDetailOutput     `json:"tags"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
	PublishedAt *time.Time            `json:"published_at,omitempty"`
}

type PostListOutput struct {
	Posts []*PostItem `json:"posts"`
	Total int64       `json:"total"`
}

type PostItem struct {
	ID          string                `json:"id"`
	Title       string                `json:"title"`
	Slug        string                `json:"slug"`
	Summary     string                `json:"summary"`
	Status      model.PostStatus      `json:"status"`
	ViewCount   int                   `json:"view_count"`
	Category    *CategoryDetailOutput `json:"category,omitempty"`
	Tags        []TagDetailOutput     `json:"tags"`
	CreatedAt   time.Time             `json:"created_at"`
	PublishedAt *time.Time            `json:"published_at,omitempty"`
}

type CategoryDetailOutput struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	PostCount int    `json:"post_count"`
}

type TagDetailOutput struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	PostCount int    `json:"post_count"`
}
