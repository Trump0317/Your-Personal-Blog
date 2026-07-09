package tag

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
)

// Tag 标签
type Tag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// TagWithCount 带文章数的标签
type TagWithCount struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	PostCount int    `json:"post_count"`
}

// Store 标签存储接口
type Store interface {
	Create(ctx context.Context, t *Tag) error
	Get(ctx context.Context, id string) (*Tag, error)
	GetByName(ctx context.Context, name string) (*Tag, error)
	GetBySlug(ctx context.Context, slug string) (*Tag, error)
	List(ctx context.Context) ([]*Tag, error)
	ListWithCount(ctx context.Context) ([]*TagWithCount, error)
	GetByIDs(ctx context.Context, ids []string) ([]*Tag, error)
	Delete(ctx context.Context, id string) error
}

// Service 标签业务逻辑
type Service struct {
	store Store
}

func NewService(s Store) *Service {
	return &Service{store: s}
}

// CreateOrGet 如果标签已存在（按 name 匹配）则返回已有 ID，否则新建
func (s *Service) CreateOrGet(ctx context.Context, name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", errors.New("tag name required")
	}

	existing, err := s.store.GetByName(ctx, name)
	if err == nil {
		return existing.ID, nil
	}

	id := uuid.New().String()
	t := &Tag{
		ID:   id,
		Name: name,
		Slug: toSlug(name, id[:8]),
	}
	if err := s.store.Create(ctx, t); err != nil {
		return "", err
	}
	return t.ID, nil
}

func (s *Service) List(ctx context.Context) ([]*Tag, error) {
	return s.store.List(ctx)
}

func (s *Service) ListWithCount(ctx context.Context) ([]*TagWithCount, error) {
	return s.store.ListWithCount(ctx)
}

func (s *Service) GetByIDs(ctx context.Context, ids []string) ([]*Tag, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	return s.store.GetByIDs(ctx, ids)
}

func (s *Service) GetBySlug(ctx context.Context, slug string) (*Tag, error) {
	return s.store.GetBySlug(ctx, slug)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.store.Delete(ctx, id)
}

func toSlug(name, suffix string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	s = strings.ReplaceAll(s, " ", "-")
	for _, r := range s {
		if r > 127 {
			return suffix
		}
	}
	return s
}
