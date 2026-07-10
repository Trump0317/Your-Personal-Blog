# Your-Personal-Blog

用 Go 语言构建的极简个人博客系统。

## 技术栈

| 层 | 选型 |
|---|------|
| 语言 | Go 1.26 |
| HTTP 框架 | Gin |
| 数据库 | SQLite（`mattn/go-sqlite3`） |
| Markdown | `gomarkdown/markdown` |
| 前端 | 原生 HTML + CSS + JavaScript（零构建） |

## 快速开始

```bash
# 启动
go run ./cmd/blog-server/

# 浏览器打开
# 首页      http://localhost:8080
# 管理后台  http://localhost:8080/admin.html
# 默认账号  admin / password
```

### 配置

通过项目根目录的 `config.json` 文件配置，环境变量优先生效。

**配置文件** (`config.json`)：
```json
{
  "server":  { "port": ":8080" },
  "database": { "path": "blog.db" },
  "admin": { "username": "admin", "password": "password" },
  "site": {
    "name": "Your Personal Blog",
    "tagline": "探索技术、设计与思考的交汇",
    "description": "...",
    "hero_kicker": "Daily Engineering Notes",
    "hero_title1": "代码.",
    "hero_title2": "匠心.",
    "hero_title3": "文化.",
    "hero_desc": "...",
    "about_title": "关于 Your Personal Blog",
    "about_text": "...",
    "footer_note": "..."
  }
}
```

自定义配置文件路径: `BLOG_CONFIG=/path/to/my-config.json`

**环境变量覆盖**（优先级高于配置文件）：

| 变量 | 说明 |
|------|------|
| `BLOG_PORT` | 监听地址 |
| `BLOG_DB` | SQLite 数据库路径 |
| `BLOG_ADMIN_USER` | 管理员用户名 |
| `BLOG_ADMIN_PASS` | 管理员密码 |
| `BLOG_SITE_NAME` | 站点名称 |
| `BLOG_SITE_TAGLINE` | 站点标语 |
| … | (所有 site.* 字段均可覆盖) |

## 项目结构

```
cmd/blog-server/main.go       # 入口：依赖装配 + 启动
internal/
  post/                        # 文章模块
    post.go                    #   Post 模型 + Store 接口 + Service
    sqlite.go                  #   SQLite 实现
    handler.go                 #   HTTP Handler（公开/管理端路由）
  category/                    # 分类模块
    category.go
    sqlite.go
    handler.go
  tag/                         # 标签模块
    tag.go
    sqlite.go
    handler.go
  config/config.go             # 配置（环境变量）
  server/server.go             # HTTP Server 启动 + 优雅关闭
  middleware/                   # 中间件
    logger.go
    auth.go
web/
  templates/                   # HTML 页面
    index.html                 #   首页
    post.html                  #   文章详情
    admin.html                 #   管理后台
  static/                      # 静态资源
    style.css                  #   样式
    app.js                     #   首页逻辑
    post.js                    #   文章详情逻辑
    admin.js                   #   管理后台逻辑
```

## API

### 公开接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/posts` | 文章列表（`?page=`、`?category=`、`?tag=`、`?q=`） |
| GET | `/api/posts/:slug` | 文章详情（含 HTML 渲染） |
| GET | `/api/categories` | 分类列表 |
| GET | `/api/tags` | 标签列表 |

### 管理端接口（Basic Auth）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/admin/posts` | 文章列表（含草稿） |
| POST | `/api/admin/posts` | 创建文章 |
| GET | `/api/admin/posts/:id` | 文章详情 |
| PUT | `/api/admin/posts/:id` | 更新文章 |
| DELETE | `/api/admin/posts/:id` | 删除文章 |
| PUT | `/api/admin/posts/:id/publish` | 发布 |
| PUT | `/api/admin/posts/:id/unpublish` | 取消发布 |
| GET | `/api/admin/categories` | 分类列表 |
| POST | `/api/admin/categories` | 创建分类 |
| PUT | `/api/admin/categories/:id` | 更新分类 |
| DELETE | `/api/admin/categories/:id` | 删除分类 |
| GET | `/api/admin/tags` | 标签列表 |
| POST | `/api/admin/tags` | 创建标签 |
| DELETE | `/api/admin/tags/:id` | 删除标签 |

## 数据库

启动时自动创建 `blog.db`（SQLite），包含四张表：

```sql
posts (id, title, slug, content, html_content, summary, category_id, published, created_at, updated_at)
categories (id, name, slug)
tags (id, name, slug)
post_tags (post_id, tag_id)
```

首次启动自动写入种子数据。
