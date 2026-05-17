package model

import "time"

// PostStatus 文章状态
type PostStatus int

const (
	PostDraft     PostStatus = iota // 草稿
	PostPublished                   // 已发布
	PostDeleted                     // 已删除
)

// Post 文章模型
type Post struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`        // 标题
	Slug        string     `json:"slug"`         // URL 别名
	Summary     string     `json:"summary"`      // 摘要
	Content     string     `json:"content"`      // Markdown 内容
	HTMLContent string     `json:"html_content"` // 渲染后的 HTML
	ViewCount   int        `json:"view_count"`   // 阅读量
	CategoryID  string     `json:"category_id"`
	TagIDs      []string   `json:"tag_ids"`
	Tags        []*Tag     `json:"tags,omitempty"` // 填充后的完整对象，用于前端展示
	Status      PostStatus `json:"status"`         // 状态
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	PublishedAt *time.Time `json:"published_at"` // 发布时间
}
