package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	config "github.com/ypb/your-personal-blog/config/blog"
	blogHttp "github.com/ypb/your-personal-blog/internal/blog/controller/http"
	"github.com/ypb/your-personal-blog/internal/blog/model"
	"github.com/ypb/your-personal-blog/internal/blog/repo"
	"github.com/ypb/your-personal-blog/internal/blog/usecase"
	"github.com/ypb/your-personal-blog/internal/pkg/httpserver"
	"github.com/ypb/your-personal-blog/internal/pkg/logger"
)

func Run(cfg *config.Config) {
	l := logger.New("debug")

	// 1. 初始化 Repository 层
	postRepo := repo.NewMemoryPostRepo()
	tagRepo := repo.NewMemoryTagRepo()
	catRepo := repo.NewMemoryCategoryRepo()

	// 2. 初始化 UseCase 层
	tagUC := usecase.NewTagUsecase(tagRepo)
	catUC := usecase.NewCategoryUsecase(catRepo)
	postUC := usecase.NewPostUsecase(postRepo, catRepo)

	// 执行演示数据初始化
	seedData(postUC, catUC, tagUC)

	// 3. 初始化 HTTP Server 和 Router
	handler := httpserver.New(l)
	blogHttp.NewRouter(handler.Engine, cfg, l, postUC, tagUC, catUC)

	// 4. 启动服务
	handler.Start()

	// 5. 等待停止信号 (Graceful Shutdown)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err := <-handler.Notify():
		l.Error(fmt.Errorf("app - Run - handler.Notify: %w", err))
	}

	// 6. 优雅关闭
	err := handler.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - handler.Shutdown: %w", err))
	}
}

func seedData(postUC usecase.Post, catUC usecase.Category, tagUC usecase.Tag) {
	ctx := context.Background()

	// 演示分类
	c1ID, _ := catUC.Create(ctx, "技术分享", "tech")
	c2ID, _ := catUC.Create(ctx, "随笔感悟", "blog")

	// 演示文章 1: Markdown 丰富内容测试
	t1, _ := tagUC.GetOrCreates(ctx, []string{"Go", "Vue", "Backend"})
	postUC.Create(ctx, &usecase.PostCreateInput{
		Title:   "你好，世界 - Go 语言博客系统",
		Slug:    "hello-world",
		Summary: "本博客系统采用 Clean Architecture 架构，使用 Go 语言开发。支持 Markdown 渲染、分类管理和标签系统。",
		Content: `## 欢迎来到我的博客

这是一篇测试文章，用于展示 **Markdown** 的渲染效果。

### 技术栈
* **Go** (Gin / Clean Architecture)
* **Vue 3** (前端驱动)
* **Tailwind CSS** (样式美化)
* **Marked.js** (Markdown 解析)

### 代码示例
` + "```go" + `
package main

import "fmt"

func main() {
    fmt.Println("Hello, Blog!")
}
` + "```" + `

> 这是一段引用文本，测试引用样式。
`,
		CategoryID: c1ID,
		TagIDs:     t1,
	})

	// 演示文章 2: 简单内容
	t2, _ := tagUC.GetOrCreates(ctx, []string{"Design", "Notes"})
	postUC.Create(ctx, &usecase.PostCreateInput{
		Title:   "关于本项目的设计初衷",
		Slug:    "about-design",
		Summary: "为什么要写这样一个博客系统？不仅仅是为了记录，更是为了实践。",
		Content: `### 为什么 es Go？
Go 语言具有优秀的并发处理能力和简洁的语法，非常适合编写轻量级的后端服务。

### 设计原则
1. **模块化**: 核心业务逻辑与前端/数据库解耦。
2. **易用性**: 无需复杂的配置，开箱即用。
3. **扩展性**: 方便后续集成真实数据库（如 MySQL/PostgreSQL）。`,
		CategoryID: c1ID,
		TagIDs:     t2,
	})

	// 演示文章 3: 生活随笔
	t3, _ := tagUC.GetOrCreates(ctx, []string{"生活", "心情"})
	postUC.Create(ctx, &usecase.PostCreateInput{
		Title:      "周末午后的咖啡馆",
		Slug:       "coffee-afternoon",
		Summary:    "找一个静谧的角落，写写代码，看看窗外。",
		Content:    "阳光洒进窗户，桌上的拿铁冒着热气。有时候慢下来也是一种进步。",
		CategoryID: c2ID,
		TagIDs:     t3,
	})

	// 将所有文章设为已发布
	posts, _, _ := postUC.List(ctx, &usecase.PostListInput{Page: 1, PageSize: 10})
	for _, p := range posts {
		postUC.Publish(ctx, p.ID)
	}

	pubStatus := model.PostPublished
	pubPosts, total, _ := postUC.List(ctx, &usecase.PostListInput{Page: 1, PageSize: 10, Status: &pubStatus})
	fmt.Printf("Initialization successful: %d posts are now published (total items: %d)\n", len(pubPosts), total)
}
