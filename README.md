# UF 框架使用文档

UF (Unified Framework) 是一个高性能、模块化的 Go 语言 Web 框架，专注于简洁性、性能和可维护性。

## ✨ 核心特性

- 🚀 **高性能**: 零反射设计，sync.Pool 对象池，性能接近原生代码
- 🔌 **协议无关**: 基于 unet 抽象层，支持 HTTP、TCP、QUIC 等多种协议
- 📦 **零依赖配置**: 自研 YAML 解析器，比标准库快 4-5 倍
- 📝 **完善日志**: 基于 slog 的高性能日志系统，支持文件轮转和压缩
- 🔗 **链路追踪**: 内置分布式追踪，雪花算法 ID 生成
- ✅ **数据验证**: 零反射的结构化验证器，支持国际化
- 🎯 **类型安全**: 编译期检查，避免运行时错误
- ⚙️ **无感配置**: init 自动注册，一行代码加载配置

## 📦 安装

```bash
go get github.com/whosafe/uf
```

**环境要求**: Go 1.25+

## 🚀 快速开始

### 最简示例

创建一个最简单的 HTTP 服务器：

```go
package main

import (
    "github.com/whosafe/uf/ucontext"
    "github.com/whosafe/uf/uprotocol/uhttp"
    "github.com/whosafe/uf/uprotocol/unet"
)

func main() {
    // 创建服务器
    server := uhttp.New()
    
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

运行程序后访问 `http://localhost:8080`，即可看到 JSON 响应。

### 带配置文件的示例

创建配置文件 `config.yaml`：

```yaml
server:
  name: "my-api"
  protocol: "http"
  address: ":8080"
  
  # 中间件配置
  middleware:
    enable_trace: true
    enable_logger: true
    enable_recovery: true
```

创建应用程序：

```go
package main

import (
    "github.com/whosafe/uf/uconfig"
    "github.com/whosafe/uf/ucontext"
    "github.com/whosafe/uf/uprotocol/uhttp"
    "github.com/whosafe/uf/uprotocol/unet"
)

func main() {
    // 加载配置文件
    uconfig.Load("config.yaml")
    
    // 创建服务器
    server := uhttp.New()
    
    // 自动应用默认中间件（根据配置文件）
    uhttp.ApplyDefaultMiddlewares(server)
    
    // 注册路由
    server.GET("/", func(ctx *ucontext.Context, req unet.Request) error {
        return req.Response().JSON(200, map[string]string{
            "message": "Server with config",
        })
    })
    
    // 启动服务器（地址来自配置文件）
    server.Start("")
}
```

## 📚 核心模块详解

### 1. unet - 网络层抽象

`unet` 提供协议无关的网络层抽象，使业务逻辑与具体协议解耦。

#### 核心设计理念

**两参数设计** - 所有处理器函数使用统一签名：

```go
type HandlerFunc func(ctx *ucontext.Context, req Request) error
```

- `*ucontext.Context` - 链路追踪上下文
- `Request` - 请求对象（包含请求数据和响应能力）

#### 核心接口

```go
// 请求接口
type Request interface {
    Protocol() Protocol
    RemoteAddr() net.Addr
    LocalAddr() net.Addr
    Get(key string) (any, bool)
    Set(key string, value any)
    Bind(obj ubind.Binder) error
    Response() Response
}

// 响应接口
type Response interface {
    JSON(code int, data any) error
    String(code int, text string) error
    Bytes(code int, data []byte) error
}

// 服务器接口
type Server interface {
    Start(addr string) error
    Stop(ctx context.Context) error
    Use(middleware ...MiddlewareFunc)
    Handle(pattern string, handler HandlerFunc)
}
```

#### 协议无关的优势

同一个处理器可以在不同协议中使用：

```go
func CreateUser(ctx *ucontext.Context, req unet.Request) error {
    var user User
    if err := req.Bind(&user); err != nil {
        return err
    }
    
    // 业务逻辑
    saveUser(ctx.Context(), &user)
    
    return req.Response().JSON(200, user)
}

// HTTP
httpServer.POST("/users", CreateUser)

// TCP (未来支持)
tcpServer.Handle(MSG_CREATE_USER, CreateUser)

// QUIC (未来支持)
quicServer.Handle("/users", CreateUser)
```

**详细文档**: [uprotocol/unet/README.md](uprotocol/unet/README.md)

---

### 2. uhttp - HTTP 服务器

`uhttp` 是一个高性能、功能完整的 HTTP 服务器框架，完全实现 `unet.Server` 接口。

#### 核心特性

- 🚀 高性能：sync.Pool 对象池，零反射数据绑定
- 🎯 链路追踪：所有日志自动包含 Trace ID
- ⚙️ 无感配置：init 自动注册，uconfig.Load() 即可
- 📝 双日志系统：访问日志和错误日志分离
- 🔒 Session 管理：支持内存和 Redis 存储
- 📁 静态文件：完整的文件服务支持
- 🍪 Cookie 操作：丰富的 Cookie 辅助方法

