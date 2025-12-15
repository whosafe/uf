package uhttp

import (
	"net/http/httptest"
	"testing"

	"github.com/whosafe/uf/ucontext"
	"github.com/whosafe/uf/uprotocol/unet"
)

// TestServerCreation 测试服务器创建
func TestServerCreation(t *testing.T) {
	// 测试默认创建
	server := New()
	if server == nil {
		t.Fatal("New() should return a server instance")
	}
	if server.config == nil {
		t.Error("Server config should not be nil")
	}
	if server.router == nil {
		t.Error("Server router should not be nil")
	}

	// 测试配置创建
	cfg := DefaultConfig()
	cfg.Name = "test-server"
	server2 := NewWithConfig(cfg)
	if server2 == nil {
		t.Fatal("NewWithConfig() should return a server instance")
	}
	if server2.config.Name != "test-server" {
		t.Errorf("Expected server name 'test-server', got '%s'", server2.config.Name)
	}
}

// TestRouteRegistration 测试路由注册
func TestRouteRegistration(t *testing.T) {
	server := New()

	// 注册路由
	called := false
	handler := func(ctx *ucontext.Context, req unet.Request) error {
		called = true
		return req.Response().JSON(200, map[string]string{"status": "ok"})
	}

	server.GET("/test", handler)

	// 创建测试请求
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// 执行请求
	server.ServeHTTP(w, req)

	// 验证
	if !called {
		t.Error("Handler should be called")
	}
	if w.Code != 200 {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}
}

// TestHTTPMethods 测试所有 HTTP 方法
func TestHTTPMethods(t *testing.T) {
	server := New()

	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			called := false
			handler := func(ctx *ucontext.Context, req unet.Request) error {
				called = true
				return req.Response().JSON(200, map[string]string{"method": method})
			}

			// 注册路由
			switch method {
			case "GET":
				server.GET("/"+method, handler)
			case "POST":
				server.POST("/"+method, handler)
			case "PUT":
				server.PUT("/"+method, handler)
			case "DELETE":
				server.DELETE("/"+method, handler)
			case "PATCH":
				server.PATCH("/"+method, handler)
			case "HEAD":
				server.HEAD("/"+method, handler)
			case "OPTIONS":
				server.OPTIONS("/"+method, handler)
			}

			// 创建测试请求
			req := httptest.NewRequest(method, "/"+method, nil)
			w := httptest.NewRecorder()

			// 执行请求
			server.ServeHTTP(w, req)

			// 验证
			if !called {
				t.Errorf("%s handler should be called", method)
			}
			if w.Code != 200 {
				t.Errorf("Expected status code 200, got %d", w.Code)
			}
		})
	}
}

// TestPathParams 测试路径参数
func TestPathParams(t *testing.T) {
	server := New()

	server.GET("/users/:id", func(ctx *ucontext.Context, req unet.Request) error {
		httpReq := req.(*Request)
		id := httpReq.Param("id")
		return req.Response().JSON(200, map[string]string{"id": id})
	})

	// 测试请求
	req := httptest.NewRequest("GET", "/users/123", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}
}

// TestQueryParams 测试查询参数
func TestQueryParams(t *testing.T) {
	server := New()

	server.GET("/search", func(ctx *ucontext.Context, req unet.Request) error {
		httpReq := req.(*Request)
		q := httpReq.Query("q")
		page := httpReq.QueryDefault("page", "1")
		return req.Response().JSON(200, map[string]string{
			"q":    q,
			"page": page,
		})
	})

	// 测试请求
	req := httptest.NewRequest("GET", "/search?q=test", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}
}

// TestRouteGroup 测试路由组
func TestRouteGroup(t *testing.T) {
	server := New()

	api := server.Group("/api")
	api.GET("/users", func(ctx *ucontext.Context, req unet.Request) error {
		return req.Response().JSON(200, map[string]string{"endpoint": "/api/users"})
	})

	// 测试请求
	req := httptest.NewRequest("GET", "/api/users", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}
}

// TestMiddleware 测试中间件
func TestMiddleware(t *testing.T) {
	server := New()

	// 中间件标记
	middlewareCalled := false
	handlerCalled := false

	// 注册中间件
	server.Use(func(next unet.HandlerFunc) unet.HandlerFunc {
		return func(ctx *ucontext.Context, req unet.Request) error {
			middlewareCalled = true
			return next(ctx, req)
		}
	})

	// 注册路由
	server.GET("/test", func(ctx *ucontext.Context, req unet.Request) error {
		handlerCalled = true
		return req.Response().JSON(200, map[string]string{"status": "ok"})
	})

	// 测试请求
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	// 验证
	if !middlewareCalled {
		t.Error("Middleware should be called")
	}
	if !handlerCalled {
		t.Error("Handler should be called")
	}
}

// Test404 测试 404 响应
func Test404(t *testing.T) {
	server := New()

	server.GET("/exists", func(ctx *ucontext.Context, req unet.Request) error {
		return req.Response().JSON(200, map[string]string{"status": "ok"})
	})

	// 测试不存在的路由
	req := httptest.NewRequest("GET", "/notfound", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Errorf("Expected status code 404, got %d", w.Code)
	}
}

// BenchmarkServerServeHTTP 性能测试
func BenchmarkServerServeHTTP(b *testing.B) {
	server := New()
	server.GET("/bench", func(ctx *ucontext.Context, req unet.Request) error {
		return req.Response().JSON(200, map[string]string{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/bench", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)
	}
}
