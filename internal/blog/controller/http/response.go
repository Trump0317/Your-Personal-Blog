package http

import "time"

// Response 统一 HTTP 响应结构
type Response struct {
	Code    int         `json:"code"`    // 业务状态码
	Message string      `json:"message"` // 提示信息
	Data    interface{} `json:"data"`    // 数据
}

// PostItemResponse 访客端：文章列表条目
type PostItemResponse struct {
	ID          string           `json:"id"`
	Title       string           `json:"title"`
	Slug        string           `json:"slug"`
	Summary     string           `json:"summary"`
	ViewCount   int              `json:"view_count"`
	Category    *CategoryBenefit `json:"category,omitempty"`
	Tags        []TagBenefit     `json:"tags,omitempty"`
	PublishedAt *time.Time       `json:"published_at"`
}

// PostDetailResponse 访客端：文章详情响应
type PostDetailResponse struct {
	ID          string           `json:"id"`
	Title       string           `json:"title"`
	Slug        string           `json:"slug"`
	Content     string           `json:"content"`      // Markdown 原文
	HTMLContent string           `json:"html_content"` // 渲染后的 HTML
	Summary     string           `json:"summary"`
	ViewCount   int              `json:"view_count"`
	Category    *CategoryBenefit `json:"category,omitempty"`
	Tags        []TagBenefit     `json:"tags,omitempty"`
	PublishedAt *time.Time       `json:"published_at"`
}

// PostListResponse 访客端：文章列表分页响应
type PostListResponse struct {
	Total int64               `json:"total"`
	Items []*PostItemResponse `json:"items"`
}

// CategoryBenefit 简化的分类信息
type CategoryBenefit struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// TagBenefit 简化的标签信息
type TagBenefit struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// TagDetailResponse 管理端：标签详情响应
type TagDetailResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	PostCount int    `json:"post_count"`
}

// CategoryDetailResponse 管理端：分类详情响应
type CategoryDetailResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	PostCount int    `json:"post_count"`
}

// PostAdminItemResponse 管理端：文章列表条目
type PostAdminItemResponse struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Status      int        `json:"status"`
	ViewCount   int        `json:"view_count"`
	CategoryID  string     `json:"category_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	PublishedAt *time.Time `json:"published_at"`
}
