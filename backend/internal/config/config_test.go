package config

import (
	"os"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()
	if cfg.Server.Port != ":8080" {
		t.Errorf("port = %q, want :8080", cfg.Server.Port)
	}
	if cfg.Database.Path != "blog.db" {
		t.Errorf("db path = %q, want blog.db", cfg.Database.Path)
	}
	if cfg.Admin.Username != "admin" {
		t.Errorf("admin user = %q, want admin", cfg.Admin.Username)
	}
	if cfg.Admin.Password != "password" {
		t.Errorf("admin pass = %q, want password", cfg.Admin.Password)
	}
}

func TestEnvOverrides(t *testing.T) {
	os.Setenv("BLOG_PORT", ":3000")
	os.Setenv("BLOG_DB", "test.db")
	os.Setenv("BLOG_ADMIN_USER", "root")
	os.Setenv("BLOG_ADMIN_PASS", "secret")
	defer func() {
		os.Unsetenv("BLOG_PORT")
		os.Unsetenv("BLOG_DB")
		os.Unsetenv("BLOG_ADMIN_USER")
		os.Unsetenv("BLOG_ADMIN_PASS")
	}()

	cfg := Default()
	if cfg.Server.Port != ":3000" {
		t.Errorf("port = %q, want :3000", cfg.Server.Port)
	}
	if cfg.Database.Path != "test.db" {
		t.Errorf("db path = %q, want test.db", cfg.Database.Path)
	}
	if cfg.Admin.Username != "root" {
		t.Errorf("admin user = %q, want root", cfg.Admin.Username)
	}
	if cfg.Admin.Password != "secret" {
		t.Errorf("admin pass = %q, want secret", cfg.Admin.Password)
	}
}

func TestEnv_Fallback(t *testing.T) {
	// env returns fallback when not set
	v := env("NONEXISTENT_VAR_12345", "fallback")
	if v != "fallback" {
		t.Errorf("env = %q, want fallback", v)
	}
}

func TestEnv_ReturnsValue(t *testing.T) {
	os.Setenv("TEST_EXISTING", "hello")
	defer os.Unsetenv("TEST_EXISTING")
	v := env("TEST_EXISTING", "fallback")
	if v != "hello" {
		t.Errorf("env = %q, want hello", v)
	}
}
