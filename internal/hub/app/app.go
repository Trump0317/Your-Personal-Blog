package app

import (
"context"
"fmt"
"os"
"os/signal"
"syscall"

config "github.com/ypb/your-personal-blog/config/hub"
"github.com/ypb/your-personal-blog/internal/hub/controller"
"github.com/ypb/your-personal-blog/internal/hub/repo/db"
"github.com/ypb/your-personal-blog/internal/hub/repo/storage"
"github.com/ypb/your-personal-blog/internal/hub/usercase"
"github.com/ypb/your-personal-blog/internal/hub/usercase/input"
"github.com/ypb/your-personal-blog/internal/pkg/httpserver"
"github.com/ypb/your-personal-blog/internal/pkg/logger"
)

func Run(cfg *config.Config) {
l := logger.New(cfg.Log.Level)

// 1. 初始化 Repository 层
fileStore, err := storage.NewLocalFileStore("uploads")
if err != nil {
l.Fatal(fmt.Errorf("app - Run - storage.NewLocalFileStore: %w", err))
}

fileRepo := db.NewMemoryFileRepo()
userRepo := db.NewMemoryUserRepo()

// 2. 初始化 UseCase 层
fileUC := usercase.NewFileUseCase(fileStore, fileRepo)
userUC := usercase.NewUserUseCase(userRepo)
adminUC := usercase.NewAdminUseCase(cfg, fileRepo, userRepo)

// 初始化测试用户 (临时代码)
_, _ = userUC.Create(context.Background(), input.UserCreate{APIKey: "default-test-key"})
l.Info("Initialized test user with APIKey: default-test-key")

// 3. 初始化 HTTP Server 和 Router
handler := httpserver.New(l)
controller.NewRouter(handler.Engine, cfg, l, fileUC, userUC, adminUC)

// 4. 启动服务
handler.Start()

// 5. 等待停止信号
interrupt := make(chan os.Signal, 1)
signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

select {
case s := <-interrupt:
l.Info("app - Run - signal: " + s.String())
case err = <-handler.Notify():
l.Error(fmt.Errorf("app - Run - handler.Notify: %w", err))
}

// 6. 优雅关闭
err = handler.Shutdown()
if err != nil {
l.Error(fmt.Errorf("app - Run - handler.Shutdown: %w", err))
}
}
