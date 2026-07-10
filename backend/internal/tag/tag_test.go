package tag

import (
	"context"
	"errors"
	"testing"
)

type mockTagStore struct {
	tags map[string]*Tag
}

func (m *mockTagStore) Create(ctx context.Context, t *Tag) error {
	if _, ok := m.tags[t.ID]; ok {
		return errors.New("duplicate")
	}
	m.tags[t.ID] = t
	return nil
}
func (m *mockTagStore) Get(ctx context.Context, id string) (*Tag, error) {
	tag, ok := m.tags[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return tag, nil
}
func (m *mockTagStore) GetByName(ctx context.Context, name string) (*Tag, error) {
	for _, t := range m.tags {
		if t.Name == name {
			return t, nil
		}
	}
	return nil, errors.New("not found")
}
func (m *mockTagStore) GetBySlug(ctx context.Context, slug string) (*Tag, error) {
	for _, t := range m.tags {
		if t.Slug == slug {
			return t, nil
		}
	}
	return nil, errors.New("not found")
}
func (m *mockTagStore) List(ctx context.Context) ([]*Tag, error) {
	var list []*Tag
	for _, t := range m.tags {
		list = append(list, t)
	}
	return list, nil
}
func (m *mockTagStore) ListWithCount(ctx context.Context) ([]*TagWithCount, error) { return nil, nil }
func (m *mockTagStore) GetByIDs(ctx context.Context, ids []string) ([]*Tag, error) {
	var list []*Tag
	for _, id := range ids {
		if t, ok := m.tags[id]; ok {
			list = append(list, t)
		}
	}
	return list, nil
}
func (m *mockTagStore) Update(ctx context.Context, t *Tag) error {
	if _, ok := m.tags[t.ID]; !ok {
		return errors.New("not found")
	}
	m.tags[t.ID] = t
	return nil
}
func (m *mockTagStore) Delete(ctx context.Context, id string) error {
	delete(m.tags, id)
	return nil
}

func TestTagService_CreateOrGet(t *testing.T) {
	svc := NewService(&mockTagStore{tags: map[string]*Tag{}})

	// first create
	id1, err := svc.CreateOrGet(context.Background(), "  go  ")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if id1 == "" {
		t.Fatal("expected non-empty id")
	}

	// second call should return same id (dedup)
	id2, err := svc.CreateOrGet(context.Background(), "go")
	if err != nil {
		t.Fatalf("get existing: %v", err)
	}
	if id1 != id2 {
		t.Errorf("expected same id, got %s and %s", id1, id2)
	}
}

func TestTagService_CreateOrGet_EmptyName(t *testing.T) {
	svc := NewService(&mockTagStore{tags: map[string]*Tag{}})
	_, err := svc.CreateOrGet(context.Background(), "   ")
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestTagService_Update(t *testing.T) {
	store := &mockTagStore{tags: map[string]*Tag{}}
	svc := NewService(store)

	id, _ := svc.CreateOrGet(context.Background(), "old-tag")
	if err := svc.Update(context.Background(), id, "new-tag"); err != nil {
		t.Fatalf("update: %v", err)
	}

	tag, _ := store.Get(context.Background(), id)
	if tag.Name != "new-tag" {
		t.Errorf("name = %q, want 'new-tag'", tag.Name)
	}
}

func TestTagService_Update_EmptyName(t *testing.T) {
	store := &mockTagStore{tags: map[string]*Tag{}}
	svc := NewService(store)

	id, _ := svc.CreateOrGet(context.Background(), "tag")
	err := svc.Update(context.Background(), id, "   ")
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestTagService_Delete(t *testing.T) {
	store := &mockTagStore{tags: map[string]*Tag{}}
	svc := NewService(store)

	id, _ := svc.CreateOrGet(context.Background(), "to-delete")
	if err := svc.Delete(context.Background(), id); err != nil {
		t.Fatalf("delete: %v", err)
	}
	_, err := store.Get(context.Background(), id)
	if err == nil {
		t.Fatal("expected error after delete")
	}
}
