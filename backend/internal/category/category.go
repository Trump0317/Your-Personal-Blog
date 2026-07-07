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

// Store 分类存储接口
type Store interface {
	Create(ctx context.Context, c *Category) error
	Get(ctx context.Context, id string) (*Category, error)
	GetBySlug(ctx context.Context, slug string) (*Category, error)
	List(ctx context.Context) ([]*Category, error)
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

	slug := toSlug(name)
	c := &Category{
		ID:   uuid.New().String(),
		Name: name,
		Slug: slug,
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

func (s *Service) Update(ctx context.Context, id, name string) error {
	c, err := s.store.Get(ctx, id)
	if err != nil {
		return err
	}
	c.Name = strings.TrimSpace(name)
	c.Slug = toSlug(c.Name)
	return s.store.Update(ctx, c)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.store.Delete(ctx, id)
}

func toSlug(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, " ", "-"))
}
