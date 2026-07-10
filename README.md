# Your Personal Blog

单体模块化个人博客系统。一个二进制 + 一个 SQLite 文件，开箱即用。

## 特性

- **零配置启动** — 首次运行自动建表并写入种子数据
- **Markdown 编辑器** — 分栏实时预览 + 工具栏 + 快捷键 + 草稿自动保存
- **Mermaid 图表** — 编辑器/文章页渲染流程图、时序图等
- **KaTeX 公式** — 行内和块级数学公式渲染
- **暗色模式** — CSS 变量切换，偏好持久化
- **RSS 2.0** — 自动生成，动态域名
- **CLI 工具** — 终端 login/push/pull/list/publish
- **JWT 认证** — 管理员密码 Argon2 哈希，7 天 token
- **配置灵活** — `config.json` 一站式配置，环境变量可覆盖

## 技术栈

| 层 | 选型 |
|---|------|
| 语言 | Go 1.24+ |
| HTTP 框架 | Gin |
| 数据库 | SQLite（`mattn/go-sqlite3`） |
| Markdown | `gomarkdown/markdown`（后端渲染） |
| 前端 | Vue 3 + Vite |
| 认证 | JWT + Argon2 |
| 图表/公式 | Mermaid + KaTeX（CDN） |

## 快速开始

```bash
# 构建并启动
make run

# 浏览器打开
# 首页      http://localhost:8080
# 管理后台  http://localhost:8080/admin
# CLI 工具  ./bin/ypb login
```

## 配置

编辑项目根目录 `config.json`：

```json
{
  "server":  { "port": ":8080" },
  "database": { "path": "blog.db" },
  "admin": {
    "username": "admin",
    "password": "password",
    "jwt_secret": ""
  },
  "site": {
    "name": "Your Personal Blog",
    "tagline": "探索技术、设计与思考的交汇",
    "description": "探索技术、设计与思考的交汇。以极简致敬繁复。",
    "hero_kicker": "Daily Engineering Notes",
    "hero_title1": "代码.",
    "hero_title2": "匠心.",
    "hero_title3": "文化.",
    "hero_desc": "探索软件架构、工程美学与人类体验的交汇点。以工匠精神，雕琢每一行代码。",
    "about_title": "关于 Your Personal Blog",
    "about_text": "这是一个用 Go + SQLite + Vue 构建的个人博客。\n写博客的目的：记录思考，沉淀知识，分享实践。",
    "footer_note": "探索技术、设计与思考的交汇"
  }
}
```

所有 `site.*` 字段均可通过同名环境变量覆盖（如 `BLOG_SITE_NAME`）。
自定义配置文件路径：`BLOG_CONFIG=/path/to/config.json`

## 部署

### 单机运行

```bash
make run
```

### Docker

```bash
make docker-build
make docker-up
```

## 项目结构

```
├── config.json                  # 站点配置
├── Makefile
├── Dockerfile
├── docker-compose.yml
├── backend/
│   ├── cmd/
│   │   ├── blog-server/main.go  # 入口 + 路由 + RSS + 种子数据
│   │   └── ypb/main.go          # CLI 工具
│   └── internal/
│       ├── post/                # 文章模块（模型/存储/服务/接口）
│       ├── category/            # 分类模块
│       ├── tag/                 # 标签模块
│       ├── config/              # 配置加载（JSON + 环境变量）
│       ├── server/              # HTTP 启动 + 优雅关闭
│       └── middleware/          # Logger / JWT 认证
├── frontend/
│   ├── src/
│   │   ├── App.vue              # 根组件（导航/页脚/暗色切换）
│   │   ├── views/               # 8 个页面
│   │   ├── api/index.js         # API 客户端
│   │   ├── composables/         # useDark / useToast / useSite / useEnhance
│   │   └── router/index.js      # Vue Router
│   └── dist/                    # 构建产物
└── doc/SRS.md                   # 需求规格说明书
```

## API

### 公开接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/posts` | 文章列表（`?page/size/category_id/tag_id/year/month/q`） |
| GET | `/api/posts/:slug` | 文章详情 |
| GET | `/api/categories` | 分类列表（含文章数） |
| GET | `/api/categories/:slug` | 分类详情 |
| GET | `/api/tags` | 标签列表（含文章数） |
| GET | `/api/tags/:slug` | 标签详情 |
| GET | `/api/stats` | 站点统计 |
| GET | `/api/archive` | 年月归档 |
| GET | `/api/site` | 站点配置（供前端使用） |
| GET | `/rss.xml` | RSS 2.0 |

### 管理端接口（JWT 保护）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/admin/login` | 登录获取 token |
| GET | `/api/admin/posts` | 文章列表（含草稿，`?page/size`） |
| POST | `/api/admin/posts` | 创建文章 |
| GET | `/api/admin/posts/:id` | 文章详情 |
| PUT | `/api/admin/posts/:id` | 更新文章 |
| DELETE | `/api/admin/posts/:id` | 删除文章 |
| PUT | `.../publish` / `.../unpublish` | 发布/取消发布 |
| GET/POST/PUT/DELETE | `/api/admin/categories[/:id]` | 分类 CRUD |
| GET/POST/PUT/DELETE | `/api/admin/tags[/:id]` | 标签 CRUD |
| POST | `/api/upload` | 图片上传 |

## CLI

```bash
./bin/ypb login               # 登录并保存 token
./bin/ypb list                # 列出文章
./bin/ypb new "标题"          # 创建草稿
./bin/ypb push article.md     # 推送本地 Markdown（支持 frontmatter）
./bin/ypb pull                # 拉取文章到 _posts/
./bin/ypb publish <id>        # 发布
./bin/ypb delete <id>         # 删除
```

## 数据库

```sql
posts       (id, title, slug, content, html_content, summary, category_id, published, created_at, updated_at)
categories  (id, name, slug)
tags        (id, name, slug)
post_tags   (post_id, tag_id)
```

启动时自动创建。首次启动写入 10 篇种子文章、4 个分类、11 个标签。

## 开发

```bash
make dev     # 后端 :8080 + 前端 :3000 (热更新)
make build   # 生产构建
make test    # 运行测试
make clean   # 清理
```
