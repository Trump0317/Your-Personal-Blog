package main

import (
	config "github.com/ypb/your-personal-blog/config/blog"
	"github.com/ypb/your-personal-blog/internal/blog/app"
)

func main() {
	// 1. 获取配置
	cfg := &config.Config{
		Admin: config.Admin{
			Username: "admin",
			Password: "password123",
		},
	}

	// 2. 运行 App (内部包含初始化、启动、等待信号及优雅关闭)
	app.Run(cfg)
}
