package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/whosafe/uf/ucontext"
	"github.com/whosafe/uf/ulogger"
	"github.com/whosafe/uf/uprotocol/ubind"
	"github.com/whosafe/uf/uprotocol/uhttp"
	"github.com/whosafe/uf/uprotocol/unet"
)

// User 用户模型
type User struct {
	ID   int
	Name string
	Age  int
}

// Bind 实现 ubind.Binder 接口
func (u *User) Bind(key string, value *ubind.Value) error {
	switch key {
	case "id":
		u.ID = value.Int()
	case "name":
		u.Name = value.Str()
	case "age":
		u.Age = value.Int()
	}
	return nil
}

func main() {
	// 初始化日志
	logger, _ := ulogger.New(&ulogger.Config{
		Level:  slog.LevelDebug,
		Stdout: true,
	})
	ulogger.SetDefault(logger)

	// 初始化雪花算法
	if err := ucontext.InitSnowflake(1); err != nil {
		ulogger.Error("初始化雪花算法失败", "error", err)
		return
	}

	// 创建 HTTP 服务器
	server := uhttp.New()

	// 注册全局中间件
	server.Use(uhttp.MiddlewareTrace())    // 链路追踪
	server.Use(uhttp.MiddlewareLogger())   // 请求日志
	server.Use(uhttp.MiddlewareRecovery()) // 异常恢复
	server.Use(uhttp.MiddlewareCORS())     // 跨域支持

	// 根路径
	server.GET("/", handleIndex)

	// API 路由
	api := server.Group("/api")
	{
		// 用户相关
		api.GET("/users", listUsers)
		api.GET("/users/:id", getUser)
		api.POST("/users", createUser)
		api.PUT("/users/:id", updateUser)
		api.DELETE("/users/:id", deleteUser)

		// 健康检查
		api.GET("/health", healthCheck)
	}

	// 启动服务器
	addr := ":8080"
	ulogger.Info("HTTP 服务器启动", "addr", addr)

	// 优雅关闭
	go func() {
		if err := server.Start(addr); err != nil {
			ulogger.Error("服务器启动失败", "error", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ulogger.Info("正在关闭服务器...")

	// 优雅关闭,最多等待 5 秒
	ctx := ucontext.New()
	timeoutCtx, cancel := ctx.WithTimeout(5 * time.Second)
	defer cancel()

	if err := server.Stop(timeoutCtx.Context()); err != nil {
		ulogger.Error("服务器关闭失败", "error", err)
	}

	ulogger.Info("服务器已关闭")
}

// handleIndex 首页
func handleIndex(ctx *ucontext.Context, req unet.Request) error {
	resp := req.Response()
	return resp.JSON(200, map[string]any{
		"message": "欢迎使用 UF HTTP Server",
		"version": "1.0.0",
		"docs":    "/api/health",
	})
}

// listUsers 获取用户列表
func listUsers(ctx *ucontext.Context, req unet.Request) error {
	users := []User{
		{ID: 1, Name: "Alice", Age: 25},
		{ID: 2, Name: "Bob", Age: 30},
		{ID: 3, Name: "Charlie", Age: 35},
	}

	ulogger.InfoCtx(ctx.Context(), "获取用户列表", "count", len(users))

	resp := req.Response()
	return resp.JSON(200, map[string]any{
		"users": users,
		"total": len(users),
	})
}

// getUser 获取单个用户
func getUser(ctx *ucontext.Context, req unet.Request) error {
	httpReq := req.(*uhttp.Request)
	id := httpReq.Param("id")

	ulogger.InfoCtx(ctx.Context(), "获取用户", "id", id)

	// 模拟查询
	user := User{ID: 1, Name: "Alice", Age: 25}

	resp := req.Response()
	return resp.JSON(200, user)
}

// createUser 创建用户
func createUser(ctx *ucontext.Context, req unet.Request) error {
	httpReq := req.(*uhttp.Request)

	var user User
	if err := httpReq.BindJSON(&user); err != nil {
		ulogger.WarnCtx(ctx.Context(), "绑定用户数据失败", "error", err)
		resp := req.Response()
		return resp.JSON(400, map[string]string{
			"error": "无效的请求数据",
		})
	}

	// 生成 ID
	user.ID = 100

	ulogger.InfoCtx(ctx.Context(), "创建用户", "user_id", user.ID, "name", user.Name)

	resp := req.Response()
	return resp.JSON(201, user)
}

// updateUser 更新用户
func updateUser(ctx *ucontext.Context, req unet.Request) error {
	httpReq := req.(*uhttp.Request)
	id := httpReq.Param("id")

	var user User
	if err := httpReq.BindJSON(&user); err != nil {
		resp := req.Response()
		return resp.JSON(400, map[string]string{
			"error": "无效的请求数据",
		})
	}

	ulogger.InfoCtx(ctx.Context(), "更新用户", "id", id, "name", user.Name)

	resp := req.Response()
	return resp.JSON(200, user)
}

// deleteUser 删除用户
func deleteUser(ctx *ucontext.Context, req unet.Request) error {
	httpReq := req.(*uhttp.Request)
	id := httpReq.Param("id")

	ulogger.InfoCtx(ctx.Context(), "删除用户", "id", id)

	resp := req.Response()
	return resp.JSON(200, map[string]string{
		"message": fmt.Sprintf("用户 %s 已删除", id),
	})
}

// healthCheck 健康检查
func healthCheck(ctx *ucontext.Context, req unet.Request) error {
	tc := ucontext.FromContext(ctx.Context())
	traceID := ""
	if tc != nil {
		traceID = tc.TraceID
	}

	resp := req.Response()
	return resp.JSON(200, map[string]any{
		"status":   "ok",
		"time":     time.Now().Format(time.RFC3339),
		"trace_id": traceID,
	})
}
