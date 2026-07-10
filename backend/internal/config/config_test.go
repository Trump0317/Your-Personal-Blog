package config

import (
	"os"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	os.Setenv("BLOG_CONFIG", "/nonexistent/path")
	defer os.Unsetenv("BLOG_CONFIG")

	cfg := Load()
	if cfg.Server.Port != ":8080" {
		t.Errorf("port = %q, want :8080", cfg.Server.Port)
	}
	if cfg.Admin.Username != "admin" {
		t.Errorf("admin user = %q, want admin", cfg.Admin.Username)
	}
	if cfg.Admin.PasswordHash == "" {
		t.Error("expected password hash to be set")
	}
	if cfg.Admin.JWTSecret == "" {
		t.Error("expected JWT secret to be auto-generated")
	}
}

func TestLoad_FileWithEnvOverride(t *testing.T) {
	content := `{"site":{"name":"FromFile"},"admin":{"username":"root","password":"secret"}}`
	tmpFile := t.TempDir() + "/test-config.json"
	os.WriteFile(tmpFile, []byte(content), 0644)
	os.Setenv("BLOG_CONFIG", tmpFile)
	os.Setenv("BLOG_SITE_NAME", "FromEnv")
	defer func() {
		os.Unsetenv("BLOG_CONFIG")
		os.Unsetenv("BLOG_SITE_NAME")
	}()

	cfg := Load()
	if cfg.Site.Name != "FromEnv" {
		t.Errorf("name = %q, want 'FromEnv' (env overrides file)", cfg.Site.Name)
	}
	if cfg.Admin.Username != "root" {
		t.Errorf("user = %q, want root", cfg.Admin.Username)
	}
	if cfg.Admin.PasswordHash == "" {
		t.Error("expected password hash from file config")
	}
}

func TestVerifyPassword(t *testing.T) {
	os.Setenv("BLOG_CONFIG", "/nonexistent/path")
	os.Setenv("BLOG_ADMIN_USER", "admin")
	os.Setenv("BLOG_ADMIN_PASS", "test123")
	defer func() {
		os.Unsetenv("BLOG_CONFIG")
		os.Unsetenv("BLOG_ADMIN_USER")
		os.Unsetenv("BLOG_ADMIN_PASS")
	}()

	cfg := Load()
	if !cfg.VerifyPassword("test123") {
		t.Error("expected password 'test123' to verify")
	}
	if cfg.VerifyPassword("wrong") {
		t.Error("expected 'wrong' password to fail")
	}
	if cfg.VerifyPassword("") {
		t.Error("expected empty password to fail")
	}
}

func TestLoad_JWTSecret(t *testing.T) {
	os.Setenv("BLOG_CONFIG", "/nonexistent/path")
	defer os.Unsetenv("BLOG_CONFIG")

	cfg1 := Load()
	cfg2 := Load()
	// 每次启动生成随机 secret
	if cfg1.Admin.JWTSecret == cfg2.Admin.JWTSecret {
		t.Error("expected different JWT secrets per Load() call")
	}
	if len(cfg1.Admin.JWTSecret) != 32 {
		t.Errorf("JWT secret length = %d, want 32", len(cfg1.Admin.JWTSecret))
	}
}
