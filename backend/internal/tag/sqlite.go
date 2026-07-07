package tag

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"sync"
)

// compile-time check
var _ Store = (*SQLiteStore)(nil)

type SQLiteStore struct {
	db *sql.DB
	mu sync.Mutex // SQLite 单写锁
}

func NewSQLiteStore(db *sql.DB) *SQLiteStore {
	return &SQLiteStore{db: db}
}

func (s *SQLiteStore) Migrate(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS tags (
			id   TEXT PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			slug TEXT NOT NULL UNIQUE
		)
	`)
	return err
}

func (s *SQLiteStore) Create(ctx context.Context, t *Tag) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO tags (id, name, slug) VALUES (?, ?, ?)`,
		t.ID, t.Name, t.Slug,
	)
	return err
}

func (s *SQLiteStore) Get(ctx context.Context, id string) (*Tag, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, name, slug FROM tags WHERE id = ?`, id)
	t := &Tag{}
	if err := row.Scan(&t.ID, &t.Name, &t.Slug); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("tag not found")
		}
		return nil, err
	}
	return t, nil
}

func (s *SQLiteStore) GetByName(ctx context.Context, name string) (*Tag, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, name, slug FROM tags WHERE name = ?`, name,
	)
	t := &Tag{}
	if err := row.Scan(&t.ID, &t.Name, &t.Slug); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("tag not found")
		}
		return nil, err
	}
	return t, nil
}

func (s *SQLiteStore) List(ctx context.Context) ([]*Tag, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, name, slug FROM tags ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*Tag
	for rows.Next() {
		t := &Tag{}
		if err := rows.Scan(&t.ID, &t.Name, &t.Slug); err != nil {
			return nil, err
		}
		list = append(list, t)
	}
	return list, rows.Err()
}

func (s *SQLiteStore) GetByIDs(ctx context.Context, ids []string) ([]*Tag, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	placeholders := strings.Repeat("?,", len(ids)-1) + "?"
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, name, slug FROM tags WHERE id IN (`+placeholders+`)`, args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*Tag
	for rows.Next() {
		t := &Tag{}
		if err := rows.Scan(&t.ID, &t.Name, &t.Slug); err != nil {
			return nil, err
		}
		list = append(list, t)
	}
	return list, rows.Err()
}

func (s *SQLiteStore) Delete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.db.ExecContext(ctx, `DELETE FROM tags WHERE id = ?`, id)
	return err
}
