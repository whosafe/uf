# UF HTTP Server

高性能、功能完整的 HTTP 服务器框架,完全实现 `unet.Server` 接口。

## ✨ 核心特性

- 🚀 **高性能**: sync.Pool 对象池,零反射数据绑定
- 🔌 **协议无关**: 完全实现 unet.Server 接口
- 🎯 **链路追踪**: 所有日志自动包含 Trace ID
- ⚙️ **无感配置**: init 自动注册,uconfig.Load() 即可
- 📝 **双日志系统**: 访问日志和错误日志分离
- 🔒 **Session 管理**: 支持内存和 Redis 存储
- 📁 **静态文件**: 完整的文件服务支持
- 🍪 **Cookie 操作**: 丰富的 Cookie 辅助方法
- 📤 **文件上传**: 大小限制和扩展名验证

## 📦 安装

```bash
go get github.com/whosafe/uf/uprotocol/uhttp
```

## 🚀 快速开始

### 基础示例

```go
package main

import (
    "github.com/whosafe/uf/uconfig"
    "github.com/whosafe/uf/ucontext"
    "github.com/whosafe/uf/uprotocol/uhttp"
    "github.com/whosafe/uf/uprotocol/unet"
)

func main() {
    // 加载配置 (自动注册)
    uconfig.Load("config.yaml")
    
    // 创建服务器
    server := uhttp.New()
    
    // 自动应用默认中间件 (Trace, Logger, Recovery)
    // 根据配置文件中的 middleware 配置决定是否启用
    uhttp.ApplyDefaultMiddlewares(server)
    
    // 或者手动注册中间件
    // server.Use(uhttp.MiddlewareTrace())
    // server.Use(uhttp.MiddlewareLogger())
    // server.Use(uhttp.MiddlewareRecovery())
    
    // 注册路由
    server.GET("/", func(ctx *ucontext.Context, req unet.Request) error {
        return req.Response().JSON(200, map[string]string{
            "message": "Hello, World!",
        })
    })
    
    // 启动服务器
    server.Start(":8080")
}
```

### 配置文件示例

```yaml
server:
  name: "my-api"
  protocol: "http"
  address: ":8080"
  
  # 超时配置
  read_timeout: 30s
  write_timeout: 30s
  idle_timeout: 120s
  
  # 静态文件服务
  static:
    enabled: true
    root: "./public"
    prefix: "/static"
    index: ["index.html"]
    browse: false
  
  # Cookie 配置
  cookie:
    domain: ""
    path: "/"
    max_age: 86400
    secure: false
    http_only: true
    same_site: "lax"
  
  # Session 配置
  session:
    enabled: true
    provider: "memory"
    cookie_name: "session_id"
    max_age: 3600
  
  # 中间件配置
  middleware:
    # 核心中间件 (默认启用)
    enable_trace: true     # 追踪中间件
    enable_logger: true    # 日志中间件
    enable_recovery: true  # 恢复中间件
    
    # CORS 中间件 (默认禁用)
    enable_cors: false
    cors:
      allow_origins: "*"
      allow_methods: "GET,POST,PUT,DELETE,PATCH,HEAD,OPTIONS"
      allow_headers: "*"
      allow_credentials: false
      expose_headers: ""
      max_age: 3600
    
    # CSRF 中间件 (默认禁用)
    enable_csrf: false
    csrf:
      token_length: 32
      cookie_name: "csrf_token"
      header_name: "X-CSRF-Token"
      form_field_name: "csrf_token"
      cookie_max_age: 3600
    
    # 超时中间件 (默认禁用)
    enable_timeout: false
    timeout: "30s"
    
    # 限流中间件 (默认禁用)
    enable_rate_limit: false
    rate_limit:
      max_requests: 100
      window: "1m"
  
  # 访问日志
  access_log:
    enabled: true
    level: "info"
    format: "json"
    output: "stdout"
  
  # 错误日志
  error_log:
    enabled: true
    level: "error"
    format: "json"
    output: "stderr"
```

## 📖 功能详解

### 路由系统

