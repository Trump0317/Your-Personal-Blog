# 架构说明

## 设计理念：模块化优先

本项目采用 **Go 社区常见的模块化分包** 方式，而非 Clean Architecture。核心原则：

> **每个包是一个独立的功能模块，自包含模型、存储、业务逻辑和 HTTP 处理。**

### 与 Clean Architecture 的对比

```
Clean Architecture（旧版 dev 分支）        模块化（新版 modular 分支）
────────────────────────────────         ────────────────────────────
model/                                    post/
  post.go                                   post.go       ← 模型 + 接口 + 业务 + handler
  tag.go                                    sqlite.go
  category.go                               handler.go
repo/                                     category/
  contracts.go    ← 接口集中定义              category.go
  memory_repo.go                            sqlite.go
usecase/                                    handler.go
  contracts.go    ← 接口重复定义             tag/
  post.go         ← 业务逻辑                  tag.go
  input.go        ← DTO                      sqlite.go
  output.go       ← DTO                      handler.go
controller/                               config/
  http/                                     config.go
    visitor.go                            middleware/
    admin.go                                auth.go
    router.go                               logger.go
    middleware.go                          server/
    request.go                              server.go
    response.go
```

| 对比维度 | Clean Architecture | 模块化 |
|---------|-------------------|--------|
| 文件数（核心功能） | ~15 个 | 10 个 |
| 包数 | 7 层 | 6 个模块 |
| 改一个功能跨几个包 | 4（model→repo→usecase→controller） | 1（post/） |
| 接口定义位置 | 独立 contracts.go | 使用时定义 |
| DTO 层 | 有（input.go/output.go） | 无（直接用模型） |
| 新人上手 | 需理解分层约定 | 打开包即看懂 |

---

## 模块详解

### `internal/post/` — 文章模块

```
post.go     定义 Post 模型、Store 接口、Service 业务逻辑、依赖接口
sqlite.go   Store 接口的 SQLite 实现（含分页、过滤、关联查询）
handler.go  HTTP Handler（公开路由 + 管理端路由）
```

**依赖方向**：`post` 包定义它需要什么接口（`CategoryStore`、`TagStore`），由 `main.go` 提供实现。

```go
// post.go — post 包定义自己需要的接口
type CategoryStore interface {
    Get(ctx context.Context, id string) (*CategoryRef, error)
    List(ctx context.Context) ([]*CategoryRef, error)
}

type TagStore interface {
    GetByIDs(ctx context.Context, ids []string) ([]TagRef, error)
}
```

### `internal/category/` — 分类模块

独立模块，不依赖其他业务包。`main.go` 通过适配器将其注入 `post.Service`。

### `internal/tag/` — 标签模块

同上。`tag.Service` 提供 `CreateOrGet` 方法（去重创建）。

### 适配器

`category.Service` 和 `tag.Service` 的返回值类型与 `post` 包期望的接口不同。`main.go` 中定义轻薄适配器做类型转换：

```go
// main.go — 组装层（Composition Root）
type categoryAdapter struct { svc *category.Service }

func (a *categoryAdapter) Get(ctx context.Context, id string) (*post.CategoryRef, error) {
    c, err := a.svc.Get(ctx, id)
    if err != nil { return nil, err }
    return &post.CategoryRef{ID: c.ID, Name: c.Name, Slug: c.Slug}, nil
}
```

这是 Go 社区的常见模式：**接口在使用者一侧定义，适配器在组装层实现。**

---

## 请求流程

以「访问者查看文章」为例：

```
浏览器                    Gin Router                post.Handler           post.Service         post.SQLiteStore
  │                          │                          │                      │                     │
  │ GET /api/posts/hello     │                          │                      │                     │
  │─────────────────────────►│                          │                      │                     │
  │                          │ ListPublic(c)            │                      │                     │
  │                          │─────────────────────────►│                      │                     │
  │                          │                          │ List(ctx, filter)    │                     │
  │                          │                          │─────────────────────►│                     │
  │                          │                          │                      │ List(ctx, filter)   │
  │                          │                          │                      │────────────────────►│
  │                          │                          │                      │                     │ SELECT * FROM posts
  │                          │                          │                      │ ◄────────────────────│
  │                          │                          │                      │                     │ + post_tags JOIN
  │                          │                          │ ◄─────────────────────│                     │
  │                          │                          │                      │                     │
  │                          │                          │ fillRefs(ctx, post)  │                     │
  │                          │                          │──────┐               │                     │
  │                          │                          │      │ catAdapter.Get│                     │
  │                          │                          │      │ tagAdapter.Get│                     │
  │                          │                          │◄─────┘               │                     │
  │                          │                          │                      │                     │
  │                          │ ◄─── JSON ───────────────│                      │                     │
  │ ◄─── 200 OK ────────────│                          │                      │                     │
```

每一步都是单一职责，没有跨越式的数据访问。

---

## 设计决策记录

### 1. 为什么不用接口集中定义？

Go 的接口是 **隐式满足** 的。定义在使用方、实现在提供方，是 Go 社区的主流做法。`post` 包说「我需要能查分类的接口」，`main.go` 给一个适配器就行——`category` 包完全不知道 `post` 包的存在。

### 2. 为什么不用 DTO？

博客系统 CRUD 逻辑简单，`Post` 结构体直接用 `json:"..."` tag 序列化就够了。当模型字段和 API 响应差异变大时（比如需要计算字段、权限裁剪），再引入 `Output` 类型不迟。

### 3. 为什么 SQLite 直接实现而不写内存版？

内存版只为演示和测试存在。既然 SQLite 可以零配置跑起来（一个文件），就没有理由先写内存版再替换。测试可以用 `:memory:` 模式的 SQLite。

### 4. 为什么前端零构建？

个人博客不需要 SPA 框架。原生 HTML + fetch API + 几百行 JS 足够。构建工具链是维护负担，除非项目复杂度证明它值得。

### 5. 为什么包不叫 `service` 而叫 `post`/`category`/`tag`？

包的命名应该反映 **领域概念**，而非 **技术角色**。打开 `internal/post/` 就知道这里管文章，打开 `internal/category/` 就知道这里管分类。如果用 `service`、`repository`、`model` 这种技术分层名，需要跳三个目录才能理解一个功能。

---

## 未来演进方向

当项目膨胀到一定程度时，自然会需要：

1. **DTO 层** — 当 API 响应和数据库模型差异过大时
2. **结构化日志** — 当需要 ELK/ Grafana 时
3. **JWT 认证** — 当需要多用户/权限控制时
4. **配置热加载** — 当运维需要不重启改配置时
5. **前端构建** — 当前端交互复杂度超过原生 JS 舒适区时

但 **不应该在需要之前实现它们**。模块化架构的优点是：新增功能只需改一个包，不会像分层架构那样需要跨 4 层修改。
