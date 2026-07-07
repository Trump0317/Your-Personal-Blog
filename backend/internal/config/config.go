package config

import (
	"os"
)

type Config struct {
	Server   Server
	Database Database
	Admin    Admin
}

type Server struct {
	Port string
}

type Database struct {
	Path string
}

type Admin struct {
	Username string
	Password string
}

// Default 返回默认配置，可通过环境变量覆盖
func Default() *Config {
	return &Config{
		Server: Server{
			Port: env("BLOG_PORT", ":8080"),
		},
		Database: Database{
			Path: env("BLOG_DB", "blog.db"),
		},
		Admin: Admin{
			Username: env("BLOG_ADMIN_USER", "admin"),
			Password: env("BLOG_ADMIN_PASS", "password"),
		},
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
