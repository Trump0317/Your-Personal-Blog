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

// Store 标签存储接口
type Store interface {
	Create(ctx context.Context, t *Tag) error
	Get(ctx context.Context, id string) (*Tag, error)
	GetByName(ctx context.Context, name string) (*Tag, error)
	List(ctx context.Context) ([]*Tag, error)
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

	slug := toSlug(name)
	t := &Tag{
		ID:   uuid.New().String(),
		Name: name,
		Slug: slug,
	}
	if err := s.store.Create(ctx, t); err != nil {
		return "", err
	}
	return t.ID, nil
}

func (s *Service) List(ctx context.Context) ([]*Tag, error) {
	return s.store.List(ctx)
}

func (s *Service) GetByIDs(ctx context.Context, ids []string) ([]*Tag, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	return s.store.GetByIDs(ctx, ids)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.store.Delete(ctx, id)
}

func toSlug(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, " ", "-"))
}
