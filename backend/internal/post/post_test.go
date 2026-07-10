package post

import (
	"context"
	"errors"
	"testing"
)

// ── Mock Stores ──

type mockPostStore struct {
	posts  map[string]*Post
	tagIDs map[string][]string
}

func (m *mockPostStore) Create(ctx context.Context, p *Post) error {
	if _, ok := m.posts[p.ID]; ok {
		return errors.New("duplicate")
	}
	m.posts[p.ID] = p
	return nil
}
func (m *mockPostStore) Get(ctx context.Context, idOrSlug string) (*Post, error) {
	for _, p := range m.posts {
		if p.ID == idOrSlug || p.Slug == idOrSlug {
			p.TagIDs = m.tagIDs[p.ID]
			return p, nil
		}
	}
	return nil, errors.New("not found")
}
func (m *mockPostStore) List(ctx context.Context, f Filter) ([]*Post, int, error) {
	var result []*Post
	for _, p := range m.posts {
		if f.Published != nil && p.Published != *f.Published {
			continue
		}
		p.TagIDs = m.tagIDs[p.ID]
		result = append(result, p)
	}
	total := len(result)
	limit := f.Limit
	if limit <= 0 {
		limit = total
	}
	offset := f.Offset
	if offset < 0 {
		offset = 0
	}
	if offset >= total {
		return nil, total, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return result[offset:end], total, nil
}
func (m *mockPostStore) Update(ctx context.Context, p *Post) error {
	if _, ok := m.posts[p.ID]; !ok {
		return errors.New("not found")
	}
	m.posts[p.ID] = p
	return nil
}
func (m *mockPostStore) Delete(ctx context.Context, id string) error {
	delete(m.posts, id)
	delete(m.tagIDs, id)
	return nil
}
func (m *mockPostStore) SetTags(ctx context.Context, postID string, tagIDs []string) error {
	m.tagIDs[postID] = tagIDs
	return nil
}
func (m *mockPostStore) Stats(ctx context.Context) (PostStats, error) {
	var c int
	for _, p := range m.posts {
		if p.Published {
			c++
		}
	}
	return PostStats{PostCount: c}, nil
}
func (m *mockPostStore) Archive(ctx context.Context) ([]ArchiveYear, error) {
	return nil, nil
}

type mockCatStore struct {
	cats map[string]*CategoryRef
}

func (m *mockCatStore) Get(ctx context.Context, id string) (*CategoryRef, error) {
	c, ok := m.cats[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return c, nil
}
func (m *mockCatStore) List(ctx context.Context) ([]*CategoryRef, error) {
	var list []*CategoryRef
	for _, c := range m.cats {
		list = append(list, c)
	}
	return list, nil
}

type mockTagStore struct {
	tags map[string]TagRef
}

func (m *mockTagStore) GetByIDs(ctx context.Context, ids []string) ([]TagRef, error) {
	var list []TagRef
	for _, id := range ids {
		if t, ok := m.tags[id]; ok {
			list = append(list, t)
		}
	}
	return list, nil
}
func (m *mockTagStore) List(ctx context.Context) ([]TagRef, error) {
	var list []TagRef
	for _, t := range m.tags {
		list = append(list, t)
	}
	return list, nil
}

// ── Tests ──

func TestService_Create(t *testing.T) {
	svc := NewService(
		&mockPostStore{posts: map[string]*Post{}, tagIDs: map[string][]string{}},
		&mockCatStore{cats: map[string]*CategoryRef{}},
		&mockTagStore{tags: map[string]TagRef{}},
	)

	id, err := svc.Create(context.Background(), CreateInput{
		Title:     "Test Post",
		Content:   "## Hello\n\nWorld",
		Published: true,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if id == "" {
		t.Fatal("expected non-empty id")
	}

	// verify created
	p, err := svc.Get(context.Background(), id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if p.Title != "Test Post" {
		t.Errorf("title = %q, want 'Test Post'", p.Title)
	}
	if !p.Published {
		t.Error("expected published=true")
	}
	if p.HTMLContent == "" {
		t.Error("expected html_content to be rendered")
	}
	if p.Summary == "" {
		t.Error("expected summary auto-extracted")
	}
	if p.Slug == "" {
		t.Error("expected slug auto-generated")
	}
}

func TestService_Create_EmptyTitle(t *testing.T) {
	svc := NewService(
		&mockPostStore{posts: map[string]*Post{}, tagIDs: map[string][]string{}},
		&mockCatStore{cats: map[string]*CategoryRef{}},
		&mockTagStore{tags: map[string]TagRef{}},
	)
	_, err := svc.Create(context.Background(), CreateInput{Title: "   "})
	if err == nil {
		t.Fatal("expected error for empty title")
	}
}

func TestService_Get_FillsRefs(t *testing.T) {
	postStore := &mockPostStore{posts: map[string]*Post{}, tagIDs: map[string][]string{}}
	catStore := &mockCatStore{cats: map[string]*CategoryRef{
		"c1": {ID: "c1", Name: "Go", Slug: "go"},
	}}
	tagStore := &mockTagStore{tags: map[string]TagRef{
		"t1": {ID: "t1", Name: "gin", Slug: "gin"},
	}}

	svc := NewService(postStore, catStore, tagStore)

	id, _ := svc.Create(context.Background(), CreateInput{
		Title:      "With Refs",
		CategoryID: "c1",
		TagIDs:     []string{"t1"},
		Published:  true,
	})

	p, err := svc.Get(context.Background(), id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if p.Category == nil {
		t.Error("expected category to be filled")
	} else if p.Category.Name != "Go" {
		t.Errorf("category name = %q, want 'Go'", p.Category.Name)
	}
	if len(p.Tags) != 1 {
		t.Errorf("tags len = %d, want 1", len(p.Tags))
	} else if p.Tags[0].Name != "gin" {
		t.Errorf("tag name = %q, want 'gin'", p.Tags[0].Name)
	}
}

func TestService_Update(t *testing.T) {
	store := &mockPostStore{posts: map[string]*Post{}, tagIDs: map[string][]string{}}
	svc := NewService(store, &mockCatStore{cats: map[string]*CategoryRef{}}, &mockTagStore{tags: map[string]TagRef{}})

	id, _ := svc.Create(context.Background(), CreateInput{Title: "Old", Published: false})

	newTitle := "New"
	err := svc.Update(context.Background(), UpdateInput{ID: id, Title: &newTitle})
	if err != nil {
		t.Fatalf("update: %v", err)
	}

	p, _ := svc.Get(context.Background(), id)
	if p.Title != "New" {
		t.Errorf("title = %q, want 'New'", p.Title)
	}
}

func TestService_Publish_Unpublish(t *testing.T) {
	store := &mockPostStore{posts: map[string]*Post{}, tagIDs: map[string][]string{}}
	svc := NewService(store, &mockCatStore{cats: map[string]*CategoryRef{}}, &mockTagStore{tags: map[string]TagRef{}})

	id, _ := svc.Create(context.Background(), CreateInput{Title: "Draft", Published: false})

	p, _ := svc.Get(context.Background(), id)
	if p.Published {
		t.Error("expected published=false")
	}

	if err := svc.Publish(context.Background(), id); err != nil {
		t.Fatalf("publish: %v", err)
	}
	p, _ = svc.Get(context.Background(), id)
	if !p.Published {
		t.Error("expected published=true after publish")
	}

	if err := svc.Unpublish(context.Background(), id); err != nil {
		t.Fatalf("unpublish: %v", err)
	}
	p, _ = svc.Get(context.Background(), id)
	if p.Published {
		t.Error("expected published=false after unpublish")
	}
}

func TestService_Delete(t *testing.T) {
	store := &mockPostStore{posts: map[string]*Post{}, tagIDs: map[string][]string{}}
	svc := NewService(store, &mockCatStore{cats: map[string]*CategoryRef{}}, &mockTagStore{tags: map[string]TagRef{}})

	id, _ := svc.Create(context.Background(), CreateInput{Title: "ToDelete"})

	if err := svc.Delete(context.Background(), id); err != nil {
		t.Fatalf("delete: %v", err)
	}
	_, err := svc.Get(context.Background(), id)
	if err == nil {
		t.Fatal("expected error after delete")
	}
}

func TestService_List(t *testing.T) {
	store := &mockPostStore{posts: map[string]*Post{}, tagIDs: map[string][]string{}}
	svc := NewService(store, &mockCatStore{cats: map[string]*CategoryRef{}}, &mockTagStore{tags: map[string]TagRef{}})

	svc.Create(context.Background(), CreateInput{Title: "A", Published: true})
	svc.Create(context.Background(), CreateInput{Title: "B", Published: false})

	pub := true
	posts, total, _ := svc.List(context.Background(), Filter{Published: &pub})
	if total != 1 {
		t.Errorf("published total = %d, want 1", total)
	}
	if len(posts) != 1 {
		t.Errorf("published len = %d, want 1", len(posts))
	}

	// list all (no filter)
	posts, total, _ = svc.List(context.Background(), Filter{})
	if total != 2 {
		t.Errorf("unfiltered total = %d, want 2", total)
	}
}

func TestService_Slug_AutoGenerate(t *testing.T) {
	svc := NewService(
		&mockPostStore{posts: map[string]*Post{}, tagIDs: map[string][]string{}},
		&mockCatStore{cats: map[string]*CategoryRef{}},
		&mockTagStore{tags: map[string]TagRef{}},
	)
	id, err := svc.Create(context.Background(), CreateInput{Title: "Hello World"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	p, _ := svc.Get(context.Background(), id)
	if p.Slug == "" {
		t.Error("expected auto-generated slug")
	}
	// should start with "hello-world-"
	if len(p.Slug) < len("hello-world-") || p.Slug[:12] != "hello-world-" {
		t.Errorf("slug = %q, expected prefix 'hello-world-'", p.Slug)
	}
}

func TestService_Markdown_Rendering(t *testing.T) {
	svc := NewService(
		&mockPostStore{posts: map[string]*Post{}, tagIDs: map[string][]string{}},
		&mockCatStore{cats: map[string]*CategoryRef{}},
		&mockTagStore{tags: map[string]TagRef{}},
	)
	id, _ := svc.Create(context.Background(), CreateInput{
		Title:   "MD Test",
		Content: "# Heading\n\n**bold** and *italic*\n\n```go\nfmt.Println(\"hi\")\n```",
	})

	p, _ := svc.Get(context.Background(), id)
	if p.HTMLContent == "" {
		t.Error("expected html_content")
	}
	if !contains(p.HTMLContent, "<h1") {
		t.Error("expected <h1> in rendered HTML")
	}
	if !contains(p.HTMLContent, "<strong>") {
		t.Error("expected <strong> in rendered HTML")
	}
	if !contains(p.HTMLContent, "<code") {
		t.Error("expected <code> in rendered HTML")
	}
}

func TestService_Summary_AutoExtract(t *testing.T) {
	svc := NewService(
		&mockPostStore{posts: map[string]*Post{}, tagIDs: map[string][]string{}},
		&mockCatStore{cats: map[string]*CategoryRef{}},
		&mockTagStore{tags: map[string]TagRef{}},
	)

	// no summary, auto-extract from content
	id, _ := svc.Create(context.Background(), CreateInput{
		Title:   "Auto Summary",
		Content: "This is some content that should be extracted as a summary.",
	})
	p, _ := svc.Get(context.Background(), id)
	if p.Summary == "" {
		t.Error("expected auto-extracted summary")
	}

	// explicit summary
	id2, _ := svc.Create(context.Background(), CreateInput{
		Title:   "Explicit",
		Content: "Some content",
		Summary: "My custom summary",
	})
	p2, _ := svc.Get(context.Background(), id2)
	if p2.Summary != "My custom summary" {
		t.Errorf("summary = %q, want 'My custom summary'", p2.Summary)
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