```go
// 基础路由
server.GET("/users", getUsers)
server.POST("/users", createUser)
server.PUT("/users/:id", updateUser)
server.DELETE("/users/:id", deleteUser)

// 路径参数
server.GET("/users/:id", func(ctx *ucontext.Context, req unet.Request) error {
    httpReq := req.(*uhttp.Request)
    id := httpReq.Param("id")
    return req.Response().JSON(200, map[string]string{"id": id})
})

// 路由组
api := server.Group("/api")
{
    api.GET("/health", healthCheck)
    
    v1 := api.Group("/v1")
    {
        v1.GET("/users", getUsers)
        v1.POST("/users", createUser)
    }
}
```

### 中间件

```go
// 方式1: 自动应用默认中间件 (推荐)
uhttp.ApplyDefaultMiddlewares(server)

// 方式2: 手动注册中间件
server.Use(uhttp.MiddlewareTrace())    // 链路追踪
server.Use(uhttp.MiddlewareLogger())   // 请求日志
server.Use(uhttp.MiddlewareRecovery()) // 异常恢复
server.Use(uhttp.MiddlewareCORS())     // 跨域支持
server.Use(uhttp.MiddlewareCSRF())     // CSRF 保护
server.Use(uhttp.MiddlewareTimeout(30 * time.Second)) // 超时控制

// 限流中间件
server.Use(uhttp.MiddlewareRateLimit())  // 默认配置
server.Use(uhttp.MiddlewareRateLimitByIP(100, time.Minute))  // 基于 IP
server.Use(uhttp.MiddlewareRateLimitByPath(50, time.Minute)) // 基于路径

// 路由级中间件
server.GET("/admin", adminHandler, authMiddleware)
```

### CSRF 保护

```go
// 方式1: 通过配置文件启用 (推荐)
// 在 config.yaml 中设置 middleware.enable_csrf: true

// 方式2: 手动启用
server.Use(uhttp.MiddlewareCSRF())

// 方式3: 自定义配置
server.Use(uhttp.MiddlewareCSRFWithConfig(uhttp.CSRFConfig{
    TokenLength:   32,
    CookieName:    "csrf_token",
    HeaderName:    "X-CSRF-Token",
    FormFieldName: "csrf_token",
    CookieMaxAge:  3600,
}))

// 在 GET 请求中获取 Token
server.GET("/form", func(ctx *ucontext.Context, req unet.Request) error {
    httpReq := req.(*uhttp.Request)
    csrfToken := httpReq.GetCSRFToken()
    
    // 在 HTML 表单中使用
    html := `<form method="POST" action="/submit">
        <input type="hidden" name="csrf_token" value="` + csrfToken + `">
        <!-- 其他表单字段 -->
    </form>`
    
    return req.Response().HTML(200, html)
})

// POST 请求自动验证 CSRF Token
server.POST("/submit", func(ctx *ucontext.Context, req unet.Request) error {
    // 如果到达这里,说明 CSRF 验证已通过
    return req.Response().JSON(200, map[string]string{"status": "ok"})
})

// AJAX 请求中使用
// fetch('/submit', {
//     method: 'POST',
//     headers: { 'X-CSRF-Token': csrfToken },
//     body: JSON.stringify(data)
// })
```

### 路由级中间件

```go
// 路由级中间件
server.GET("/admin", adminHandler, authMiddleware)
```

### 请求处理

```go
func handler(ctx *ucontext.Context, req unet.Request) error {
    httpReq := req.(*uhttp.Request)
    
    // 获取路径参数
    id := httpReq.Param("id")
    
    // 获取查询参数
    name := httpReq.Query("name")
    
    // 绑定 JSON
    var data struct {
        Name string `json:"name"`
        Age  int    `json:"age"`
    }
    if err := httpReq.BindJSON(&data); err != nil {
        return err
    }
    
    // 绑定 Form
    if err := httpReq.BindForm(&data); err != nil {
        return err
    }
    
    return req.Response().JSON(200, data)
}
```

### 响应处理

```go
resp := req.Response()

// JSON 响应
resp.JSON(200, map[string]string{"status": "ok"})

// 字符串响应
resp.String(200, "Hello, World!")

// 字节响应
resp.Bytes(200, []byte("data"))

// 设置状态码
resp.Status(404)

// 设置 Header
resp.Header("Content-Type", "application/json")
```

