package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLogger(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Logger())
	r.GET("/test", func(c *gin.Context) { c.String(200, "ok") })

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("status = %d, want 200", w.Code)
	}
	// Logger 只做日志输出，不应影响请求/响应
}

func TestLogger_Panic(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(gin.Recovery()) // 防止 panic 传播
	r.Use(Logger())
	r.GET("/panic", func(c *gin.Context) { panic("test panic") })

	req := httptest.NewRequest("GET", "/panic", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 500 {
		t.Errorf("status = %d, want 500", w.Code)
	}
}

func TestLogger_StatusCodes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		path    string
		handler func(*gin.Context)
		want    int
	}{
		{"/200", func(c *gin.Context) { c.Status(200) }, 200},
		{"/301", func(c *gin.Context) { c.Redirect(301, "/") }, 301},
		{"/400", func(c *gin.Context) { c.Status(400) }, 400},
		{"/500", func(c *gin.Context) { c.Status(500) }, 500},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			r := gin.New()
			r.Use(Logger())
			r.GET(tt.path, tt.handler)

			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.want {
				t.Errorf("status = %d, want %d", w.Code, tt.want)
			}
		})
	}
}
