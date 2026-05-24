package model

import "time"

// Category 分类模型
type Category struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Priority  int       `json:"priority"`   // 优先级，数值越小优先展示
	PostCount int       `json:"post_count"` // 关联的文章数量
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
