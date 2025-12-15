package uhttp

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/whosafe/uf/ucontext"
	"github.com/whosafe/uf/uprotocol/unet"
)

// TestResponseJSON 测试 JSON 响应
func TestResponseJSON(t *testing.T) {
	server := New()

	server.GET("/test", func(ctx *ucontext.Context, req unet.Request) error {
		data := map[string]string{"message": "hello"}
		return req.Response().JSON(200, data)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}
	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", w.Header().Get("Content-Type"))
	}
}

// TestResponseString 测试字符串响应
func TestResponseString(t *testing.T) {
	server := New()

	server.GET("/test", func(ctx *ucontext.Context, req unet.Request) error {
		return req.Response().String(200, "Hello, World!")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}
	if w.Body.String() != "Hello, World!" {
		t.Errorf("Expected body 'Hello, World!', got '%s'", w.Body.String())
	}
}

// TestResponseBytes 测试字节响应
func TestResponseBytes(t *testing.T) {
	server := New()

	testData := []byte("test data")
	server.GET("/test", func(ctx *ucontext.Context, req unet.Request) error {
		return req.Response().Bytes(200, testData)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}
	if !bytes.Equal(w.Body.Bytes(), testData) {
		t.Errorf("Expected body %v, got %v", testData, w.Body.Bytes())
	}
}