### 静态文件服务

```go
// 方式1: 配置文件 (推荐)
// 在 config.yaml 中配置 static 即可

// 方式2: 代码注册
server.Static("/static", "./public")

// 方式3: 自定义配置
server.StaticWithConfig(&uhttp.StaticConfig{
    Root:   "./public",
    Prefix: "/static",
    Index:  []string{"index.html"},
    Browse: false,
})

// 单文件服务
server.File("/favicon.ico", "./public/favicon.ico")
```

### Cookie 操作

```go
// 设置 Cookie
resp.SetCookieValue("user", "alice", 3600)

// 读取 Cookie
user, _ := httpReq.GetCookie("user")

// 删除 Cookie
resp.DeleteCookie("user")

// 安全 Cookie
resp.SetSecureCookie("token", "xxx", 3600, "example.com")

// 会话 Cookie
resp.SetSessionCookie(name string, id string, path string, domain string, age int, secure bool, only bool, site http.SameSite)
```

### Session 管理

```go
// 方式1: 配置文件 (推荐)
// 在 config.yaml 中配置 session 即可

// 方式2: 手动创建
sessionMgr := uhttp.NewSessionManager(
    uhttp.NewMemoryStore(),
    "session_id",
    3600,
)

// 在处理器中使用
func loginHandler(ctx *ucontext.Context, req unet.Request) error {
    httpReq := req.(*uhttp.Request)
    httpResp := req.Response().(*uhttp.Response)
    
    // 获取 Session 管理器
    sessionMgr := httpReq.Server().SessionManager()
    
    // 启动会话
    session, _ := sessionMgr.Start(httpReq, httpResp)
    
    // 设置数据
    session.Set("user_id", 123)
    session.Set("username", "alice")
    session.Save()
    
    return httpResp.JSON(200, map[string]string{"status": "ok"})
}

func profileHandler(ctx *ucontext.Context, req unet.Request) error {
    httpReq := req.(*uhttp.Request)
    httpResp := req.Response().(*uhttp.Response)
    
    sessionMgr := httpReq.Server().SessionManager()
    session, _ := sessionMgr.Start(httpReq, httpResp)
    
    // 读取数据
    userID, _ := session.Get("user_id")
    username, _ := session.Get("username")
    
    return httpResp.JSON(200, map[string]any{
        "user_id": userID,
        "username": username,
    })
}

// Redis 存储
import "github.com/redis/go-redis/v9"

redisClient := redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
})

store := uhttp.NewRedisStore(redisClient, "session:", 30*time.Minute)
sessionMgr := uhttp.NewSessionManager(store, "session_id", 3600)
```

### 文件上传

```go
func uploadHandler(ctx *ucontext.Context, req unet.Request) error {
    httpReq := req.(*uhttp.Request)
    
    // 单文件上传
    file, _ := httpReq.FormFile("file")
    httpReq.SaveUploadedFile(file, "./uploads/" + file.Filename)
    
    // 多文件上传
    paths, _ := httpReq.SaveUploadedFiles("files", "./uploads")
    
    // 配置化上传
    path, _ := httpReq.SaveUploadedFileWithConfig(file, &uhttp.FileUploadConfig{
        MaxSize:     10 << 20, // 10MB
        AllowedExts: []string{".jpg", ".png", ".gif"},
        UploadDir:   "./uploads",
    })
    
    return req.Response().JSON(200, map[string]any{
        "path": path,
    })
}
```

## 🔧 配置说明

### 服务器配置

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| name | string | "uhttp-server" | 服务名称 |
| protocol | string | "http" | 协议类型 |
| address | string | ":8080" | 监听地址 |
| read_timeout | duration | 30s | 读取超时 |
| write_timeout | duration | 30s | 写入超时 |
| idle_timeout | duration | 120s | 空闲超时 |
| max_header_bytes | int | 1MB | 最大请求头 |
| max_body_bytes | int | 10MB | 最大请求体 |
| keep_alive | bool | true | 启用 Keep-Alive |

