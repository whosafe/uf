package template

// internal/router/router.go 模板
const RouterTemplate = `package router

import (
	"{{.ModulePath}}/internal/handler"
	"{{.ModulePath}}/internal/middleware"
	"github.com/whosafe/uf/uprotocol/uhttp"
)

// New 注册所有路由
func New(server *uhttp.Server) {
	// 公开路由（无需认证）
	server.GET("/", handler.Index)
	server.GET("/health", handler.Health)

	// API 路由组（需要认证）
	api := server.Group("/api")
	api.Use(middleware.RequireAuth())
	{
		// 示例：需要认证的路由
		// api.GET("/profile", handler.Profile)
		// api.POST("/logout", handler.Logout)
	}

	// 管理员路由组（需要认证 + 管理员权限）
	// admin := server.Group("/admin", middleware.RequireAuth(), middleware.RequireAdmin())
	// {
	//     admin.GET("/users", handler.AdminUsers)
	//     admin.GET("/settings", handler.AdminSettings)
	// }
}
`

// internal/handler/index.go 模板
const HandlerIndexTemplate = `package handler

import (
	"time"

	"github.com/whosafe/uf/ucontext"
	"github.com/whosafe/uf/uprotocol/unet"
)

// Index 首页
func Index(ctx *ucontext.Context, req unet.Request) error {
	return req.Response().JSON(200, map[string]string{
		"message": "欢迎使用 {{.ProjectName}}",
		"version": "1.0.0",
	})
}

// Health 健康检查
func Health(ctx *ucontext.Context, req unet.Request) error {
	tc := ucontext.FromContext(ctx.Context())
	traceID := ""
	if tc != nil {
		traceID = tc.TraceID
	}

	return req.Response().JSON(200, map[string]any{
		"status":   "ok",
		"time":     time.Now().Format(time.RFC3339),
		"trace_id": traceID,
	})
}
`

// 带示例版 main.go
const HTTPExampleMain = `package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"{{.ModulePath}}/internal/router"
	"github.com/whosafe/uf/uconfig"
	"github.com/whosafe/uf/ucontext"
	"github.com/whosafe/uf/ulogger"
	"github.com/whosafe/uf/uprotocol/uhttp"
)

func main() {
	// 加载配置文件
	if err := uconfig.Load("config/config.yaml"); err != nil {
		panic("加载配置文件失败: " + err.Error())
	}

	// 初始化雪花算法
	if err := ucontext.InitSnowflake(1); err != nil {
		ulogger.Error("初始化雪花算法失败", "error", err)
		os.Exit(1)
	}

	// 创建并配置服务器
	server := uhttp.New()

	// 注册路由
	router.New(server)

	// 启动服务器
	go func() {
		if err := server.Start(); err != nil {
			ulogger.Error("服务器启动失败", "error", err)
		}
	}()

	// 优雅关闭
	gracefulShutdown(server)
}

// gracefulShutdown 优雅关闭服务器
func gracefulShutdown(server *uhttp.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ulogger.Info("正在关闭服务器...")

	ctx := ucontext.New()
	timeoutCtx, cancel := ctx.WithTimeout(5 * time.Second)
	defer cancel()

	if err := server.Stop(timeoutCtx.Context()); err != nil {
		ulogger.Error("服务器关闭失败", "error", err)
	}

	ulogger.Info("服务器已关闭")
}
`
