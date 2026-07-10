# 软件需求规格说明书（SRS）
# Your Personal Blog (YPB)
**版本**：V2.0
**状态**：已发布
**日期**：2026-07-10

---

## 1. 文档变更记录

| 日期 | 版本 | 说明 |
|------|------|------|
| 2026-04-13 | V1.4 | 初稿（存算分离架构，已废弃） |
| 2026-07-10 | V2.0 | 基于实际开发修订：单体模块化架构 |

---

## 2. 引言
### 2.1 目的
本文档为 **Your Personal Blog (YPB)** 轻量级博客系统的软件需求规格说明书，定义系统功能、非功能需求和验收标准。

### 2.2 范围
- 后端 API 服务（Go + Gin + SQLite）
- 前端 SPA（Vue 3 + Vite）
- 命令行工具（CLI）
- 容器化部署

### 2.3 术语
| 术语 | 说明 |
|------|------|
| YPB | Your Personal Blog |
| SPA | 单页应用，Vue Router 前端路由 |
| Basic Auth | HTTP 基本认证，管理端口鉴权 |

---

## 3. 总体描述
### 3.1 产品定位
面向个人开发者的**单体模块化**博客系统。一个二进制 + 一个 SQLite 文件即可运行，零外部依赖。

### 3.2 架构原则
```
┌─────────────────────────────────────────┐
│                blog-server               │
│  ┌───────┐ ┌──────────┐ ┌──────────┐   │
│  │ post  │ │ category │ │   tag    │   │
│  │ module│ │ module   │ │  module  │   │
│  └───────┘ └──────────┘ └──────────┘   │
│  ┌──────────┐ ┌─────────────┐           │
│  │ config   │ │ middleware  │           │
│  │ (json)   │ │ (auth/log)  │           │
│  └──────────┘ └─────────────┘           │
│  ┌──────────────────────┐               │
│  │    Vue 3 SPA (dist/) │               │
│  └──────────────────────┘               │
│              ↕ SQLite (blog.db)          │
└─────────────────────────────────────────┘
```

- **模块化打包**：每个业务域（post/category/tag）自包含模型、存储、服务、HTTP handler
- **接口在使用方定义**：Go 隐式接口，适配器在组装层实现
- **零配置启动**：首次运行自动建表 + 种子数据

---

## 4. 功能需求
### 4.1 访客端（公开访问）
| 功能 | 状态 | 说明 |
|------|------|------|
| 首页 | ✅ | 文章列表 + 分页 + 分类卡片 + 站点统计 |
| 文章详情 | ✅ | Markdown 渲染 HTML + 分类/标签导航 |
| 归档页 | ✅ | 年/月统计 + 按年月筛选文章列表 |
| 分类页 | ✅ | 分类列表 + 按分类浏览文章 |
| 标签页 | ✅ | 标签列表 + 按标签浏览文章 |
| 关于页 | ✅ | 静态页面，内容通过 config.json 定制 |
| RSS | ✅ | RSS 2.0，动态域名 |
| 搜索 | ✅ | 标题/内容 LIKE 查询 |
| 暗色模式 | ✅ | CSS 变量 + localStorage 持久化 |

### 4.2 作者端（Basic Auth 保护）
| 功能 | 状态 | 说明 |
|------|------|------|
| 登录 | ✅ | Basic Auth + localStorage 持久化 |
| 文章 CRUD | ✅ | 创建/编辑/删除/发布/取消发布 |
| Markdown 编辑器 | ✅ | 分栏实时预览 + 工具栏 + 快捷键 |
| 草稿自动保存 | ✅ | 3 秒 debounce → localStorage |
| 图片上传 | ✅ | multipart → `uploads/` 目录 |
| 分类管理 | ✅ | 增删改 |
| 标签管理 | ✅ | 增删改（去重创建） |
| 筛选 | ✅ | 搜索 + 状态（已发布/草稿）+ 分类 |

### 4.3 自定义配置
| 功能 | 状态 | 说明 |
|------|------|------|
| 配置文件 | ✅ | `config.json`（端口/数据库/管理员/站点文案） |
| 环境变量覆盖 | ✅ | 所有配置项均支持 `BLOG_*` 环境变量 |
| 前端文案 | ✅ | 11 个字段（站点名/标语/Hero/关于/页脚） |

### 4.4 CLI 命令行工具（待开发）
```
ypb login          → 鉴权并保存 Token
ypb new <title>    → 快速创建文章
ypb push <file>    → 推送本地 Markdown 到博客
ypb pull           → 拉取远程文章到本地
ypb list           → 列出文章
ypb publish <id>   → 发布文章
```

