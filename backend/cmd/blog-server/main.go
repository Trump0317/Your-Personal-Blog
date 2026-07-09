package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

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

	// 上传
	r.Static("/uploads", "./uploads")
	r.POST("/api/upload", middleware.BasicAuth(cfg.Admin.Username, cfg.Admin.Password), uploadHandler)

	// RSS
	r.GET("/rss.xml", func(c *gin.Context) { rssHandler(c, postSvc) })

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

// ── RSS ──
func rssHandler(c *gin.Context, postSvc *post.Service) {
	pub := true
	posts, _, err := postSvc.List(c.Request.Context(), post.Filter{Published: &pub, Limit: 20})
	if err != nil {
		c.String(500, "error")
		return
	}

	// 从请求 host 或环境变量推导站点 URL
	baseURL := resolveBaseURL(c)

	var items strings.Builder
	for _, p := range posts {
		items.WriteString(fmt.Sprintf(`
<item>
  <title>%s</title>
  <link>%s/post/%s</link>
  <description>%s</description>
  <pubDate>%s</pubDate>
</item>`,
			escXML(p.Title), baseURL, escXML(p.Slug), escXML(p.Summary),
			p.CreatedAt.Format(time.RFC1123Z)))
	}

	rss := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
<channel>
  <title>My Blog</title>
  <link>%s</link>
  <description>探索技术、设计与思考的交汇</description>%s
</channel>
</rss>`, baseURL, items.String())

	c.Header("Content-Type", "application/rss+xml")
	c.String(200, rss)
}

func resolveBaseURL(c *gin.Context) string {
	if u := os.Getenv("BLOG_BASE_URL"); u != "" {
		return strings.TrimRight(u, "/")
	}
	scheme := "http"
	if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	return scheme + "://" + c.Request.Host
}

func escXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

// ── 图片上传 ──
func uploadHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "no file"})
		return
	}

	os.MkdirAll("uploads", 0755)
	filename := fmt.Sprintf("uploads/%d_%s", time.Now().UnixMilli(), file.Filename)
	if err := c.SaveUploadedFile(file, filename); err != nil {
		c.JSON(500, gin.H{"error": "upload failed"})
		return
	}
	c.JSON(200, gin.H{"url": "/" + filename})
}

func seed(ctx context.Context, postSvc *post.Service, catSvc *category.Service, tagSvc *tag.Service) {
	pub := true
	posts, total, _ := postSvc.List(ctx, post.Filter{Published: &pub, Limit: 1})
	if total > 0 || len(posts) > 0 {
		return
	}

	// 分类
	catGo, _ := catSvc.Create(ctx, "Go 语言")
	catFE, _ := catSvc.Create(ctx, "前端开发")
	catArch, _ := catSvc.Create(ctx, "系统架构")
	catLife, _ := catSvc.Create(ctx, "生活随笔")

	// 标签
	t := func(name string) string { id, _ := tagSvc.CreateOrGet(ctx, name); return id }
	tGo, tGin, tSQLite := t("Go"), t("Gin"), t("SQLite")
	tVue, tVite, tCSS := t("Vue"), t("Vite"), t("CSS")
	tMicroservice, tDocker, tCI := t("微服务"), t("Docker"), t("CI/CD")
	tThinking, tReading := t("思考"), t("阅读")

	// 10 篇文章
	articles := []struct {
		title, slug, summary, content string
		cat                           string
		tags                          []string
	}{
		{"你好世界：用 Go 写一个博客系统", "hello-world-go-blog",
			"从零开始，用 Go + Gin + SQLite 构建一个模块化的个人博客。",
			`## 为什么用 Go 写博客

Go 语言的高并发和简洁语法，让它成为构建 Web 服务的最佳选择之一。本项目采用模块化架构，每个领域（文章、分类、标签）自包含模型、存储和 HTTP 处理。

### 技术选型

- **Gin** — 高性能 HTTP 框架
- **SQLite** — 零配置嵌入式数据库
- **Vue 3 + Vite** — 响应式前端
- **gomarkdown** — Markdown 渲染

### 核心代码

` + "```go" + `
func main() {
    r := gin.Default()
    r.GET("/api/posts", handleListPosts)
    r.Run(":8080")
}
` + "```" + `

> 好的架构不是过度设计，而是在正确的抽象层做正确的取舍。`, catGo, []string{tGo, tGin, tSQLite}},

		{"Gin 中间件完全指南：从 Logger 到 Auth", "gin-middleware-guide",
			"深入理解 Gin 中间件机制，手写日志、认证、恢复三大中间件。",
			`## 中间件是什么

Gin 的中间件本质上是一个 ` + "`gin.HandlerFunc`" + `，在请求到达最终 handler 之前和之后执行。

### Logger 中间件

记录每个请求的方法、路径、状态码和耗时：

` + "```go" + `
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Next()
        log.Printf("[%s] %s %d %v",
            c.Request.Method,
            c.Request.URL.Path,
            c.Writer.Status(),
            time.Since(start))
    }
}
` + "```" + `

