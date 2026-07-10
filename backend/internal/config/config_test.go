package config

import (
	"os"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	// 确保没有 config.json 干扰
	os.Setenv("BLOG_CONFIG", "/nonexistent/path")
	defer os.Unsetenv("BLOG_CONFIG")

	cfg := Load()
	if cfg.Server.Port != ":8080" {
		t.Errorf("port = %q, want :8080", cfg.Server.Port)
	}
	if cfg.Database.Path != "blog.db" {
		t.Errorf("db path = %q, want blog.db", cfg.Database.Path)
	}
	if cfg.Admin.Username != "admin" {
		t.Errorf("admin user = %q, want admin", cfg.Admin.Username)
	}
	if cfg.Site.Name != "Your Personal Blog" {
		t.Errorf("site name = %q", cfg.Site.Name)
	}
}

func TestLoad_EnvOverrides(t *testing.T) {
	os.Setenv("BLOG_CONFIG", "/nonexistent/path")
	os.Setenv("BLOG_PORT", ":3000")
	os.Setenv("BLOG_DB", "test.db")
	os.Setenv("BLOG_ADMIN_USER", "root")
	os.Setenv("BLOG_SITE_NAME", "My Site")
	defer func() {
		for _, k := range []string{"BLOG_CONFIG", "BLOG_PORT", "BLOG_DB", "BLOG_ADMIN_USER", "BLOG_SITE_NAME"} {
			os.Unsetenv(k)
		}
	}()

	cfg := Load()
	if cfg.Server.Port != ":3000" {
		t.Errorf("port = %q, want :3000", cfg.Server.Port)
	}
	if cfg.Database.Path != "test.db" {
		t.Errorf("db = %q, want test.db", cfg.Database.Path)
	}
	if cfg.Admin.Username != "root" {
		t.Errorf("user = %q, want root", cfg.Admin.Username)
	}
	if cfg.Site.Name != "My Site" {
		t.Errorf("site name = %q, want 'My Site'", cfg.Site.Name)
	}
}

func TestLoad_FromFile(t *testing.T) {
	// 写临时配置文件
	content := `{"server":{"port":":9999"},"site":{"name":"FileBlog","tagline":"from file"}}`
	tmpFile := t.TempDir() + "/test-config.json"
	os.WriteFile(tmpFile, []byte(content), 0644)
	os.Setenv("BLOG_CONFIG", tmpFile)
	defer os.Unsetenv("BLOG_CONFIG")

	cfg := Load()
	if cfg.Server.Port != ":9999" {
		t.Errorf("port = %q, want :9999", cfg.Server.Port)
	}
	if cfg.Site.Name != "FileBlog" {
		t.Errorf("name = %q, want 'FileBlog'", cfg.Site.Name)
	}
	if cfg.Site.Tagline != "from file" {
		t.Errorf("tagline = %q, want 'from file'", cfg.Site.Tagline)
	}
	// 未指定字段应使用默认值
	if cfg.Database.Path != "blog.db" {
		t.Errorf("db = %q, want blog.db (default)", cfg.Database.Path)
	}
}

func TestLoad_FileWithEnvOverride(t *testing.T) {
	content := `{"site":{"name":"FromFile"}}`
	tmpFile := t.TempDir() + "/test-config2.json"
	os.WriteFile(tmpFile, []byte(content), 0644)
	os.Setenv("BLOG_CONFIG", tmpFile)
	os.Setenv("BLOG_SITE_NAME", "FromEnv")
	defer func() {
		os.Unsetenv("BLOG_CONFIG")
		os.Unsetenv("BLOG_SITE_NAME")
	}()

	cfg := Load()
	// 环境变量优先级高于文件
	if cfg.Site.Name != "FromEnv" {
		t.Errorf("name = %q, want 'FromEnv' (env overrides file)", cfg.Site.Name)
	}
}