### 4.5 待开发功能
| 优先级 | 功能 |
|--------|------|
| P0 | Argon2 密码哈希（替代明文 Basic Auth） |
| P1 | JWT Token 认证 |
| P2 | 编辑器支持 Mermaid 图表 |
| P2 | 编辑器支持 KaTeX 数学公式 |
| P3 | 文章定时发布 |
| P3 | 自动生成文章目录（TOC） |

---

## 5. 非功能需求
### 5.1 性能
- 单页渲染 < 500ms（本地 SQLite）
- 前端构建产物总大小 < 150KB（gzip 后约 40KB）

### 5.2 安全
- 管理接口 Basic Auth 保护
- 后续版本：管理员密码 Argon2 哈希存储
- 后续版本：管理接口强制 HTTPS

### 5.3 部署
- 单个二进制 33MB
- 前端静态文件内嵌 dist/ 目录
- SQLite 零配置
- Dockerfile（多阶段构建）

### 5.4 项目规格
| 指标 | 数值 |
|------|------|
| Go 代码 | ~3,500 行 |
| 前端代码 | ~1,800 行 |
| 测试用例 | 52 个 |
| 测试覆盖率 | config/middleware 100%，业务模块 ~50% |
| 外部依赖 | gin, gomarkdown, uuid, go-sqlite3, vue, vue-router (6 个) |

---

## 6. API 接口

### 6.1 公开接口
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/posts` | 文章列表（?page/size/category_id/tag_id/year/month/q） |
| GET | `/api/posts/:slug` | 文章详情（仅已发布） |
| GET | `/api/categories` | 分类列表（含文章数） |
| GET | `/api/categories/:slug` | 分类详情 |
| GET | `/api/tags` | 标签列表（含文章数） |
| GET | `/api/tags/:slug` | 标签详情 |
| GET | `/api/stats` | 站点统计 |
| GET | `/api/archive` | 年月归档 |
| GET | `/api/site` | 站点配置 |
| GET | `/rss.xml` | RSS 2.0 |

### 6.2 管理端接口（Basic Auth）
| 方法 | 路径 | 说明 |
|------|------|------|
| GET/POST | `/api/admin/posts` | 文章列表/创建 |
| GET/PUT/DELETE | `/api/admin/posts/:id` | 文章详情/更新/删除 |
| PUT | `/api/admin/posts/:id/publish` | 发布 |
| PUT | `/api/admin/posts/:id/unpublish` | 取消发布 |
| GET/POST/PUT/DELETE | `/api/admin/categories[/:id]` | 分类 CRUD |
| GET/POST/PUT/DELETE | `/api/admin/tags[/:id]` | 标签 CRUD |
| POST | `/api/upload` | 图片上传 |

---

## 7. 验收标准
### 7.1 功能验收
- [x] 所有公开接口可用
- [x] 所有管理端接口可用，Basic Auth 正确拦截
- [x] 前端 8 个页面全部可访问，SPA 路由正确
- [x] 种子数据首次启动自动写入
- [x] config.json 自定义文案生效
- [x] 暗色模式切换正常

### 7.2 质量验收
- [x] `go vet` 零警告
- [x] 52 个单元测试全部通过
- [x] 前端 `npm run build` 零错误
- [x] 无 TODO/FIXME 遗留

### 7.3 部署验收
- [x] `make build` 一键构建
- [x] `./bin/blog-server` 直接启动
- [ ] Docker 部署（待实现）

---

## 8. 项目结构
```
Your-Personal-Blog/
├── config.json                  # 站点配置
├── Makefile                     # dev/build/run/clean
├── backend/
│   ├── cmd/blog-server/main.go  # 入口 + 路由 + 种子数据 + RSS
│   └── internal/
│       ├── post/                # 文章模块
│       │   ├── post.go          #   模型 + Store接口 + Service
│       │   ├── sqlite.go        #   SQLite 实现
│       │   └── handler.go       #   HTTP Handler
│       ├── category/            # 分类模块 (同上结构)
│       ├── tag/                 # 标签模块 (同上结构)
│       ├── config/              # 配置 (JSON文件 + 环境变量)
│       ├── server/              # HTTP 启动 + 优雅关闭
│       └── middleware/          # Logger + BasicAuth
├── frontend/
│   ├── src/
│   │   ├── App.vue              # 根组件 (导航/页脚/暗色切换)
│   │   ├── views/               # 8 个页面组件
│   │   ├── api/index.js         # API 客户端封装
│   │   ├── composables/         # useDark / useToast / useSite
│   │   ├── router/index.js      # Vue Router
│   │   └── style.css            # 全局样式 + 暗色变量
│   └── dist/                    # 构建产物
└── doc/
    └── SRS.md                   # 本文档
```
