package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ypb/your-personal-blog/internal/blog/model"
	"github.com/ypb/your-personal-blog/internal/blog/repo"
)

type postRepo struct {
	db *sql.DB
}

func NewPostRepo(db *sql.DB) repo.PostRepo {
	return &postRepo{db: db}
}

func (r *postRepo) Create(ctx context.Context, p *model.Post) error {
	query := `INSERT INTO posts (title, slug, content, html_content, summary, cover_image, status, is_top, view_count, category_id, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	now := time.Now()
	res, err := r.db.ExecContext(ctx, query,
		p.Title, p.Slug, p.Content, p.HTMLContent, p.Summary, p.CoverImage, p.Status, p.IsTop, p.ViewCount, p.CategoryID, now, now)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	p.ID = id
	p.CreatedAt = now
	p.UpdatedAt = now
	return nil
}

func (r *postRepo) GetByID(ctx context.Context, id int64) (*model.Post, error) {
	p := &model.Post{}
	query := `SELECT id, title, slug, content, html_content, summary, cover_image, status, is_top, view_count, category_id, created_at, updated_at FROM posts WHERE id = ?`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.Title, &p.Slug, &p.Content, &p.HTMLContent, &p.Summary, &p.CoverImage, &p.Status, &p.IsTop, &p.ViewCount, &p.CategoryID, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("post not found")
	}
	return p, err
}

func (r *postRepo) GetBySlug(ctx context.Context, slug string) (*model.Post, error) {
	p := &model.Post{}
	query := `SELECT id, title, slug, content, html_content, summary, cover_image, status, is_top, view_count, category_id, created_at, updated_at FROM posts WHERE slug = ?`
	err := r.db.QueryRowContext(ctx, query, slug).Scan(
		&p.ID, &p.Title, &p.Slug, &p.Content, &p.HTMLContent, &p.Summary, &p.CoverImage, &p.Status, &p.IsTop, &p.ViewCount, &p.CategoryID, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("post not found")
	}
	return p, err
}

func (r *postRepo) Update(ctx context.Context, p *model.Post) error {
	query := `UPDATE posts SET title=?, slug=?, content=?, html_content=?, summary=?, cover_image=?, status=?, is_top=?, view_count=?, category_id=?, updated_at=? WHERE id=?`
	now := time.Now()
	_, err := r.db.ExecContext(ctx, query,
		p.Title, p.Slug, p.Content, p.HTMLContent, p.Summary, p.CoverImage, p.Status, p.IsTop, p.ViewCount, p.CategoryID, now, p.ID)
	if err == nil {
		p.UpdatedAt = now
	}
	return err
}

func (r *postRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM posts WHERE id = ?", id)
	return err
}

func (r *postRepo) List(ctx context.Context, offset, limit int, status *model.PostStatus) ([]*model.Post, int64, error) {
	var total int64
	countQuery := "SELECT count(*) FROM posts"
	listQuery := `SELECT id, title, slug, content, html_content, summary, cover_image, status, is_top, view_count, category_id, created_at, updated_at FROM posts`

	args := []interface{}{}
	if status != nil {
		countQuery += " WHERE status = ?"
		listQuery += " WHERE status = ?"
		args = append(args, *status)
	}

	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	listQuery += " ORDER BY is_top DESC, created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var posts []*model.Post
	for rows.Next() {
		p := &model.Post{}
		err := rows.Scan(&p.ID, &p.Title, &p.Slug, &p.Content, &p.HTMLContent, &p.Summary, &p.CoverImage, &p.Status, &p.IsTop, &p.ViewCount, &p.CategoryID, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		posts = append(posts, p)
	}
	return posts, total, nil
}
