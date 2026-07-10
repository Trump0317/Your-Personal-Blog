package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestBasicAuth_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(BasicAuth("admin", "password"))
	r.GET("/", func(c *gin.Context) { c.Status(200) })

	req := httptest.NewRequest("GET", "/", nil)
	req.SetBasicAuth("admin", "password")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("status = %d, want 200", w.Code)
	}
}

func TestBasicAuth_WrongPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(BasicAuth("admin", "password"))
	r.GET("/", func(c *gin.Context) { c.Status(200) })

	req := httptest.NewRequest("GET", "/", nil)
	req.SetBasicAuth("admin", "wrong")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("status = %d, want 401", w.Code)
	}
}

func TestBasicAuth_NoAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(BasicAuth("admin", "password"))
	r.GET("/", func(c *gin.Context) { c.Status(200) })

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("status = %d, want 401", w.Code)
	}
	if w.Header().Get("WWW-Authenticate") == "" {
		t.Error("expected WWW-Authenticate header")
	}
}

func TestBasicAuth_WrongUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(BasicAuth("admin", "password"))
	r.GET("/", func(c *gin.Context) { c.Status(200) })

	req := httptest.NewRequest("GET", "/", nil)
	req.SetBasicAuth("hacker", "password")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("status = %d, want 401", w.Code)
	}
}