### 中间件链

Gin 使用 ` + "`c.Next()`" + ` 实现洋葱模型：请求 → Logger → Auth → Handler → Logger，响应原路返回。

### 最佳实践

1. 全局中间件放在路由组外层
2. 认证中间件只保护需要鉴权的路由
3. Recovery 必须放在第一位`, catGo, []string{tGo, tGin}},

		{"Vue 3 Composition API 实战：构建博客前端", "vue3-composition-api",
			"用 Composition API + `<script setup>` 重写博客前端，代码量减少 40%。",
			`## 为什么选择 Composition API

Options API 适合小型组件，但当组件逻辑变复杂时，` + "`data`、`methods`、`computed`" + ` 分散在不同区域，难以维护。

### 对比

**Options API:**
` + "```js" + `
export default {
  data() { return { posts: [] } },
  methods: { async fetchPosts() { ... } },
  mounted() { this.fetchPosts() }
}
` + "```" + `

**Composition API:**
` + "```js" + `
const posts = ref([])
async function fetchPosts() { ... }
onMounted(() => fetchPosts())
` + "```" + `

### 实际收益

- 逻辑按功能组织，而不是按选项
- ` + "`<script setup>`" + ` 语法糖减少样板代码
- 更好的 TypeScript 支持`, catFE, []string{tVue, tVite}},

		{"Tailwind 之外的另一种选择：手写现代 CSS", "modern-css-without-tailwind",
			"CSS 自定义属性、Grid、clamp() — 原生 CSS 已经足够强大，不需要额外框架。",
			`## CSS 的现代化

2026 年的原生 CSS 已经非常强大，配合 CSS 自定义属性可以实现完整的主题系统。

### 主题变量

` + "```css" + `
:root {
  --bg: #fafaf9;
  --text: #0f172a;
  --accent: #f97316;
}

.dark {
  --bg: #0f172a;
  --text: #f1f5f9;
}
` + "```" + `

### 为什么不用 Tailwind

