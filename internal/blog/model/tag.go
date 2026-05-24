package model

import "time"

// Tag 标签模型
type Tag struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	PostCount int       `json:"post_count"` // 关联的文章数量
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
