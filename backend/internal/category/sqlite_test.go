package category

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupCatDB(t *testing.T) *SQLiteStore {
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

func TestCategorySQLite_Create(t *testing.T) {
	store := setupCatDB(t)
	err := store.Create(context.Background(), &Category{ID: "c1", Name: "Go", Slug: "go"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	// duplicate name
	err = store.Create(context.Background(), &Category{ID: "c2", Name: "Go", Slug: "go2"})
	if err == nil {
		t.Fatal("expected error for duplicate name")
	}
}

func TestCategorySQLite_Get(t *testing.T) {
	store := setupCatDB(t)
	store.Create(context.Background(), &Category{ID: "c1", Name: "Go", Slug: "go"})

	c, err := store.Get(context.Background(), "c1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if c.Name != "Go" {
		t.Errorf("name = %q, want 'Go'", c.Name)
	}

	_, err = store.Get(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent")
	}
}

func TestCategorySQLite_GetBySlug(t *testing.T) {
	store := setupCatDB(t)
	store.Create(context.Background(), &Category{ID: "c1", Name: "Go", Slug: "go"})

	c, err := store.GetBySlug(context.Background(), "go")
	if err != nil {
		t.Fatalf("getBySlug: %v", err)
	}
	if c.Name != "Go" {
		t.Errorf("name = %q, want 'Go'", c.Name)
	}

	_, err = store.GetBySlug(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent slug")
	}
}

func TestCategorySQLite_List(t *testing.T) {
	store := setupCatDB(t)
	store.Create(context.Background(), &Category{ID: "c1", Name: "B", Slug: "b"})
	store.Create(context.Background(), &Category{ID: "c2", Name: "A", Slug: "a"})

	list, err := store.List(context.Background())
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("len = %d, want 2", len(list))
	}
	// should be sorted by name
	if list[0].Name != "A" {
		t.Errorf("first = %q, want 'A'", list[0].Name)
	}
}

func TestCategorySQLite_Update(t *testing.T) {
	store := setupCatDB(t)
	store.Create(context.Background(), &Category{ID: "c1", Name: "Old", Slug: "old"})

	if err := store.Update(context.Background(), &Category{ID: "c1", Name: "New", Slug: "new"}); err != nil {
		t.Fatalf("update: %v", err)
	}

	c, _ := store.Get(context.Background(), "c1")
	if c.Name != "New" || c.Slug != "new" {
		t.Errorf("name=%s slug=%s, want New/new", c.Name, c.Slug)
	}
}

func TestCategorySQLite_Delete(t *testing.T) {
	store := setupCatDB(t)
	store.Create(context.Background(), &Category{ID: "c1", Name: "X", Slug: "x"})

	if err := store.Delete(context.Background(), "c1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	_, err := store.Get(context.Background(), "c1")
	if err == nil {
		t.Fatal("expected error after delete")
	}
}
