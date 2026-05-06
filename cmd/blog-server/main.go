package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"

	config "github.com/ypb/your-personal-blog/config/blog"
	blogHttp "github.com/ypb/your-personal-blog/internal/blog/controller/http"
	"github.com/ypb/your-personal-blog/internal/blog/repo/db"
	"github.com/ypb/your-personal-blog/internal/blog/usecase"
	"github.com/ypb/your-personal-blog/internal/pkg/logger"
)

func main() {
	// 1. 初始化数据库 (使用 SQLite 内存模式方便演示)
	sqldb, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		log.Fatal(err)
	}
	defer sqldb.Close()

	// 2. 创建表结构
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		slug TEXT NOT NULL,
		content TEXT NOT NULL,
		html_content TEXT,
		summary TEXT,
		cover_image TEXT,
		status INTEGER DEFAULT 0,
		is_top BOOLEAN DEFAULT 0,
		view_count INTEGER DEFAULT 0,
		category_id INTEGER,
		created_at DATETIME,
		updated_at DATETIME
	);
	INSERT INTO posts (title, slug, content, status, created_at, updated_at) 
	VALUES ('我的第一篇博客', 'hello-world', '这是内容...', 1, datetime('now'), datetime('now'));
	`
	_, err = sqldb.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	// 3. 初始化组件
	l := logger.New("debug")
	postRepo := db.NewPostRepo(sqldb)
	postUC := usecase.NewPostUsecase(postRepo)

	// 4. 初始化 Gin 路由
	r := gin.Default()

	cfg := &config.Config{} // 暂时传空配置
	blogHttp.NewRouter(r, cfg, l, postUC)

	// 5. 启动服务
	log.Println("Server running on http://localhost:8080")
	log.Println("Access http://localhost:8080/blog/visitor/index to view the blog")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
