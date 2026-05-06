package main

import (
	config "github.com/ypb/your-personal-blog/config/hub"
	"github.com/ypb/your-personal-blog/internal/hub/app"
)

func main() {
	// 简单模拟配置加载，这里硬编码管理员信息
	cfg := &config.Config{
		Log: config.Log{
			Level: "debug",
		},
		Admin: config.Admin{
			Username: "admin",
			Password: "admin-password",
			APIKey:   "super-admin-key",
		},
	}

	app.Run(cfg)
}
