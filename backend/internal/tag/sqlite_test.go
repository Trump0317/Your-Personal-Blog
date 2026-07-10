package tag

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTagDB(t *testing.T) *SQLiteStore {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	db.Exec("PRAGMA foreign_keys=ON")

	store := NewSQLiteStore(db)
	if err := store.Migrate(context.Background()); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return store
}

func TestTagSQLite_Create(t *testing.T) {
	store := setupTagDB(t)
	err := store.Create(context.Background(), &Tag{ID: "t1", Name: "go", Slug: "go"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	// duplicate name
	err = store.Create(context.Background(), &Tag{ID: "t2", Name: "go", Slug: "go2"})
	if err == nil {
		t.Fatal("expected error for duplicate name")
	}
}

func TestTagSQLite_Get(t *testing.T) {
	store := setupTagDB(t)
	store.Create(context.Background(), &Tag{ID: "t1", Name: "go", Slug: "go"})

	tag, err := store.Get(context.Background(), "t1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if tag.Name != "go" {
		t.Errorf("name = %q, want 'go'", tag.Name)
	}
}

func TestTagSQLite_GetByName(t *testing.T) {
	store := setupTagDB(t)
	store.Create(context.Background(), &Tag{ID: "t1", Name: "go", Slug: "go"})

	tag, err := store.GetByName(context.Background(), "go")
	if err != nil {
		t.Fatalf("getByName: %v", err)
	}
	if tag.ID != "t1" {
		t.Errorf("id = %q, want 't1'", tag.ID)
	}
}

func TestTagSQLite_GetBySlug(t *testing.T) {
	store := setupTagDB(t)
	store.Create(context.Background(), &Tag{ID: "t1", Name: "go", Slug: "go"})

	tag, err := store.GetBySlug(context.Background(), "go")
	if err != nil {
		t.Fatalf("getBySlug: %v", err)
	}
	if tag.Name != "go" {
		t.Errorf("name = %q, want 'go'", tag.Name)
	}
}

func TestTagSQLite_List(t *testing.T) {
	store := setupTagDB(t)
	store.Create(context.Background(), &Tag{ID: "t1", Name: "B", Slug: "b"})
	store.Create(context.Background(), &Tag{ID: "t2", Name: "A", Slug: "a"})

	list, err := store.List(context.Background())
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("len = %d, want 2", len(list))
	}
	if list[0].Name != "A" {
		t.Errorf("first = %q, want 'A' (sorted)", list[0].Name)
	}
}

func TestTagSQLite_GetByIDs(t *testing.T) {
	store := setupTagDB(t)
	store.Create(context.Background(), &Tag{ID: "t1", Name: "a", Slug: "a"})
	store.Create(context.Background(), &Tag{ID: "t2", Name: "b", Slug: "b"})
	store.Create(context.Background(), &Tag{ID: "t3", Name: "c", Slug: "c"})

	list, err := store.GetByIDs(context.Background(), []string{"t1", "t3"})
	if err != nil {
		t.Fatalf("getByIDs: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("len = %d, want 2", len(list))
	}

	// empty ids
	list, err = store.GetByIDs(context.Background(), []string{})
	if err != nil {
		t.Fatalf("getByIDs empty: %v", err)
	}
	if list != nil {
		t.Error("expected nil for empty ids")
	}
}

func TestTagSQLite_Update(t *testing.T) {
	store := setupTagDB(t)
	store.Create(context.Background(), &Tag{ID: "t1", Name: "Old", Slug: "old"})

	err := store.Update(context.Background(), &Tag{ID: "t1", Name: "New", Slug: "new"})
	if err != nil {
		t.Fatalf("update: %v", err)
	}

	tag, _ := store.Get(context.Background(), "t1")
	if tag.Name != "New" || tag.Slug != "new" {
		t.Errorf("name=%s slug=%s, want New/new", tag.Name, tag.Slug)
	}
}

func TestTagSQLite_Delete(t *testing.T) {
	store := setupTagDB(t)
	store.Create(context.Background(), &Tag{ID: "t1", Name: "x", Slug: "x"})

	if err := store.Delete(context.Background(), "t1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	_, err := store.Get(context.Background(), "t1")
	if err == nil {
		t.Fatal("expected error after delete")
	}
}
