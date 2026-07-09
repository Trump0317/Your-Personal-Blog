package category

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
)

// Category 分类
type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// CategoryWithCount 带文章数的分类
type CategoryWithCount struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	PostCount int    `json:"post_count"`
}

// Store 分类存储接口
type Store interface {
	Create(ctx context.Context, c *Category) error
	Get(ctx context.Context, id string) (*Category, error)
	GetBySlug(ctx context.Context, slug string) (*Category, error)
	List(ctx context.Context) ([]*Category, error)
	ListWithCount(ctx context.Context) ([]*CategoryWithCount, error)
	Update(ctx context.Context, c *Category) error
	Delete(ctx context.Context, id string) error
}

// Service 分类业务逻辑
type Service struct {
	store Store
}

func NewService(s Store) *Service {
	return &Service{store: s}
}

func (s *Service) Create(ctx context.Context, name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", errors.New("category name required")
	}

	id := uuid.New().String()
	c := &Category{
		ID:   id,
		Name: name,
		Slug: toSlug(name, id[:8]),
	}
	if err := s.store.Create(ctx, c); err != nil {
		return "", err
	}
	return c.ID, nil
}

func (s *Service) Get(ctx context.Context, id string) (*Category, error) {
	return s.store.Get(ctx, id)
}

func (s *Service) GetBySlug(ctx context.Context, slug string) (*Category, error) {
	return s.store.GetBySlug(ctx, slug)
}

func (s *Service) List(ctx context.Context) ([]*Category, error) {
	return s.store.List(ctx)
}

func (s *Service) ListWithCount(ctx context.Context) ([]*CategoryWithCount, error) {
	return s.store.ListWithCount(ctx)
}

func (s *Service) Update(ctx context.Context, id, name string) error {
	c, err := s.store.Get(ctx, id)
	if err != nil {
		return err
	}
	c.Name = strings.TrimSpace(name)
	c.Slug = toSlug(c.Name, c.ID[:8])
	return s.store.Update(ctx, c)
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
