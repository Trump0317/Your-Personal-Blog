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
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`         // URL 别名
	Content     string     `json:"content"`      // Markdown 内容
	HTMLContent string     `json:"html_content"` // 渲染后的 HTML
	Summary     string     `json:"summary"`      // 摘要
	CoverImage  string     `json:"cover_image"`  // 封面图
	Status      PostStatus `json:"status"`
	IsTop       bool       `json:"is_top"`     // 是否置顶
	ViewCount   int        `json:"view_count"` // 阅读量
	CategoryID  int64      `json:"category_id"`
	Tags        []Tag      `json:"tags"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	PublishedAt *time.Time `json:"published_at"` // 发布时间
}
