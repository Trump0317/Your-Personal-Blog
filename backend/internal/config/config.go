package config

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"log"
	"os"

	"golang.org/x/crypto/argon2"
)

type Config struct {
	Server   Server   `json:"server"`
	Database Database `json:"database"`
	Admin    Admin    `json:"admin"`
	Site     Site     `json:"site"`
}

type Server struct {
	Port string `json:"port"`
}

type Database struct {
	Path string `json:"path"`
}

type Admin struct {
	Username     string `json:"username"`
	Password     string `json:"password"` // 从配置文件读取明文，启动后清除
	PasswordHash string `json:"-"`         // Argon2 哈希（内存中）
	PasswordSalt string `json:"-"`         // 盐值（内存中）
	JWTSecret    string `json:"jwt_secret"`
}

type Site struct {
	Name        string `json:"name"`
	Tagline     string `json:"tagline"`
	Description string `json:"description"`
	FooterNote  string `json:"footer_note"`
	HeroKicker  string `json:"hero_kicker"`
	HeroTitle1  string `json:"hero_title1"`
	HeroTitle2  string `json:"hero_title2"`
	HeroTitle3  string `json:"hero_title3"`
	HeroDesc    string `json:"hero_desc"`
	AboutTitle  string `json:"about_title"`
	AboutText   string `json:"about_text"`
}

// Load 从 config.json 加载配置，环境变量可覆盖
func Load() *Config {
	cfg := defaults()

	path := os.Getenv("BLOG_CONFIG")
	if path == "" {
		path = "config.json"
	}

	if data, err := os.ReadFile(path); err == nil {
		if err := json.Unmarshal(data, cfg); err != nil {
			log.Printf("config: parse error in %s: %v, using defaults", path, err)
		}
	} else {
		log.Printf("config: %s not found, using defaults", path)
	}

	// 环境变量覆盖
	applyEnv(cfg)

	// 哈希密码
	cfg.hashPassword()

	// 生成 JWT secret（如果未配置）
	if cfg.Admin.JWTSecret == "" {
		cfg.Admin.JWTSecret = randomString(32)
	}

	return cfg
}

func (cfg *Config) hashPassword() {
	if cfg.Admin.Password == "" {
		return
	}
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		log.Printf("config: failed to generate salt, password disabled")
		cfg.Admin.Password = ""
		return
	}
	cfg.Admin.PasswordSalt = base64.StdEncoding.EncodeToString(salt)
	cfg.Admin.PasswordHash = base64.StdEncoding.EncodeToString(
		argon2.IDKey([]byte(cfg.Admin.Password), salt, 1, 64*1024, 4, 32))
	cfg.Admin.Password = "" // 清除明文
}

// VerifyPassword 验证密码
func (cfg *Config) VerifyPassword(password string) bool {
	if cfg.Admin.PasswordSalt == "" || cfg.Admin.PasswordHash == "" {
		return false
	}
	salt, _ := base64.StdEncoding.DecodeString(cfg.Admin.PasswordSalt)
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	expected, _ := base64.StdEncoding.DecodeString(cfg.Admin.PasswordHash)
	return subtle.ConstantTimeCompare(hash, expected) == 1
}

func randomString(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:n]
}

func defaults() *Config {
	return &Config{
		Server:   Server{Port: ":8080"},
		Database: Database{Path: "blog.db"},
		Admin:    Admin{Username: "admin", Password: "password"},
		Site: Site{
			Name:       "Your Personal Blog",
			Tagline:    "探索技术、设计与思考的交汇",
			Description: "探索技术、设计与思考的交汇。以极简致敬繁复。",
			FooterNote: "探索技术、设计与思考的交汇",
			HeroKicker: "Daily Engineering Notes",
			HeroTitle1: "代码.",
			HeroTitle2: "匠心.",
			HeroTitle3: "文化.",
			HeroDesc:   "探索软件架构、工程美学与人类体验的交汇点。以工匠精神，雕琢每一行代码。",
			AboutTitle: "关于 Your Personal Blog",
			AboutText:  "这是一个用 Go + SQLite + Vue 构建的个人博客。静态分发优先，无运行时依赖。\n写博客的目的：记录思考，沉淀知识，分享实践。",
		},
	}
}

func applyEnv(cfg *Config) {
	if v := os.Getenv("BLOG_PORT"); v != "" {
		cfg.Server.Port = v
	}
	if v := os.Getenv("BLOG_DB"); v != "" {
		cfg.Database.Path = v
	}
	if v := os.Getenv("BLOG_ADMIN_USER"); v != "" {
		cfg.Admin.Username = v
	}
	if v := os.Getenv("BLOG_ADMIN_PASS"); v != "" {
		cfg.Admin.Password = v
	}
	if v := os.Getenv("BLOG_SITE_NAME"); v != "" {
		cfg.Site.Name = v
	}
	if v := os.Getenv("BLOG_SITE_TAGLINE"); v != "" {
		cfg.Site.Tagline = v
	}
	if v := os.Getenv("BLOG_SITE_DESC"); v != "" {
		cfg.Site.Description = v
	}
	if v := os.Getenv("BLOG_FOOTER_NOTE"); v != "" {
		cfg.Site.FooterNote = v
	}
	if v := os.Getenv("BLOG_HERO_KICKER"); v != "" {
		cfg.Site.HeroKicker = v
	}
	if v := os.Getenv("BLOG_HERO_TITLE1"); v != "" {
		cfg.Site.HeroTitle1 = v
	}
	if v := os.Getenv("BLOG_HERO_TITLE2"); v != "" {
		cfg.Site.HeroTitle2 = v
	}
	if v := os.Getenv("BLOG_HERO_TITLE3"); v != "" {
		cfg.Site.HeroTitle3 = v
	}
	if v := os.Getenv("BLOG_HERO_DESC"); v != "" {
		cfg.Site.HeroDesc = v
	}
	if v := os.Getenv("BLOG_ABOUT_TITLE"); v != "" {
		cfg.Site.AboutTitle = v
	}
	if v := os.Getenv("BLOG_ABOUT_TEXT"); v != "" {
		cfg.Site.AboutText = v
	}
}
