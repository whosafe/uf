package uhttp

import (
	"net/http/httptest"
	"testing"

	"github.com/whosafe/uf/ucontext"
	"github.com/whosafe/uf/uprotocol/unet"
)

// TestRouterBasicRoute 测试基础路由
func TestRouterBasicRoute(t *testing.T) {
	server := New()

	called := false
	server.GET("/users", func(ctx *ucontext.Context, req unet.Request) error {
		called = true
		return req.Response().JSON(200, map[string]string{"endpoint": "/users"})
	})

	req := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	if !called {
		t.Error("Route handler should be called")
	}
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestRouterPathParams 测试路径参数
func TestRouterPathParams(t *testing.T) {
	server := New()

	var capturedID string
	server.GET("/users/:id", func(ctx *ucontext.Context, req unet.Request) error {
		httpReq := req.(*Request)
		capturedID = httpReq.Param("id")
		return req.Response().JSON(200, map[string]string{"id": capturedID})
	})

	req := httptest.NewRequest("GET", "/users/123", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	if capturedID != "123" {
		t.Errorf("Expected id=123, got %s", capturedID)
	}
}

// TestRouterMultipleParams 测试多个路径参数
func TestRouterMultipleParams(t *testing.T) {
	server := New()

	var userID, postID string
	server.GET("/users/:userId/posts/:postId", func(ctx *ucontext.Context, req unet.Request) error {
		httpReq := req.(*Request)
		userID = httpReq.Param("userId")
		postID = httpReq.Param("postId")
		return req.Response().JSON(200, map[string]string{
			"userId": userID,
			"postId": postID,
		})
	})

	req := httptest.NewRequest("GET", "/users/123/posts/456", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	if userID != "123" {
		t.Errorf("Expected userId=123, got %s", userID)
	}
	if postID != "456" {
		t.Errorf("Expected postId=456, got %s", postID)
	}
}

// TestRouterNotFound 测试路由未找到
func TestRouterNotFound(t *testing.T) {
	server := New()

	server.GET("/users", func(ctx *ucontext.Context, req unet.Request) error {
		return req.Response().JSON(200, map[string]string{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/posts", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

// TestRouterMethodNotAllowed 测试方法不匹配
func TestRouterMethodNotAllowed(t *testing.T) {
	server := New()

	server.GET("/users", func(ctx *ucontext.Context, req unet.Request) error {
		return req.Response().JSON(200, map[string]string{"status": "ok"})
	})

	req := httptest.NewRequest("POST", "/users", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

// TestRouterStaticPriority 测试静态路由优先级
func TestRouterStaticPriority(t *testing.T) {
	server := New()

	var endpoint string
	server.GET("/users/new", func(ctx *ucontext.Context, req unet.Request) error {
		endpoint = "static"
		return req.Response().JSON(200, map[string]string{"type": "static"})
	})

	server.GET("/users/:id", func(ctx *ucontext.Context, req unet.Request) error {
		endpoint = "dynamic"
		return req.Response().JSON(200, map[string]string{"type": "dynamic"})
	})

	// 测试静态路由
	req := httptest.NewRequest("GET", "/users/new", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	if endpoint != "static" {
		t.Error("Static route should have priority")
	}

	// 测试动态路由
	req = httptest.NewRequest("GET", "/users/123", nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)

	if endpoint != "dynamic" {
		t.Error("Dynamic route should match")
	}
}