### 静态文件配置

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| enabled | bool | false | 是否启用 |
| root | string | - | 静态文件根目录 |
| prefix | string | - | URL 前缀 |
| index | []string | ["index.html"] | 索引文件 |
| browse | bool | false | 允许目录浏览 |

### Cookie 配置

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| domain | string | "" | Cookie 域 |
| path | string | "/" | Cookie 路径 |
| max_age | int | 86400 | 最大存活时间(秒) |
| secure | bool | false | 仅 HTTPS |
| http_only | bool | true | 禁止 JS 访问 |
| same_site | string | "lax" | SameSite 策略 |

### Session 配置

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| enabled | bool | false | 是否启用 |
| provider | string | "memory" | 存储类型 |
| cookie_name | string | "session_id" | Cookie 名称 |
| max_age | int | 3600 | 过期时间(秒) |

### 日志配置

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| enabled | bool | true | 是否启用 |
| level | string | "info" | 日志级别 |
| format | string | "text" | 格式 (json/text) |
| output | string | "stdout" | 输出 (stdout/stderr/file) |
| file_path | string | - | 文件路径 |
| max_size | int | 100 | 最大文件大小(MB) |
| max_backups | int | 10 | 最大备份数量 |
| max_age | int | 30 | 最大保留天数 |
| compress | bool | true | 是否压缩 |

## 🎯 性能优化

### sync.Pool 对象池

Request 和 Response 对象使用 sync.Pool 复用,减少 GC 压力:

```go
// 自动管理,无需手动操作
// 每个请求结束后自动释放回对象池
```

### 零反射数据绑定

使用 `ubind` 进行数据绑定,避免反射开销:

```go
// 高性能的数据绑定
httpReq.BindJSON(&data)
httpReq.BindForm(&data)
```

### 链路追踪

所有日志自动包含 Trace ID,便于问题追踪:

```json
{
  "level": "info",
  "msg": "HTTP Request",
  "trace_id": "1234567890abcdef",
  "span_id": "abcdef1234567890",
  "method": "GET",
  "path": "/api/users",
  "status": 200,
  "duration_ms": 15
}
```

## 📚 API 文档

### Server

- `New() *Server` - 创建服务器
- `NewWithConfig(cfg *Config) *Server` - 使用配置创建
- `Start(addr string) error` - 启动服务器
- `Stop(ctx context.Context) error` - 停止服务器
- `Use(middlewares ...unet.MiddlewareFunc)` - 注册全局中间件
- `GET/POST/PUT/DELETE/PATCH/HEAD/OPTIONS(path string, handler unet.HandlerFunc, middlewares ...unet.MiddlewareFunc)` - 注册路由
- `Group(prefix string) *Group` - 创建路由组
- `Static(prefix, root string)` - 注册静态文件服务
- `File(path, filepath string)` - 注册单文件服务

### Request

- `Method() string` - 获取请求方法
- `Path() string` - 获取请求路径
- `Param(key string) string` - 获取路径参数
- `Query(key string) string` - 获取查询参数
- `Header(key string) string` - 获取请求头
- `Cookie(name string) (*http.Cookie, error)` - 获取 Cookie
- `GetCookie(name string) (string, error)` - 获取 Cookie 值
- `BindJSON(v any) error` - 绑定 JSON
- `BindForm(v any) error` - 绑定 Form
- `FormFile(name string) (*multipart.FileHeader, error)` - 获取上传文件
- `SaveUploadedFile(file *multipart.FileHeader, dst string) error` - 保存文件

### Response

- `Status(code int)` - 设置状态码
- `Header(key, value string)` - 设置响应头
- `JSON(code int, v any) error` - JSON 响应
- `String(code int, s string) error` - 字符串响应
- `Bytes(code int, b []byte) error` - 字节响应
- `SetCookie(cookie *http.Cookie)` - 设置 Cookie
- `SetCookieValue(name, value string, maxAge int)` - 快速设置 Cookie
- `DeleteCookie(name string)` - 删除 Cookie

## 📝 示例项目

查看 `example/uhttp` 目录获取完整示例。

## 🤝 贡献

欢迎提交 Issue 和 Pull Request!

## 📄 许可证

MIT License
