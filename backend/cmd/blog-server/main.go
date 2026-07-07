package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"

	"github.com/ypb/your-personal-blog/internal/category"
	"github.com/ypb/your-personal-blog/internal/config"
	"github.com/ypb/your-personal-blog/internal/middleware"
	"github.com/ypb/your-personal-blog/internal/post"
	"github.com/ypb/your-personal-blog/internal/server"
	"github.com/ypb/your-personal-blog/internal/tag"
)

func main() {
	cfg := config.Default()

	// 1. 数据库
	db, err := sql.Open("sqlite3", cfg.Database.Path)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	// WAL 模式提升并发性能
	db.Exec("PRAGMA journal_mode=WAL")
	db.Exec("PRAGMA foreign_keys=ON")

	// 2. 迁移
	catStore := category.NewSQLiteStore(db)
	tagStore := tag.NewSQLiteStore(db)
	postStore := post.NewSQLiteStore(db)

	must(catStore.Migrate(context.Background()))
	must(tagStore.Migrate(context.Background()))
	must(postStore.Migrate(context.Background()))

	// 3. 服务层
	catSvc := category.NewService(catStore)
	tagSvc := tag.NewService(tagStore)

	// 适配器：将 category.Service / tag.Service 转为 post 包期望的接口
	catAdapter := &categoryAdapter{svc: catSvc}
	tagAdapter := &tagAdapter{svc: tagSvc}

	postSvc := post.NewService(postStore, catAdapter, tagAdapter)

	// 4. 种子数据
	seed(context.Background(), postSvc, catSvc, tagSvc)

	// 5. 路由
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())

	// 生产模式：serve Vue build 产物；开发模式：Vite dev server 代理前端
	distDir := resolveDistDir()
	r.Static("/assets", filepath.Join(distDir, "assets"))
	r.GET("/favicon.ico", func(c *gin.Context) { c.File(filepath.Join(distDir, "favicon.ico")) })
	r.GET("/", func(c *gin.Context) { c.File(filepath.Join(distDir, "index.html")) })
	// SPA fallback：非 /api 路径全部返回 index.html
	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.JSON(404, gin.H{"error": "not found"})
			return
		}
		c.File(filepath.Join(distDir, "index.html"))
	})

	api := r.Group("/api")

	// 公开路由
	postH := post.NewHandler(postSvc)
	catH := category.NewHandler(catSvc)
	tagH := tag.NewHandler(tagSvc)

	postH.RegisterPublic(api)
	catH.RegisterPublic(api)
	tagH.RegisterPublic(api)

	// 管理端路由（Basic Auth 保护）
	admin := api.Group("/admin")
	admin.Use(middleware.BasicAuth(cfg.Admin.Username, cfg.Admin.Password))

	postH.RegisterAdmin(admin)
	catH.RegisterAdmin(admin)
	tagH.RegisterAdmin(admin)

	// 6. 启动
	server.Run(r, cfg.Server.Port)
}

// resolveDistDir 查找前端构建产物目录
func resolveDistDir() string {
	candidates := []string{
		"frontend/dist",       // 从项目根目录运行
		"../frontend/dist",    // 从 backend/ 目录运行
	}
	for _, d := range candidates {
		if fi, err := os.Stat(d); err == nil && fi.IsDir() {
			return d
		}
	}
	return "frontend/dist" // fallback
}

// ── 适配器 ──

type categoryAdapter struct {
	svc *category.Service
}

func (a *categoryAdapter) Get(ctx context.Context, id string) (*post.CategoryRef, error) {
	c, err := a.svc.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &post.CategoryRef{ID: c.ID, Name: c.Name, Slug: c.Slug}, nil
}

func (a *categoryAdapter) List(ctx context.Context) ([]*post.CategoryRef, error) {
	cats, err := a.svc.List(ctx)
	if err != nil {
		return nil, err
	}
	refs := make([]*post.CategoryRef, len(cats))
	for i, c := range cats {
		refs[i] = &post.CategoryRef{ID: c.ID, Name: c.Name, Slug: c.Slug}
	}
	return refs, nil
}

type tagAdapter struct {
	svc *tag.Service
}

func (a *tagAdapter) GetByIDs(ctx context.Context, ids []string) ([]post.TagRef, error) {
	tags, err := a.svc.GetByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	refs := make([]post.TagRef, len(tags))
	for i, t := range tags {
		refs[i] = post.TagRef{ID: t.ID, Name: t.Name, Slug: t.Slug}
	}
	return refs, nil
}

func (a *tagAdapter) List(ctx context.Context) ([]post.TagRef, error) {
	tags, err := a.svc.List(ctx)
	if err != nil {
		return nil, err
	}
	refs := make([]post.TagRef, len(tags))
	for i, t := range tags {
		refs[i] = post.TagRef{ID: t.ID, Name: t.Name, Slug: t.Slug}
	}
	return refs, nil
}

// ── 辅助 ──

func must(err error) {
	if err != nil {
		log.Fatalf("migration error: %v", err)
	}
}

func seed(ctx context.Context, postSvc *post.Service, catSvc *category.Service, tagSvc *tag.Service) {
	// 检查是否已有数据
	pub := true
	posts, total, _ := postSvc.List(ctx, post.Filter{Published: &pub, Limit: 1})
	if total > 0 || len(posts) > 0 {
		return
	}

	// 创建分类
	cat1ID, _ := catSvc.Create(ctx, "技术分享")
	cat2ID, _ := catSvc.Create(ctx, "随笔感悟")

	// 创建标签
	tagNames := []string{"Go", "后端", "前端", "设计", "生活"}
	tagIDs := make([]string, len(tagNames))
	for i, name := range tagNames {
		id, _ := tagSvc.CreateOrGet(ctx, name)
		tagIDs[i] = id
	}

	// 创建文章
	postSvc.Create(ctx, post.CreateInput{
		Title:      "你好，世界 - 用 Go 写博客",
		Slug:       "hello-world",
		Content:    "## 欢迎\n\n这是我的第一篇博客，使用 **Go** 语言和 **SQLite** 构建。\n\n### 技术栈\n- Go + Gin\n- SQLite\n- Markdown 渲染\n\n```go\nfmt.Println(\"Hello, Blog!\")\n```",
		Summary:    "用 Go 语言构建的极简博客系统",
		CategoryID: cat1ID,
		TagIDs:     []string{tagIDs[0], tagIDs[1]},
		Published:  true,
	})

	postSvc.Create(ctx, post.CreateInput{
		Title:      "周末午后的咖啡馆",
		Slug:       "coffee-afternoon",
		Content:    "阳光洒进窗户，桌上的拿铁冒着热气。\n\n有时候慢下来，也是一种进步。",
		Summary:    "找一个静谧的角落，写写代码，看看窗外。",
		CategoryID: cat2ID,
		TagIDs:     []string{tagIDs[4]},
		Published:  true,
	})

	log.Println("seed data created")
}
