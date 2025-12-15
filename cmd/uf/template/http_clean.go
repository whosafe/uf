package template

// HTTP 纯净版 main.go
const HTTPCleanMain = `package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/whosafe/uf/uconfig"
	"github.com/whosafe/uf/ucontext"
	"github.com/whosafe/uf/ulogger"
	"github.com/whosafe/uf/uprotocol/uhttp"
	"github.com/whosafe/uf/uprotocol/unet"
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
	server := setupServer()

	// 启动服务器
	go func() {
		if err := server.Start(); err != nil {
			ulogger.Error("服务器启动失败", "error", err)
		}
	}()

	// 优雅关闭
	gracefulShutdown(server)
}

// setupServer 创建并配置服务器
func setupServer() *uhttp.Server {
	server := uhttp.New()

	// 注册路由
	registerRoutes(server)

	return server
}

// registerRoutes 注册路由
func registerRoutes(server *uhttp.Server) {
	server.GET("/", handleIndex)
	server.GET("/health", handleHealth)
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

// handleIndex 首页
func handleIndex(ctx *ucontext.Context, req unet.Request) error {
	return req.Response().JSON(200, map[string]string{
		"message": "欢迎使用 {{.ProjectName}}",
		"version": "1.0.0",
	})
}

// handleHealth 健康检查
func handleHealth(ctx *ucontext.Context, req unet.Request) error {
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

// HTTP 纯净版 config.yaml
const HTTPCleanConfig = `server:
  name: "{{.ProjectName}}"
  protocol: "http"
  address: ":8080"
  
  read_timeout: 30s
  write_timeout: 30s
  idle_timeout: 120s
  
  max_body_bytes: 10485760  # 10MB
  
  middleware:
    enable_trace: true
    enable_logger: true
    enable_recovery: true
`

// go.mod 模板
const GoModTemplate = `module {{.ModulePath}}

go 1.25

require github.com/whosafe/uf v0.0.1
`

// .gitignore 模板
const GitignoreTemplate = `# Binaries
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary
*.test

# Output
*.out

# Go workspace file
go.work

# IDE
.idea/
.vscode/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db

# Logs
logs/
*.log

# Build
dist/
build/
`

// README 模板
const ReadmeTemplate = `# {{.ProjectName}}

使用 UF 框架构建的 {{.Protocol}} 应用。

## 快速开始

### 安装依赖

` + "```bash" + `
go mod tidy
` + "```" + `

### 运行项目

` + "```bash" + `
go run main.go
` + "```" + `

服务器将在 http://localhost:8080 启动（可在 config/config.yaml 中修改）。

### 测试 API

` + "```bash" + `
# 访问首页
curl http://localhost:8080/

# 健康检查
curl http://localhost:8080/health
` + "```" + `

## 项目结构

` + "```" + `
.
├── main.go          # 主入口
├── config/          # 配置文件
│   └── config.yaml
├── internal/        # 私有代码
│   ├── consts/      # 常量定义
│   ├── handler/     # HTTP 处理器
│   ├── model/       # 数据模型
│   ├── dao/         # 数据访问层
│   ├── router/      # 路由定义
│   ├── service/     # 业务服务层
│   ├── logic/       # 业务逻辑层
│   └── middleware/  # 中间件
├── utility/         # 工具函数
├── hack/            # 脚本和工具
├── go.mod
└── README.md
` + "```" + `

## 代码结构说明

**main 函数**：
1. 加载配置文件
2. 初始化雪花算法
3. 创建并配置服务器
4. 启动服务器
5. 优雅关闭

**辅助函数**：
- **setupServer()**: 创建服务器并注册路由
- **registerRoutes()**: 路由注册
- **gracefulShutdown()**: 优雅关闭处理

## 配置文件

配置文件位于 ` + "`config/config.yaml`" + `，可以配置：
- 服务器地址和端口
- 超时设置
- 中间件开关
- 日志配置等

## 文档

- [UF 框架文档](https://github.com/yourusername/uf)

## License

MIT
`
