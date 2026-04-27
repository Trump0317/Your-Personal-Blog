package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ypb/your-personal-blog/internal/hub/handler"
	"github.com/ypb/your-personal-blog/internal/hub/repository/sqlite"
	"github.com/ypb/your-personal-blog/internal/hub/service"
	"github.com/ypb/your-personal-blog/internal/hub/storage"
)

func main() {
	log.Println("Storage Hub 正在启动...")

	// 1. 初始化最底层的基础设施 (即使是空实现或 LocalStorage)
	storageEngine, _ := storage.NewLocalStorage("./uploads")
	fileRepo := sqlite.NewFileRepository(nil) // 传入 nil 作为占位符，实际应用中应传入 *sql.DB 实例
	clientRepo := sqlite.NewClientRepository(nil)

	// 2. 注入依赖到 Service 层
	authSvc := service.NewAuthService(clientRepo)
	fileSvc := service.NewFileService(fileRepo, storageEngine)

	// 3. 注入 Service 到 Handler 层
	fileHandler := handler.NewFileHandler(fileSvc, authSvc)

	// 4. 设置路由还启动
	r := gin.Default()
	fileHandler.RegisterRoutes(r)

	log.Println("正在监听端口 :8080...")
	r.Run(":8080")
}
