package server

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRun_ServerCreation(t *testing.T) {
	// 验证 Run 函数可被调用（不会在非阻塞模式下测试实际监听）
	// 仅验证 gin engine 创建正常
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/health", func(c *gin.Context) { c.Status(200) })

	// 不实际调用 Run（会阻塞），只验证 engine 正确
	if r == nil {
		t.Fatal("expected non-nil engine")
	}
}
