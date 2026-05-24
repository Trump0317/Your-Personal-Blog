package http

// PostListRequest 访客端：文章列表查询
type PostListRequest struct {
	Page       int    `form:"page" binding:"omitempty,min=1"`
	PageSize   int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	CategoryID string `form:"category_id"`
	TagID      string `form:"tag_id"`
	Query      string `form:"q"`
}

// LoginRequest 管理端：登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// CreatePostRequest 管理端：创建文章请求
type CreatePostRequest struct {
	Title      string   `json:"title" binding:"required"`
	Slug       string   `json:"slug"`
	Content    string   `json:"content" binding:"required"`
	Summary    string   `json:"summary"`
	CategoryID string   `json:"category_id"`
	Status     string   `json:"status"` // "published" 或 "draft"
	TagIDs     []string `json:"tag_ids"`
}

// UpdatePostRequest 管理端：更新文章请求
type UpdatePostRequest struct {
	Title      *string   `json:"title"`
	Slug       *string   `json:"slug"`
	Content    *string   `json:"content"`
	Summary    *string   `json:"summary"`
	CategoryID *string   `json:"category_id"`
	Status     *string   `json:"status"`
	TagIDs     *[]string `json:"tag_ids"`
}

// CreateTagRequest 管理端：创建标签请求
type CreateTagRequest struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug"`
}

// CreateCategoryRequest 管理端：创建分类请求
type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug"`
}

// UpdateCategoryRequest 管理端：更新分类请求
type UpdateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug"`
}