#### 路由系统

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

#### 中间件系统

```go
// 方式1: 自动应用默认中间件（推荐）
uhttp.ApplyDefaultMiddlewares(server)

// 方式2: 手动注册中间件
server.Use(uhttp.MiddlewareTrace())    // 链路追踪
server.Use(uhttp.MiddlewareLogger())   // 请求日志
server.Use(uhttp.MiddlewareRecovery()) // 异常恢复
server.Use(uhttp.MiddlewareCORS())     // 跨域支持

// 限流中间件
server.Use(uhttp.MiddlewareRateLimit())  // 默认配置
server.Use(uhttp.MiddlewareRateLimitByIP(100, time.Minute))  // 基于 IP

// 路由级中间件
server.GET("/admin", adminHandler, authMiddleware)
```

#### 请求处理

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
    
    return req.Response().JSON(200, data)
}
```

#### Session 管理

```go
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
```

#### 文件上传

```go
func uploadHandler(ctx *ucontext.Context, req unet.Request) error {
    httpReq := req.(*uhttp.Request)
    
    // 获取上传文件
    file, _ := httpReq.FormFile("file")
    
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

#### 静态文件服务

```go
// 方式1: 配置文件（推荐）
// 在 config.yaml 中配置 static 即可

// 方式2: 代码注册
server.Static("/static", "./public")

// 单文件服务
server.File("/favicon.ico", "./public/favicon.ico")
```

**详细文档**: [uprotocol/uhttp/README.md](uprotocol/uhttp/README.md)

---

### 3. uconfig - 配置管理

`uconfig` 是一个零依赖、高性能的配置加载库，采用自研 YAML 解析器。

#### 核心特性

- 🚀 极致性能：比 `yaml.v3` 快 4-5 倍，内存占用减少 40%
- 📦 零依赖：无需任何第三方库
- 🎯 类型安全：强类型 API，编译期检查
- 🔧 灵活解析：支持自定义解析逻辑
- 🔌 回调机制：支持未知配置项的被动回调

#### 基本使用

**1. 定义配置结构**

```go
type ServerConfig struct {
    Host string
    Port []int
}

func (s *ServerConfig) UnmarshalYAML(key string, value *uconfig.Node) error {
    switch key {
    case "host":
        s.Host = value.String()
    case "port":
        s.Port = make([]int, 0)
        return value.Iter(func(i int, v *uconfig.Node) error {
            s.Port = append(s.Port, uconv.ToIntDef(v, 0))
            return nil
        })
    }
    return nil
}
```

**2. 注册和加载**

```go
var srvCfg ServerConfig

// 注册配置解析器
uconfig.Register("server", srvCfg.UnmarshalYAML)

// 加载配置文件
if err := uconfig.Load("config.yaml"); err != nil {
    log.Fatal(err)
}

fmt.Printf("Server: %+v\n", srvCfg)
```

#### 嵌套结构解析

```go
type DatabaseConfig struct {
    DSN     string
    MaxOpen int
    Logger  LoggerConfig
}

func (d *DatabaseConfig) UnmarshalYAML(key string, value *uconfig.Node) error {
    switch key {
    case "dsn":
        d.DSN = value.String()
    case "max_open":
        d.MaxOpen = uconv.ToIntDef(value, 0)
    case "logger":
        // 递归解析嵌套结构
        return value.Decode(&d.Logger)
    }
    return nil
}
```

**详细文档**: [uconfig/README.md](uconfig/README.md)

---

### 4. ulogger - 日志系统

`ulogger` 是一个功能完善的日志库，基于标准库 `log/slog`，支持文件轮转、压缩、多输出等企业级特性。

#### 核心特性

- 🎯 基于 slog：完全兼容 Go 1.21+ 的标准日志接口
- 🔄 文件轮转：支持按大小和时间自动轮转
- 📦 备份管理：自动清理过期备份
- 🗜️ 压缩支持：可选的 gzip 压缩（0-9 级别）
- 📤 多输出：同时输出到文件和终端
- 🎨 多格式支持：支持文本、JSON 和自定义格式

#### 基本使用

```go
import (
    "log/slog"
    "github.com/whosafe/uf/ulogger"
)

// 使用默认配置
logger, _ := ulogger.New(ulogger.DefaultConfig())
defer logger.Close()

logger.Info("Hello, ulogger!")
logger.Warn("This is a warning", "key", "value")
```

#### 输出到文件

```go
config := &ulogger.Config{
    Path:   "./logs",
    File:   "app.log",
    Level:  slog.LevelInfo,
    Stdout: true, // 同时输出到终端
}

logger, _ := ulogger.New(config)
defer logger.Close()

logger.Info("日志已写入文件")
```

#### JSON 格式日志

```go
config := &ulogger.Config{
    Path:   "./logs",
    File:   "app.json",
    Format: "json", // 使用 JSON 格式
    Level:  slog.LevelInfo,
}

logger, _ := ulogger.New(config)
logger.Info("API 请求", "method", "GET", "status", 200)
// 输出: {"time":"2025-12-16T15:50:00+08:00","level":"INFO","msg":"API 请求","attrs":{"method":"GET","status":200}}
```

#### 文件轮转

```go
config := &ulogger.Config{
    Path:                 "./logs",
    File:                 "app.log",
    RotateSize:           100 * 1024 * 1024, // 100MB
    RotateBackupLimit:    10,
    RotateBackupExpire:   7 * 24 * 3600, // 7 天
    RotateBackupCompress: 6,
}

logger, _ := ulogger.New(config)
```

#### 全局 Logger

```go
// 设置全局 Logger
logger, _ := ulogger.New(config)
ulogger.SetDefault(logger)

// 使用全局函数
ulogger.Info("global message")
ulogger.Error("error occurred", "error", err)

// 标准库 slog 也会使用我们的 Logger
slog.Info("this uses our logger too")
```

**详细文档**: [ulogger/README.md](ulogger/README.md)

---

### 5. ucontext - 链路追踪

`ucontext` 是一个轻量级的分布式链路追踪包，支持雪花算法 ID 生成、采样控制、HTTP 传播和 Logger 集成。

#### 核心特性

- 🔢 雪花算法：分布式唯一 ID 生成
- 🔗 链路追踪：Trace ID、Span ID、Parent Span ID
- 📊 采样控制：可配置的采样率，支持强制采样
- 🌐 HTTP 传播：跨服务传递追踪信息
- 📝 Logger 集成：自动注入追踪信息到日志

#### 基本使用

```go
import (
    "context"
    "github.com/whosafe/uf/ucontext"
)

func main() {
    // 初始化雪花算法（worker ID）
    ucontext.InitSnowflake(1)

    // 创建追踪上下文
    ctx := ucontext.NewContext(context.Background())
    tc := ucontext.FromContext(ctx)

    println("Trace ID:", tc.TraceID)
    println("Span ID:", tc.SpanID)
}
```

#### Logger 集成

```go
ctx := ucontext.NewContext(context.Background())

// 使用全局 Logger 的 Context 方法
ulogger.InfoCtx(ctx, "处理请求", "key", "value")
ulogger.DebugCtx(ctx, "调试信息")

// 输出包含 trace_id 和 span_id:
// 2025-12-16 15:50:00 [INFO] main.go:10 处理请求 key=value trace_id=522314532622700544 span_id=522314532622700545
```

#### HTTP 传播

```go
// 服务端：提取追踪信息
func handler(w http.ResponseWriter, r *http.Request) {
    tc := ucontext.ExtractHTTPHeaders(r.Header)
    ctx := ucontext.WithContext(r.Context(), tc)
    
    logger := ucontext.LoggerFromContext(ctx)
    logger.Info("处理请求")
}

// 客户端：注入追踪信息
func callAPI(ctx context.Context, url string) {
    req, _ := http.NewRequest("GET", url, nil)
    tc := ucontext.FromContext(ctx)
    ucontext.InjectHTTPHeaders(req.Header, tc)
    
    resp, _ := http.DefaultClient.Do(req)
    // ...
}
```

#### 嵌套 Span

```go
func processRequest(ctx context.Context) {
    logger := ucontext.LoggerFromContext(ctx)
    logger.Info("开始处理请求")
    
    // 创建子 Span
    childCtx := ucontext.NewSpan(ctx)
    queryDatabase(childCtx)
    
    logger.Info("请求处理完成")
}

func queryDatabase(ctx context.Context) {
    logger := ucontext.LoggerFromContext(ctx)
    logger.Info("查询数据库")
    // 子 Span 的日志会包含父子关系
}
```

**详细文档**: [ucontext/README.md](ucontext/README.md)

---

### 6. ubind - 数据绑定

`ubind` 是一个零反射的数据绑定解决方案，用于协议无关的请求数据解析。

#### 核心特性

- ✅ 零反射：手动实现，高性能
- ✅ 类型安全：编译时检查
- ✅ 协议无关：支持 HTTP、TCP、QUIC 等多种协议
- ✅ 自动识别：自动识别 JSON/Form/Binary 格式
- ✅ 支持嵌套：完整支持嵌套对象和数组

#### 基本使用

```go
import "github.com/whosafe/uf/uprotocol/ubind"

type User struct {
    ID   int
    Name string
    Age  int
}

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
    jsonData := []byte(`{"id":1,"name":"Alice","age":25}`)
    
    val := ubind.Parse(jsonData)  // 自动识别 JSON
    var user User
    ubind.Bind(val, &user)
    
    // user.ID = 1, user.Name = "Alice", user.Age = 25
}
```

#### 嵌套对象

```go
type Address struct {
    City   string
    Street string
}

func (a *Address) Bind(key string, value *ubind.Value) error {
    switch key {
    case "city":
        a.City = value.Str()
    case "street":
        a.Street = value.Str()
    }
    return nil
}

type User struct {
    ID      int
    Name    string
    Address Address
}

func (u *User) Bind(key string, value *ubind.Value) error {
    switch key {
    case "id":
        u.ID = value.Int()
    case "name":
        u.Name = value.Str()
    case "address":
        if value.IsObject() {
            return ubind.Bind(value, &u.Address)
        }
    }
    return nil
}
```

**详细文档**: [uprotocol/ubind/README.md](uprotocol/ubind/README.md)

---

### 7. umarshal - JSON 序列化

`umarshal` 是一个高性能、零反射的 JSON 序列化库。

#### 核心特性

- ✅ 高性能：比标准库快 20-30%
- ✅ 安全：完整的字符串转义处理
- ✅ 对象池：复用 Writer 对象，减少 GC 压力
- ✅ 零反射：通过接口实现自定义序列化

#### 基本使用

```go
import "github.com/whosafe/uf/uprotocol/umarshal"

// 序列化基础类型
data, _ := umarshal.Marshal("hello")
// 输出: "hello"

data, _ = umarshal.Marshal(123)
// 输出: 123
```

#### 自定义序列化

```go
type User struct {
    ID   int
    Name string
    Age  int
}

// 实现 Marshaler 接口
func (u *User) MarshalJSON(w *umarshal.Writer) error {
    w.WriteObjectStart()
    w.WriteObjectField("id")
    w.WriteInt(u.ID)
    w.WriteComma()
    w.WriteObjectField("name")
    w.WriteString(u.Name)
    w.WriteComma()
    w.WriteObjectField("age")
    w.WriteInt(u.Age)
    w.WriteObjectEnd()
    return nil
}

// 使用
user := &User{ID: 1, Name: "Alice", Age: 25}
data, _ := umarshal.Marshal(user)
// 输出: {"id":1,"name":"Alice","age":25}
```

#### 使用 Writer（更高性能）

```go
// 从对象池获取 Writer
w := umarshal.AcquireWriter()
defer umarshal.ReleaseWriter(w)

// 手动构建 JSON
w.WriteObjectStart()
w.WriteObjectField("status")
w.WriteString("ok")
w.WriteComma()
w.WriteObjectField("code")
w.WriteInt(0)
w.WriteObjectEnd()

// 获取结果
result := w.Bytes()
```

**详细文档**: [uprotocol/umarshal/README.md](uprotocol/umarshal/README.md)

---

### 8. uvalidator - 数据验证

`uvalidator` 是一个基于结构化规则的高性能验证器，零反射。

#### 核心特性

- ✅ 结构化：每个规则都是独立的结构体
- ✅ 零反射：手动实现验证逻辑
- ✅ 高性能：接近手写代码
- ✅ 易扩展：添加新规则只需创建新文件
- ✅ 国际化：支持多语言错误消息

#### 基本使用

```go
import (
    "github.com/whosafe/uf/uvalidator"
    "github.com/whosafe/uf/uvalidator/rule"
)

type CreateUserRequest struct {
    Username string
    Email    string
    Age      int
}

func (r *CreateUserRequest) Validate() error {
    var errs uvalidator.ValidationErrors
    
    // Username 验证
    requiredRule := rule.NewRequired()
    if !requiredRule.Validate(r.Username) {
        errs = append(errs, uvalidator.NewFieldError(
            "Username",
            requiredRule.Name(),
            r.Username,
            requiredRule.GetMessage("Username", nil),
        ))
    }
    
    minRule := rule.NewMin(3)
    if !minRule.Validate(r.Username) {
        errs = append(errs, uvalidator.NewFieldError(
            "Username",
            minRule.Name(),
            r.Username,
            minRule.GetMessage("Username", map[string]string{"type": "string"}),
        ))
    }
    
    // Email 验证
    emailRule := rule.NewEmail()
    if !emailRule.Validate(r.Email) {
        errs = append(errs, uvalidator.NewFieldError(
            "Email",
            emailRule.Name(),
            r.Email,
            emailRule.GetMessage("Email", nil),
        ))
    }
    
    if errs.HasErrors() {
        return errs
    }
    return nil
}
```

#### 国际化支持

```go
// 设置全局语言为中文
uvalidator.SetLanguage(uvalidator.LanguageZH)

// 错误消息会显示中文
// "Username不能为空" 而不是 "Username is required"
```

#### 内置规则

- **基础规则**: Required, Min, Max, Len, Between
- **比较规则**: Gt, Gte, Lt, Lte
- **字符串规则**: Email, URL, Phone, Alpha, AlphaNum, Regex, UUID, JSON
- **网络规则**: IP, IPv4, IPv6, MAC, Domain, Port
- **中国特色规则**: IDCard, BankCard, UnifiedSocialCreditCode, PostalCode
- **安全规则**: StrongPassword, NoHTML, NoSQL, NoXSS

**详细文档**: [uvalidator/README.md](uvalidator/README.md)

---

### 9. udb/postgresql - PostgreSQL 数据库层

`udb/postgresql` 是一个零反射、高性能的 PostgreSQL 数据库层,基于 `github.com/jackc/pgx/v5`。

#### 核心特性

- 🚀 **零反射设计**: 手动实现 Scanner 接口,性能接近原生 pgx
- 🔗 **链路追踪**: 自动集成 ucontext,所有操作包含 trace_id
- 📊 **强大的查询构建器**: 支持 JOIN、GROUP BY、HAVING、DISTINCT 等复杂查询
- 💼 **完整的 CRUD 构建器**: Insert、Update、Delete 构建器,支持链式调用
- 🔄 **事务支持**: 完整的事务管理(Begin/Commit/Rollback)
- 📝 **慢查询日志**: 自动记录慢查询,支持可配置阈值
- ⚙️ **配置驱动**: 通过 uconfig 加载配置,支持连接池、查询超时等

#### 基本使用

**1. 配置文件**

```yaml
database:
  postgres:
    host: "localhost"
    port: 5432
    username: "postgres"
    password: "your_password"
    database: "myapp"
    ssl_mode: "disable"
    
    pool:
      max_conns: 25
      min_conns: 5
      max_conn_lifetime: "1h"
    
    query:
      slow_query_threshold: "1s"
    
    log:
      enabled: true
      slow_query: true
```

**2. 定义数据模型**

```go
type User struct {
    ID        int64
    Username  string
    Email     string
    Age       int
    CreatedAt time.Time
}

// 实现 Scanner 接口(零反射)
func (u *User) Scan(key string, value any) error {
    switch key {
    case "id":
        u.ID = uconv.ToInt64Def(value, 0)
    case "username":
        u.Username = uconv.ToString(value)
    case "email":
        u.Email = uconv.ToString(value)
    case "age":
        u.Age = uconv.ToIntDef(value, 0)
    case "created_at":
        u.CreatedAt = uconv.ToTimeDef(value, time.Time{})
    }
    return nil
}
```

**3. 连接和查询**

```go
import "github.com/whosafe/uf/udb/postgresql"

// 创建连接
conn, _ := postgresql.New(postgresql.GetConfig())
defer conn.Close()

ctx := ucontext.NewContext(context.Background())

// 查询单条记录
var user User
err := conn.Query(ctx).
    Table("users").
    Where("id = ?", 1).
    Scan(&user)

// 查询多条记录
results, _ := conn.Query(ctx).
    Table("users").
    Where("age > ?", 18).
    OrderBy("created_at").
    Limit(10).
    ScanAll(func() postgresql.Scanner { return &User{} })
```

#### 查询构建器

```go
// JOIN 查询
conn.Query(ctx).
    Select("u.id", "u.username", "COUNT(o.id) as order_count").
    Table("users u").
    LeftJoin("orders o", "u.id = o.user_id").
    GroupBy("u.id", "u.username").
    Having("COUNT(o.id) > ?", 0).
    OrderByDesc("order_count").
    ScanAll(newUserStats)

// WHERE 条件
conn.Query(ctx).
    Table("users").
    Where("age > ?", 18).
    WhereIn("status", []any{"active", "pending"}).
    WhereLike("username", "admin%").
    ScanAll(newUser)
```

#### CRUD 构建器

```go
// 插入
affected, _ := conn.Insert(ctx).
    Table("users").
    Columns("username", "email", "age").
    Values("alice", "alice@example.com", 25).
    Exec()

// 插入并返回数据
var newUser User
conn.Insert(ctx).
    Table("users").
    Columns("username", "email").
    Values("bob", "bob@example.com").
    ExecReturning(&newUser)

// 更新
affected, _ := conn.Update(ctx).
    Table("users").
    Set("age", 26).
    Set("email", "newemail@example.com").
    Where("id = ?", 1).
    Exec()

// 删除
affected, _ := conn.Delete(ctx).
    Table("users").
    Where("id = ?", 1).
    Exec()
```

#### 事务处理

```go
// 开始事务
tx, _ := conn.Begin(ctx)

// 执行操作
tx.Insert(ctx).
    Table("users").
    Columns("username", "email").
    Values("alice", "alice@example.com").
    Exec()

tx.Update(ctx).
    Table("accounts").
    Set("balance", "balance - 100").
    Where("user_id = ?", 1).
    Exec()

// 提交或回滚
if err != nil {
    tx.Rollback()
} else {
    tx.Commit()
}
```

**详细文档**: [udb/postgresql/README.md](udb/postgresql/README.md)

---

### 10. udb/redis - Redis 客户端封装

`udb/redis` 是一个零反射、高性能的 Redis 客户端封装，基于 `github.com/redis/go-redis/v9`。

**核心特性**：

- ✅ **零反射设计** - 所有类型转换使用 `uconv`，避免反射开销
- 🔗 **链路追踪集成** - 使用 `*ucontext.Context` 支持分布式追踪
- 📝 **统一日志** - 集成 `ulogger`，记录所有 Redis 命令和慢查询
- ⚡ **完整功能** - 支持字符串、哈希、列表、集合、有序集合、Pipeline、事务、Pub/Sub
- 🌐 **中文错误** - 所有错误提示均为中文
- 🎯 **类型安全** - 编译期检查，避免运行时错误

**快速开始**：

```go
import "github.com/whosafe/uf/udb/redis"

// 创建配置
config := redis.DefaultConfig()
config.Host = "localhost"
config.Port = 6379

// 创建连接
conn, err := redis.New(config)
if err != nil {
    log.Fatal(err)
}
defer conn.Close()

// 使用 Redis
ctx := ucontext.New()

// 字符串操作
conn.Set(ctx, "key", "value", 10*time.Minute)
value, _ := conn.Get(ctx, "key")

// 哈希操作
conn.HSet(ctx, "user:1", "name", "Alice", "age", 25)
name, _ := conn.HGet(ctx, "user:1", "name")

// Pipeline 批量操作
conn.Pipelined(ctx, func(pipe redis.Pipeliner) error {
    pipe.Set(ctx.Context(), "key1", "value1", 0)
    pipe.Set(ctx.Context(), "key2", "value2", 0)
    return nil
})
```

**配置示例**：

```yaml
database:
  redis:
    host: "localhost"
    port: 6379
    db: 0
    password: ""
    pool:
      pool_size: 10
      min_idle_conn: 5
    query:
      default_timeout: 5s
      slow_query_threshold: 100ms
    log:
      enabled: true
      level: "debug"
```

**详细文档**: [udb/redis/README.md](udb/redis/README.md)

---

## 📖 完整示例

### RESTful API 示例

```go
package main

import (
    "github.com/whosafe/uf/uconfig"
    "github.com/whosafe/uf/ucontext"
    "github.com/whosafe/uf/ulogger"
    "github.com/whosafe/uf/uprotocol/ubind"
    "github.com/whosafe/uf/uprotocol/uhttp"
    "github.com/whosafe/uf/uprotocol/unet"
    "github.com/whosafe/uf/uvalidator"
    "github.com/whosafe/uf/uvalidator/rule"
)

// 用户模型
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Age  int    `json:"age"`
}

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

// 创建用户请求
type CreateUserRequest struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func (r *CreateUserRequest) Bind(key string, value *ubind.Value) error {
    switch key {
    case "name":
        r.Name = value.Str()
    case "age":
        r.Age = value.Int()
    }
    return nil
}

func (r *CreateUserRequest) Validate() error {
    var errs uvalidator.ValidationErrors
    
    // 验证 Name
    requiredRule := rule.NewRequired()
    if !requiredRule.Validate(r.Name) {
        errs = append(errs, uvalidator.NewFieldError(
            "Name", requiredRule.Name(), r.Name,
            requiredRule.GetMessage("Name", nil),
        ))
    }
    
    // 验证 Age
    minRule := rule.NewMin(1)
    if !minRule.Validate(r.Age) {
        errs = append(errs, uvalidator.NewFieldError(
            "Age", minRule.Name(), r.Age,
            minRule.GetMessage("Age", map[string]string{"type": "number"}),
        ))
    }
    
    if errs.HasErrors() {
        return errs
    }
    return nil
}

func main() {
    // 加载配置
    uconfig.Load("config.yaml")
    
    // 创建服务器
    server := uhttp.New()
    
    // 应用默认中间件
    uhttp.ApplyDefaultMiddlewares(server)
    
    // 注册路由
    api := server.Group("/api/v1")
    {
        api.GET("/users", listUsers)
        api.POST("/users", createUser)
        api.GET("/users/:id", getUser)
        api.PUT("/users/:id", updateUser)
        api.DELETE("/users/:id", deleteUser)
    }
    
    // 启动服务器
    ulogger.Info("服务器启动", "address", ":8080")
    server.Start(":8080")
}

func listUsers(ctx *ucontext.Context, req unet.Request) error {
    httpResp := req.Response().(*uhttp.Response)
    
    users := []User{
        {ID: 1, Name: "Alice", Age: 25},
        {ID: 2, Name: "Bob", Age: 30},
    }
    
    return httpResp.Success(users)
}

func createUser(ctx *ucontext.Context, req unet.Request) error {
    httpResp := req.Response().(*uhttp.Response)
    
    var createReq CreateUserRequest
    if err := req.Bind(&createReq); err != nil {
        return httpResp.BadRequest("无效的请求数据")
    }
    
    if err := createReq.Validate(); err != nil {
        return httpResp.BadRequest(err.Error())
    }
    
    user := User{
        ID:   3,
        Name: createReq.Name,
        Age:  createReq.Age,
    }
    
    ulogger.InfoCtx(ctx.Context(), "创建用户", "user", user.Name)
    
    return httpResp.SuccessWithMessage("创建成功", user)
}

func getUser(ctx *ucontext.Context, req unet.Request) error {
    httpReq := req.(*uhttp.Request)
    httpResp := req.Response().(*uhttp.Response)
    
    id := httpReq.Param("id")
    
    user := User{ID: 1, Name: "Alice", Age: 25}
    
    ulogger.InfoCtx(ctx.Context(), "获取用户", "id", id)
    
    return httpResp.Success(user)
}

func updateUser(ctx *ucontext.Context, req unet.Request) error {
    httpReq := req.(*uhttp.Request)
    httpResp := req.Response().(*uhttp.Response)
    
    id := httpReq.Param("id")
    
    var updateReq CreateUserRequest
    if err := req.Bind(&updateReq); err != nil {
        return httpResp.BadRequest("无效的请求数据")
    }
    
    if err := updateReq.Validate(); err != nil {
        return httpResp.BadRequest(err.Error())
    }
    
    ulogger.InfoCtx(ctx.Context(), "更新用户", "id", id)
    
    return httpResp.SuccessWithMessage("更新成功", nil)
}

func deleteUser(ctx *ucontext.Context, req unet.Request) error {
    httpReq := req.(*uhttp.Request)
    httpResp := req.Response().(*uhttp.Response)
    
    id := httpReq.Param("id")
    
    ulogger.InfoCtx(ctx.Context(), "删除用户", "id", id)
    
    return httpResp.SuccessWithMessage("删除成功", nil)
}
```

---

## 🎯 最佳实践

### 1. 项目结构建议

```
myapp/
├── cmd/
│   └── server/
│       └── main.go          # 程序入口
├── internal/
│   ├── handler/             # 处理器
│   │   ├── user.go
│   │   └── auth.go
│   ├── model/               # 数据模型
│   │   └── user.go
│   ├── service/             # 业务逻辑
│   │   └── user.go
│   └── middleware/          # 自定义中间件
│       └── auth.go
├── config/
│   ├── config.yaml          # 配置文件
│   ├── config.dev.yaml      # 开发环境配置
│   └── config.prod.yaml     # 生产环境配置
├── logs/                    # 日志目录
├── uploads/                 # 上传文件目录
└── go.mod
```

### 2. 配置文件组织

**开发环境** (`config.dev.yaml`):

```yaml
server:
  address: ":8080"
  middleware:
    enable_trace: true
    enable_logger: true
    enable_recovery: true

logger:
  level: "debug"
  stdout: true
  path: ""  # 仅输出到终端
```

**生产环境** (`config.prod.yaml`):

```yaml
server:
  address: ":8080"
  middleware:
    enable_trace: true
    enable_logger: true
    enable_recovery: true
    enable_cors: true

logger:
  level: "info"
  stdout: false
  path: "/var/log/myapp"
  file: "app.log"
  rotate_size: 104857600  # 100MB
  rotate_backup_limit: 30
  rotate_backup_compress: 6
```

### 3. 错误处理

使用统一的错误响应格式：

```go
func handler(ctx *ucontext.Context, req unet.Request) error {
    httpResp := req.Response().(*uhttp.Response)
    
    // 业务逻辑
    user, err := getUserByID(id)
    if err != nil {
        // 记录错误日志
        ulogger.ErrorCtx(ctx.Context(), "获取用户失败", "error", err)
        
        // 返回错误响应
        return httpResp.InternalServerError("服务器内部错误")
    }
    
    return httpResp.Success(user)
}
```

### 4. 性能优化

**使用对象池**：

```go
// Request 和 Response 对象自动使用 sync.Pool
// 无需手动管理

// umarshal 使用对象池
w := umarshal.AcquireWriter()
defer umarshal.ReleaseWriter(w)
```

**避免反射**：

```go
// ✅ 推荐：实现 Bind 接口
func (u *User) Bind(key string, value *ubind.Value) error {
    // 手动绑定
}

// ❌ 不推荐：使用反射
json.Unmarshal(data, &user)
```

### 5. 日志最佳实践

```go
// 使用结构化日志
ulogger.InfoCtx(ctx, "处理请求",
    "method", "POST",
    "path", "/api/users",
    "user_id", userID,
)

// 使用 Context 方法自动注入 trace_id
ulogger.InfoCtx(ctx.Context(), "业务处理", "key", "value")

// 错误日志包含详细信息
ulogger.ErrorCtx(ctx.Context(), "数据库查询失败",
    "error", err,
    "sql", query,
    "params", params,
)
```

---

## ❓ 常见问题

### Q1: 如何配置 HTTPS？

**A**: 在配置文件中设置 TLS 证书：

```yaml
server:
  protocol: "https"
  address: ":443"
  tls:
    cert_file: "/path/to/cert.pem"
    key_file: "/path/to/key.pem"
```

### Q2: 如何启用 CORS？

**A**: 在配置文件中启用 CORS 中间件：

```yaml
server:
  middleware:
    enable_cors: true
    cors:
      allow_origins: "*"
      allow_methods: "GET,POST,PUT,DELETE,PATCH,HEAD,OPTIONS"
      allow_headers: "*"
      allow_credentials: false
```

### Q3: 如何使用 Redis 存储 Session？

**A**: 创建 Redis 存储并配置 Session 管理器：

```go
import "github.com/redis/go-redis/v9"

redisClient := redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
})

store := uhttp.NewRedisStore(redisClient, "session:", 30*time.Minute)
sessionMgr := uhttp.NewSessionManager(store, "session_id", 3600)

// 在服务器中设置
server.SetSessionManager(sessionMgr)
```

### Q4: 如何自定义中间件？

**A**: 实现 `MiddlewareFunc` 接口：

```go
func AuthMiddleware() unet.MiddlewareFunc {
    return func(next unet.HandlerFunc) unet.HandlerFunc {
        return func(ctx *ucontext.Context, req unet.Request) error {
            httpReq := req.(*uhttp.Request)
            httpResp := req.Response().(*uhttp.Response)
            
            // 验证 token
            token := httpReq.Header("Authorization")
            if token == "" {
                return httpResp.Unauthorized("未授权")
            }
            
            // 验证通过，继续处理
            return next(ctx, req)
        }
    }
}

// 使用
server.Use(AuthMiddleware())
```

### Q5: 如何处理文件上传大小限制？

**A**: 使用 `FileUploadConfig` 配置：

```go
path, err := httpReq.SaveUploadedFileWithConfig(file, &uhttp.FileUploadConfig{
    MaxSize:     10 << 20, // 10MB
    AllowedExts: []string{".jpg", ".png", ".gif"},
    UploadDir:   "./uploads",
})
```

### Q6: 如何配置日志轮转？

**A**: 在配置文件或代码中设置轮转参数：

```yaml
logger:
  path: "./logs"
  file: "app.log"
  rotate_size: 104857600      # 100MB
  rotate_backup_limit: 10     # 保留 10 个备份
  rotate_backup_expire: 604800 # 7 天
  rotate_backup_compress: 6   # gzip 压缩级别
```

### Q7: 如何在微服务间传递 Trace ID？

**A**: 使用 `ucontext` 的 HTTP 传播功能：

```go
// 服务 A（调用方）
func callServiceB(ctx context.Context) {
    req, _ := http.NewRequest("GET", "http://service-b/api", nil)
    tc := ucontext.FromContext(ctx)
    ucontext.InjectHTTPHeaders(req.Header, tc)
    
    resp, _ := http.DefaultClient.Do(req)
}

// 服务 B（接收方）
func handler(w http.ResponseWriter, r *http.Request) {
    tc := ucontext.ExtractHTTPHeaders(r.Header)
    ctx := ucontext.WithContext(r.Context(), tc)
    
    // 使用带追踪信息的 context
    ulogger.InfoCtx(ctx, "处理请求")
}
```

---

## 📚 模块详细文档

- [unet - 网络层抽象](uprotocol/unet/README.md)
- [uhttp - HTTP 服务器](uprotocol/uhttp/README.md)
- [uconfig - 配置管理](uconfig/README.md)
- [ulogger - 日志系统](ulogger/README.md)
- [ucontext - 链路追踪](ucontext/README.md)
- [ubind - 数据绑定](uprotocol/ubind/README.md)
- [umarshal - JSON 序列化](uprotocol/umarshal/README.md)
- [uvalidator - 数据验证](uvalidator/README.md)
- [udb/postgresql - PostgreSQL 数据库层](udb/postgresql/README.md)
- [udb/redis - Redis 客户端封装](udb/redis/README.md)

---

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License

---

## 🔗 示例项目

查看 `example/` 目录获取完整示例：

- `example/uhttp/01_basic` - 基础 HTTP 服务器
- `example/uhttp/02_middleware` - 中间件使用
- `example/uhttp/03_restful` - RESTful API
- `example/uhttp/04_advanced` - 高级功能(Session、文件上传等)
- `example/udb/postgresql/manual_test` - PostgreSQL 完整手动测试
- `example/udb/redis` - Redis 客户端使用示例
