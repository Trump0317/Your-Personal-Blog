package post

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func setupPostDB(t *testing.T) *SQLiteStore {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	db.Exec("PRAGMA journal_mode=WAL")
	db.Exec("PRAGMA foreign_keys=ON")

	store := NewSQLiteStore(db)
	if err := store.Migrate(context.Background()); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return store
}

func createPost(t *testing.T, store *SQLiteStore, title, slug string, published bool) *Post {
	t.Helper()
	p := &Post{
		ID:          "id-" + slug,
		Title:       title,
		Slug:        slug,
		Content:     "# " + title,
		HTMLContent: "<h1>" + title + "</h1>",
		Summary:     "summary of " + title,
		CategoryID:  "",
		Published:   published,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := store.Create(context.Background(), p); err != nil {
		t.Fatalf("create post: %v", err)
	}
	return p
}

func TestSQLiteStore_Create(t *testing.T) {
	store := setupPostDB(t)

	err := store.Create(context.Background(), &Post{
		ID:        "p1",
		Title:     "Test",
		Slug:      "test",
		Content:   "# Test",
		Published: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	// duplicate slug
	err = store.Create(context.Background(), &Post{
		ID: "p2", Title: "Test 2", Slug: "test",
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	})
	if err == nil {
		t.Fatal("expected error for duplicate slug")
	}
}

func TestSQLiteStore_Get(t *testing.T) {
	store := setupPostDB(t)
	createPost(t, store, "Hello", "hello", true)

	// get by id
	p, err := store.Get(context.Background(), "id-hello")
	if err != nil {
		t.Fatalf("get by id: %v", err)
	}
	if p.Title != "Hello" {
		t.Errorf("title = %q, want %q", p.Title, "Hello")
	}
	if len(p.TagIDs) != 0 {
		t.Errorf("tag_ids = %v, want empty", p.TagIDs)
	}

	// get by slug
	p, err = store.Get(context.Background(), "hello")
	if err != nil {
		t.Fatalf("get by slug: %v", err)
	}
	if p.Title != "Hello" {
		t.Errorf("title = %q, want %q", p.Title, "Hello")
	}

	// not found
	_, err = store.Get(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent")
	}
}

func TestSQLiteStore_List(t *testing.T) {
	store := setupPostDB(t)
	createPost(t, store, "Pub1", "pub1", true)
	createPost(t, store, "Pub2", "pub2", true)
	createPost(t, store, "Draft1", "draft1", false)

	// list all
	posts, total, err := store.List(context.Background(), Filter{Limit: 10})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total != 3 {
		t.Errorf("total = %d, want 3", total)
	}

	// filter published only
	pub := true
	posts, total, err = store.List(context.Background(), Filter{Published: &pub, Limit: 10})
	if err != nil {
		t.Fatalf("list published: %v", err)
	}
	if total != 2 {
		t.Errorf("published total = %d, want 2", total)
	}
	for _, p := range posts {
		if !p.Published {
			t.Errorf("got unpublished post: %s", p.Title)
		}
	}

	// filter by query
	posts, total, err = store.List(context.Background(), Filter{Query: "Pub", Published: &pub, Limit: 10})
	if err != nil {
		t.Fatalf("list query: %v", err)
	}
	if total != 2 {
		t.Errorf("query total = %d, want 2", total)
	}

	// pagination
	posts, total, err = store.List(context.Background(), Filter{Limit: 1, Offset: 1})
	if err != nil {
		t.Fatalf("list paginated: %v", err)
	}
	if total != 3 {
		t.Errorf("total = %d, want 3", total)
	}
	if len(posts) != 1 {
		t.Errorf("len = %d, want 1", len(posts))
	}
}

func TestSQLiteStore_Update(t *testing.T) {
	store := setupPostDB(t)
	createPost(t, store, "Old", "old", true)

	p, _ := store.Get(context.Background(), "id-old")
	p.Title = "New Title"
	p.Content = "updated"
	p.HTMLContent = "<p>updated</p>"
	p.Published = false
	p.UpdatedAt = time.Now()

	if err := store.Update(context.Background(), p); err != nil {
		t.Fatalf("update: %v", err)
	}

	got, _ := store.Get(context.Background(), "id-old")
	if got.Title != "New Title" {
		t.Errorf("title = %q, want %q", got.Title, "New Title")
	}
	if got.Published {
		t.Error("expected published=false")
	}
}

func TestSQLiteStore_Delete(t *testing.T) {
	store := setupPostDB(t)
	createPost(t, store, "Gone", "gone", true)

	if err := store.Delete(context.Background(), "id-gone"); err != nil {
		t.Fatalf("delete: %v", err)
	}

	_, err := store.Get(context.Background(), "id-gone")
	if err == nil {
		t.Fatal("expected error after delete")
	}
}

func TestSQLiteStore_SetTags(t *testing.T) {
	store := setupPostDB(t)
	createPost(t, store, "Tagged", "tagged", true)

	if err := store.SetTags(context.Background(), "id-tagged", []string{"t1", "t2"}); err != nil {
		t.Fatalf("set tags: %v", err)
	}

	ids, err := store.getTagIDs(context.Background(), "id-tagged")
	if err != nil {
		t.Fatalf("get tag ids: %v", err)
	}
	if len(ids) != 2 {
		t.Errorf("tag count = %d, want 2", len(ids))
	}

	// verify Get also loads tags
	p, _ := store.Get(context.Background(), "id-tagged")
	if len(p.TagIDs) != 2 {
		t.Errorf("post tag_ids = %d, want 2", len(p.TagIDs))
	}

	// update tags (replace)
	if err := store.SetTags(context.Background(), "id-tagged", []string{"t3"}); err != nil {
		t.Fatalf("set tags update: %v", err)
	}
	p, _ = store.Get(context.Background(), "id-tagged")
	if len(p.TagIDs) != 1 || p.TagIDs[0] != "t3" {
		t.Errorf("updated tag_ids = %v, want [t3]", p.TagIDs)
	}
}

func TestSQLiteStore_Stats(t *testing.T) {
	store := setupPostDB(t)
	createPost(t, store, "A", "a", true)
	createPost(t, store, "B", "b", true)
	createPost(t, store, "C", "c", false)

	stats, err := store.Stats(context.Background())
	if err != nil {
		t.Fatalf("stats: %v", err)
	}
	if stats.PostCount != 2 {
		t.Errorf("post_count = %d, want 2 (only published)", stats.PostCount)
	}
}

func TestSQLiteStore_Archive(t *testing.T) {
	store := setupPostDB(t)

	// 用原始 SQL 插入以控制 created_at (Update 不改 created_at)
	db := store.db
	db.ExecContext(context.Background(),
		`INSERT INTO posts (id, title, slug, content, html_content, summary, published, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"old-id", "Old Post", "old-post", "content", "<p>content</p>", "summary", 1,
		"2024-03-15T00:00:00Z", "2024-03-15T00:00:00Z")

	archive, err := store.Archive(context.Background())
	if err != nil {
		t.Fatalf("archive: %v", err)
	}
	if len(archive) == 0 {
		t.Fatal("expected archive entries")
	}
	if archive[0].Year != 2024 {
		t.Errorf("year = %d, want 2024", archive[0].Year)
	}
}

func TestSQLiteStore_FilterByTag(t *testing.T) {
	store := setupPostDB(t)
	createPost(t, store, "Go Post", "go-post", true)
	createPost(t, store, "JS Post", "js-post", true)

	store.SetTags(context.Background(), "id-go-post", []string{"tag-go"})
	store.SetTags(context.Background(), "id-js-post", []string{"tag-js"})

	pub := true
	posts, total, _ := store.List(context.Background(), Filter{TagID: "tag-go", Published: &pub, Limit: 10})
	if total != 1 {
		t.Errorf("tag filter total = %d, want 1", total)
	}
	if len(posts) > 0 && posts[0].Title != "Go Post" {
		t.Errorf("tag filtered post = %s, want 'Go Post'", posts[0].Title)
	}
}

func TestSQLiteStore_FilterByCategory(t *testing.T) {
	store := setupPostDB(t)
	p1 := createPost(t, store, "Cat A", "cat-a", true)
	p2 := createPost(t, store, "Cat B", "cat-b", true)
	p1.CategoryID = "cat-1"
	p2.CategoryID = "cat-2"
	store.Update(context.Background(), p1)
	store.Update(context.Background(), p2)

	pub := true
	_, total, _ := store.List(context.Background(), Filter{CategoryID: "cat-1", Published: &pub, Limit: 10})
	if total != 1 {
		t.Errorf("category filter total = %d, want 1", total)
	}
}
