package category

import (
	"context"
	"database/sql"
	"errors"
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
	_, err := s.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS categories (
			id   TEXT PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			slug TEXT NOT NULL UNIQUE
		)
	`)
	return err
}

func (s *SQLiteStore) Create(ctx context.Context, c *Category) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO categories (id, name, slug) VALUES (?, ?, ?)`,
		c.ID, c.Name, c.Slug,
	)
	return err
}

func (s *SQLiteStore) Get(ctx context.Context, id string) (*Category, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, name, slug FROM categories WHERE id = ?`, id,
	)
	c := &Category{}
	if err := row.Scan(&c.ID, &c.Name, &c.Slug); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return c, nil
}

func (s *SQLiteStore) GetBySlug(ctx context.Context, slug string) (*Category, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, name, slug FROM categories WHERE slug = ?`, slug,
	)
	c := &Category{}
	if err := row.Scan(&c.ID, &c.Name, &c.Slug); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return c, nil
}

func (s *SQLiteStore) List(ctx context.Context) ([]*Category, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, name, slug FROM categories ORDER BY name`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*Category
	for rows.Next() {
		c := &Category{}
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, rows.Err()
}

func (s *SQLiteStore) Update(ctx context.Context, c *Category) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.db.ExecContext(ctx,
		`UPDATE categories SET name = ?, slug = ? WHERE id = ?`,
		c.Name, c.Slug, c.ID,
	)
	return err
}

func (s *SQLiteStore) ListWithCount(ctx context.Context) ([]*CategoryWithCount, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT c.id, c.name, c.slug, COUNT(p.id) as post_count
		FROM categories c
		LEFT JOIN posts p ON p.category_id = c.id AND p.published = 1
		GROUP BY c.id
		ORDER BY post_count DESC, c.name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*CategoryWithCount
	for rows.Next() {
		cc := &CategoryWithCount{}
		if err := rows.Scan(&cc.ID, &cc.Name, &cc.Slug, &cc.PostCount); err != nil {
			return nil, err
		}
		list = append(list, cc)
	}
	return list, rows.Err()
}

func (s *SQLiteStore) Delete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.db.ExecContext(ctx, `DELETE FROM categories WHERE id = ?`, id)
	return err
}
