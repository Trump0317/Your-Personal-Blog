package post

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"sync"
)

var _ Store = (*SQLiteStore)(nil)

type SQLiteStore struct {
	db *sql.DB
	mu sync.Mutex
}

func NewSQLiteStore(db *sql.DB) *SQLiteStore {
	return &SQLiteStore{db: db}
}

func (s *SQLiteStore) Migrate(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS posts (
			id           TEXT PRIMARY KEY,
			title        TEXT NOT NULL,
			slug         TEXT NOT NULL UNIQUE,
			content      TEXT NOT NULL DEFAULT '',
			html_content TEXT NOT NULL DEFAULT '',
			summary      TEXT NOT NULL DEFAULT '',
			category_id  TEXT NOT NULL DEFAULT '',
			published    INTEGER NOT NULL DEFAULT 0,
			created_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS post_tags (
			post_id TEXT NOT NULL,
			tag_id  TEXT NOT NULL,
			PRIMARY KEY (post_id, tag_id)
		)
	`)
	return err
}

func (s *SQLiteStore) Create(ctx context.Context, p *Post) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO posts (id, title, slug, content, html_content, summary, category_id, published, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		p.ID, p.Title, p.Slug, p.Content, p.HTMLContent, p.Summary, p.CategoryID, boolToInt(p.Published), p.CreatedAt, p.UpdatedAt,
	)
	return err
}

func (s *SQLiteStore) Get(ctx context.Context, idOrSlug string) (*Post, error) {
	const q = `SELECT id, title, slug, content, html_content, summary, category_id, published, created_at, updated_at FROM posts WHERE id = ? OR slug = ?`
	row := s.db.QueryRowContext(ctx, q, idOrSlug, idOrSlug)
	return scanPost(row)
}

func (s *SQLiteStore) List(ctx context.Context, f Filter) ([]*Post, int, error) {
	var conditions []string
	var args []interface{}

	if f.CategoryID != "" {
		conditions = append(conditions, "category_id = ?")
		args = append(args, f.CategoryID)
	}
	if f.Published != nil {
		conditions = append(conditions, "published = ?")
		args = append(args, boolToInt(*f.Published))
	}
	if f.Query != "" {
		conditions = append(conditions, "(title LIKE ? OR content LIKE ?)")
		q := "%" + f.Query + "%"
		args = append(args, q, q)
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	// count
	var total int
	countQ := "SELECT COUNT(*) FROM posts " + where
	if err := s.db.QueryRowContext(ctx, countQ, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// list
	limit := f.Limit
	if limit <= 0 {
		limit = 10
	}
	offset := f.Offset
	if offset < 0 {
		offset = 0
	}
	listQ := "SELECT id, title, slug, content, html_content, summary, category_id, published, created_at, updated_at FROM posts " +
		where + " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	listArgs := append(args, limit, offset)
	rows, err := s.db.QueryContext(ctx, listQ, listArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		p, err := scanPostFromRows(rows)
		if err != nil {
			return nil, 0, err
		}
		// 填充 tag_ids
		tagIDs, _ := s.getTagIDs(ctx, p.ID)
		p.TagIDs = tagIDs
		posts = append(posts, p)
	}
	return posts, total, rows.Err()
}

func (s *SQLiteStore) Update(ctx context.Context, p *Post) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.ExecContext(ctx, `
		UPDATE posts SET title=?, slug=?, content=?, html_content=?, summary=?, category_id=?, published=?, updated_at=?
		WHERE id=?
	`,
		p.Title, p.Slug, p.Content, p.HTMLContent, p.Summary, p.CategoryID, boolToInt(p.Published), p.UpdatedAt, p.ID,
	)
	return err
}

func (s *SQLiteStore) Delete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, _ = s.db.ExecContext(ctx, `DELETE FROM post_tags WHERE post_id = ?`, id)
	_, err := s.db.ExecContext(ctx, `DELETE FROM posts WHERE id = ?`, id)
	return err
}

func (s *SQLiteStore) SetTags(ctx context.Context, postID string, tagIDs []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, _ = s.db.ExecContext(ctx, `DELETE FROM post_tags WHERE post_id = ?`, postID)
	for _, tid := range tagIDs {
		_, err := s.db.ExecContext(ctx,
			`INSERT OR IGNORE INTO post_tags (post_id, tag_id) VALUES (?, ?)`,
			postID, tid,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SQLiteStore) getTagIDs(ctx context.Context, postID string) ([]string, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT tag_id FROM post_tags WHERE post_id = ?`, postID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func scanPost(row *sql.Row) (*Post, error) {
	p := &Post{}
	var published int
	err := row.Scan(&p.ID, &p.Title, &p.Slug, &p.Content, &p.HTMLContent,
		&p.Summary, &p.CategoryID, &published, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}
	p.Published = published != 0
	return p, nil
}

func scanPostFromRows(rows *sql.Rows) (*Post, error) {
	p := &Post{}
	var published int
	err := rows.Scan(&p.ID, &p.Title, &p.Slug, &p.Content, &p.HTMLContent,
		&p.Summary, &p.CategoryID, &published, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	p.Published = published != 0
	return p, nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