1. HTML 可读性：` + "`class=\"flex items-center gap-4 p-6 rounded-xl shadow-lg bg-white\"`" + ` vs ` + "`class=\"card\"`" + `
2. 样式和结构分离的哲学
3. 更小的 CSS 体积（本项目 ~12KB）

> 选择合适的工具，而不是流行的工具。`, catFE, []string{tCSS, tVue}},

		{"SQLite 在生产环境的正确打开方式", "sqlite-in-production",
			"SQLite 不是玩具。WAL 模式、连接池、备份策略 — 个人博客到中小型应用的最佳实践。",
			`## 打破偏见

很多人认为 SQLite 只适合开发环境。实际上，SQLite 在以下场景表现优异：

### 适用场景

- 个人博客、小型 CMS
- 嵌入式设备、移动应用
- 单机数据分析
- 边缘计算节点

### 关键配置

` + "```sql" + `
-- WAL 模式：读写不互斥
PRAGMA journal_mode=WAL;

-- 外键约束
PRAGMA foreign_keys=ON;

-- 缓存大小
PRAGMA cache_size = -8000; -- 8MB
` + "```" + `

### 备份策略

` + "```bash" + `
sqlite3 blog.db ".backup backup.db"
` + "```" + `

> 单机应用用 SQLite，多机应用才需要 PostgreSQL。`, catArch, []string{tSQLite, tGo}},

		{"微服务不是银弹：从单体到模块化的演进之路", "microservice-is-not-silver-bullet",
			"微服务很热门，但不适合所有场景。看看什么时候该拆，什么时候不该拆。",
			`## 先问为什么

微服务解决的是**组织扩展**问题，而非技术问题。如果你的团队只有 5 个人，微服务只会增加不必要的复杂度。

### 模块化单体的优势

1. **部署简单** — 一个二进制，一次部署
2. **调试方便** — 一个 IDE 搞定所有
3. **事务简单** — 数据库事务天然 ACID
4. **性能更好** — 进程内调用，零网络开销

### 什么时候该拆

- 团队超过 20 人，不同小组负责不同模块
- 独立部署节奏（A 模块每周发版，B 模块每月）
- 独立扩缩容需求
- 技术栈异构（Go + Python + Rust）

> 架构不是一次性设计出来的，是随着需求演进的。`, catArch, []string{tMicroservice, tDocker}},

		{"用 Docker 部署 Go 应用：从 Dockerfile 到 Compose", "dockerize-go-app",
			"多阶段构建减小镜像体积，Docker Compose 编排前后端服务。",
			`## 多阶段构建

Go 是编译型语言，只需要运行时二进制。多阶段构建可以把编译依赖从最终镜像中剥离：

` + "```dockerfile" + `
# 构建阶段
FROM golang:1.26 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o server ./cmd/blog-server/

# 运行阶段
FROM alpine:3.20
COPY --from=builder /app/server /server
COPY frontend/dist /frontend/dist
EXPOSE 8080
CMD ["/server"]
` + "```" + `

最终镜像只有 ~15MB，而不是 ~800MB。`, catArch, []string{tGo, tDocker, tCI}},

		{"程序员的写作之道", "writing-for-engineers",
			"代码是给机器看的，文档是给人看的。写好技术文章的几个简单原则。",
			`## 为什么要写

写作是程序员的第二门语言。好的技术文章可以：

- 沉淀知识，避免遗忘
- 帮助他人，建立影响力
- 倒逼深入理解（费曼学习法）

### 写作原则

1. **先说结论** — 读者没耐心从头读到尾
2. **用代码说话** — 一个可运行的示例胜过三段描述
3. **短段落** — 每段不超过 4 句话
4. **删除冗余** — 写完读一遍，删掉 30%

### 推荐的写作流程

> 草稿 → 放一天 → 删改 → 配图 → 发布

好文章是改出来的，不是写出来的。`, catLife, []string{tThinking}},

		{"我的书架：2026 年上半年阅读清单", "reading-list-2026-h1",
			"五本改变我思维方式的书，涵盖技术、设计和哲学。",
			`## 技术之外

程序员不能只读技术书。以下五本书让我在 2026 年上半年获得了新的视角：

### 推荐书单

1. **《A Philosophy of Software Design》** — John Ousterhout
   - 深刻 vs 浅薄模块，每个程序员都该读

2. **《Designing Data-Intensive Applications》** — Martin Kleppmann
   - 终于读完了，分布式系统必读

3. **《The Design of Everyday Things》** — Don Norman
   - 用户体验不只是 UI 设计师的事

4. **《Working in Public》** — Nadia Eghbal
   - 开源社区的运作方式

5. **《Thinking, Fast and Slow》** — Daniel Kahneman
   - 理解自己的认知偏差

> 广度比深度更难积累，但最终影响更大。`, catLife, []string{tReading, tThinking}},

		{"CI/CD 入门：用 GitHub Actions 自动部署博客", "ci-cd-github-actions",
			"Push 代码自动测试、构建、部署到 VPS，省去手动操作。",
			`## 自动化部署

每次改完代码要手动 ` + "`go build`、`scp`" + ` 上传、重启服务，太麻烦。用 GitHub Actions 自动化：

` + "```yaml" + `
name: Deploy
on:
  push:
    branches: [main]
jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.26' }
      - run: cd backend && go build -o ../bin/server ./cmd/blog-server/
      - run: cd frontend && npm ci && npm run build
      - name: Deploy to VPS
        run: |
          scp bin/server user@host:/app/
          scp -r frontend/dist user@host:/app/
          ssh user@host "systemctl restart blog"
` + "```" + `

> 手动操作超过三次的事情，就值得自动化。`, catArch, []string{tCI, tDocker}},
	}

	for _, a := range articles {
		_, _ = postSvc.Create(ctx, post.CreateInput{
			Title: a.title, Slug: a.slug, Summary: a.summary,
			Content: a.content, CategoryID: a.cat,
			TagIDs: a.tags, Published: true,
		})
	}

	log.Println("seed: 10 articles, 4 categories, 12 tags")
}
