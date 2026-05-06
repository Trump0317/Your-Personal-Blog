package output

import (
	"time"

	"github.com/ypb/your-personal-blog/internal/blog/model"
)

type PostDetail struct {
	ID          int64            `json:"id"`
	Title       string           `json:"title"`
	Slug        string           `json:"slug"`
	Content     string           `json:"content"`      // Markdown 原文 (用于编辑)
	HTMLContent string           `json:"html_content"` // 渲染后的 HTML (用于展示)
	Summary     string           `json:"summary"`
	Status      model.PostStatus `json:"status"`
	IsTop       bool             `json:"is_top"`
	ViewCount   int              `json:"view_count"`
	CategoryID  int64            `json:"category_id"`
	Category    *model.Category  `json:"category,omitempty"`
	Tags        []model.Tag      `json:"tags"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	PublishedAt *time.Time       `json:"published_at,omitempty"`

	// 导航信息 (可选，用于前端展示)
	PrevPost *model.Post `json:"prev_post,omitempty"`
	NextPost *model.Post `json:"next_post,omitempty"`
}
