package category

import (
	"context"
	"errors"
	"testing"
)

type mockCatStore struct {
	cats map[string]*Category
}

func (m *mockCatStore) Create(ctx context.Context, c *Category) error {
	if _, ok := m.cats[c.ID]; ok {
		return errors.New("duplicate")
	}
	m.cats[c.ID] = c
	return nil
}
func (m *mockCatStore) Get(ctx context.Context, id string) (*Category, error) {
	c, ok := m.cats[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return c, nil
}
func (m *mockCatStore) GetBySlug(ctx context.Context, slug string) (*Category, error) {
	for _, c := range m.cats {
		if c.Slug == slug {
			return c, nil
		}
	}
	return nil, errors.New("not found")
}
func (m *mockCatStore) List(ctx context.Context) ([]*Category, error) {
	var list []*Category
	for _, c := range m.cats {
		list = append(list, c)
	}
	return list, nil
}
func (m *mockCatStore) ListWithCount(ctx context.Context) ([]*CategoryWithCount, error) { return nil, nil }
func (m *mockCatStore) Update(ctx context.Context, c *Category) error {
	if _, ok := m.cats[c.ID]; !ok {
		return errors.New("not found")
	}
	m.cats[c.ID] = c
	return nil
}
func (m *mockCatStore) Delete(ctx context.Context, id string) error {
	delete(m.cats, id)
	return nil
}

func TestCategoryService_Create(t *testing.T) {
	svc := NewService(&mockCatStore{cats: map[string]*Category{}})

	id, err := svc.Create(context.Background(), "  Go 语言  ")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if id == "" {
		t.Fatal("expected non-empty id")
	}

	c, _ := svc.Get(context.Background(), id)
	if c.Name != "Go 语言" {
		t.Errorf("name = %q, want 'Go 语言'", c.Name)
	}
	if c.Slug == "" {
		t.Error("expected auto-generated slug")
	}
}

func TestCategoryService_Create_EmptyName(t *testing.T) {
	svc := NewService(&mockCatStore{cats: map[string]*Category{}})
	_, err := svc.Create(context.Background(), "   ")
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestCategoryService_Update(t *testing.T) {
	store := &mockCatStore{cats: map[string]*Category{}}
	svc := NewService(store)

	id, _ := svc.Create(context.Background(), "Old")
	if err := svc.Update(context.Background(), id, "New"); err != nil {
		t.Fatalf("update: %v", err)
	}

	c, _ := svc.Get(context.Background(), id)
	if c.Name != "New" {
		t.Errorf("name = %q, want 'New'", c.Name)
	}
}

func TestCategoryService_Delete(t *testing.T) {
	store := &mockCatStore{cats: map[string]*Category{}}
	svc := NewService(store)

	id, _ := svc.Create(context.Background(), "ToDelete")
	if err := svc.Delete(context.Background(), id); err != nil {
		t.Fatalf("delete: %v", err)
	}
	_, err := svc.Get(context.Background(), id)
	if err == nil {
		t.Fatal("expected error after delete")
	}
}
