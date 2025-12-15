package main

import (
	"time"

	"github.com/whosafe/uf/ucontext"
	"github.com/whosafe/uf/uprotocol/uhttp"
	"github.com/whosafe/uf/uprotocol/unet"
)

func main() {
	server := uhttp.New()

	// 注册全局中间件
	server.Use(uhttp.MiddlewareTrace())                           // 链路追踪
	server.Use(uhttp.MiddlewareLogger())                          // 请求日志
	server.Use(uhttp.MiddlewareRecovery())                        // 异常恢复
	server.Use(uhttp.MiddlewareCORS())                            // 跨域支持
	server.Use(uhttp.MiddlewareRateLimitByIP(100, 1*time.Minute)) // IP 限流
	server.Use(uhttp.MiddlewareTimeout(30 * time.Second))         // 超时控制

	// 路由
	server.GET("/", indexHandler)
	server.GET("/health", healthHandler)

	// 启动服务器
	server.Start(":8080")
}

func indexHandler(ctx *ucontext.Context, req unet.Request) error {
	return req.Response().JSON(200, map[string]string{
		"message": "Server with middlewares",
	})
}

func healthHandler(ctx *ucontext.Context, req unet.Request) error {
	return req.Response().JSON(200, map[string]any{
		"status": "healthy",
		"time":   time.Now().Unix(),
	})
}
