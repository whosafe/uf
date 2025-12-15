# ucontext - 分布式链路追踪上下文包

`ucontext` 是一个轻量级的分布式链路追踪包，支持雪花算法 ID 生成、采样控制、HTTP 传播和 Logger 集成。

## ✨ 特性

- 🔢 **雪花算法**: 分布式唯一 ID 生成
- 🔗 **链路追踪**: Trace ID、Span ID、Parent Span ID
- 📊 **采样控制**: 可配置的采样率，支持强制采样
- 🌐 **HTTP 传播**: 跨服务传递追踪信息
- 📝 **Logger 集成**: 自动注入追踪信息到日志
- ⚡ **高性能**: 并发安全，低开销
- 🔧 **易于使用**: 简洁的 API，开箱即用

## 📦 安装

```bash
go get github.com/whosafe/uf/ucontext
```

## 🚀 快速开始

### 基本使用

```go
package main

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

### Logger 集成

```go
import (
    "github.com/whosafe/uf/ucontext"
    "github.com/whosafe/uf/ulogger"
)

func main() {
    ctx := ucontext.NewContext(context.Background())

    // 使用全局 Logger 的 Context 方法
    ulogger.InfoCtx(ctx, "处理请求", "key", "value")
    ulogger.DebugCtx(ctx, "调试信息")
    ulogger.WarnCtx(ctx, "警告信息")
    ulogger.ErrorCtx(ctx, "错误信息")
}

// 输出包含 trace_id 和 span_id:
// 2025-12-12 15:31:28 [INFO] main.go:10 处理请求 key=value trace_id=522314532622700544 span_id=522314532622700545
```

### HTTP 传播

```go
// 服务端：提取追踪信息
func handler(w http.ResponseWriter, r *http.Request) {
    tc := ucontext.ExtractHTTPHeaders(r.Header)
    ctx := ucontext.WithContext(r.Context(), tc)

    // 使用带追踪信息的 context
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

### 嵌套 Span

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

### 采样控制

```go
// 设置采样率为 50%
ucontext.SetSamplingRate(0.5)

// 创建上下文（50% 概率被采样）
ctx := ucontext.NewContext(context.Background())

// 检查是否被采样
if ucontext.IsSampled(ctx) {
    // 记录详细追踪信息
}

// 强制采样（忽略采样率）
ctx = ucontext.ForceSample(ctx)
```

## 📚 API 文档

### 核心函数

#### ID 生成

```go
// 初始化雪花算法（应在程序启动时调用一次）
func InitSnowflake(workerID int64) error

// 生成唯一 ID
func GenerateID() string
```

#### Context 操作

```go
// 创建新的追踪上下文
func NewTraceContext() *TraceContext

// 创建带追踪信息的 context
func NewContext(parent context.Context) context.Context

// 创建子 Span
func NewSpan(parent context.Context) context.Context

// 从 context 提取追踪信息
func FromContext(ctx context.Context) *TraceContext

// 将追踪信息注入 context
func WithContext(ctx context.Context, tc *TraceContext) context.Context
```

#### 采样控制

```go
// 设置采样率 (0.0 - 1.0)
func SetSamplingRate(rate float64)

// 获取当前采样率
func GetSamplingRate() float64

// 强制采样
func ForceSample(ctx context.Context) context.Context

// 检查是否被采样
func IsSampled(ctx context.Context) bool
```

#### HTTP 传播

```go
// 注入到 HTTP Header
func InjectHTTPHeaders(header http.Header, tc *TraceContext)

// 从 HTTP Header 提取
func ExtractHTTPHeaders(header http.Header) *TraceContext

// HTTP 中间件
func HTTPMiddleware(next http.Handler) http.Handler
```

### TraceContext 结构

```go
type TraceContext struct {
    TraceID      string            // 追踪 ID
    SpanID       string            // 当前 Span ID
    ParentSpanID string            // 父 Span ID
    RequestID    string            // 请求 ID
    StartTime    time.Time         // 开始时间
    Sampled      bool              // 是否采样
    Metadata     map[string]string // 元数据
}

// 设置元数据
func (tc *TraceContext) SetMetadata(key, value string)

// 获取元数据
func (tc *TraceContext) GetMetadata(key string) string

// 获取持续时间
func (tc *TraceContext) Duration() time.Duration
```

## 🎯 使用场景

### 1. 微服务链路追踪

```go
// 服务 A
func serviceA(w http.ResponseWriter, r *http.Request) {
    ctx := ucontext.NewContext(r.Context())
    ulogger.InfoCtx(ctx, "服务 A 收到请求")

    // 调用服务 B
    callServiceB(ctx)
}

// 服务 B
func serviceB(w http.ResponseWriter, r *http.Request) {
    tc := ucontext.ExtractHTTPHeaders(r.Header)
    ctx := ucontext.WithContext(r.Context(), tc)
    ulogger.InfoCtx(ctx, "服务 B 收到请求")
    // 日志会包含相同的 Trace ID
}
```

### 2. 数据库操作追踪

```go
func handleRequest(ctx context.Context) {
    // 查询用户
    userCtx := ucontext.NewSpan(ctx)
    user := queryUser(userCtx)

    // 查询订单
    orderCtx := ucontext.NewSpan(ctx)
    orders := queryOrders(orderCtx)
}
```

### 3. 性能监控

```go
func processTask(ctx context.Context) {
    tc := ucontext.FromContext(ctx)
    defer func() {
        duration := tc.Duration()
        ulogger.InfoCtx(ctx, "任务完成", "duration_ms", duration.Milliseconds())
    }()

    // 执行任务
}
```

## 🔧 配置

### 雪花算法配置

```go
// Worker ID 范围: 0-1023
// 建议每个服务实例使用不同的 Worker ID
ucontext.InitSnowflake(1)
```

### 采样率配置

```go
// 生产环境建议 10%-30%
ucontext.SetSamplingRate(0.2)

// 开发环境建议 100%
ucontext.SetSamplingRate(1.0)
```

## 📊 性能

- ID 生成: ~100万/秒
- Context 操作: 纳秒级
- 采样判断: 纳秒级
- 并发安全: 无锁设计（除 ID 生成）

## 🤝 与其他组件集成

### ulogger 集成

自动支持，使用 `ulogger.InfoCtx(ctx, ...)` 等方法即可自动注入追踪信息。

```go
ctx := ucontext.NewContext(context.Background())
ulogger.InfoCtx(ctx, "处理请求", "user", "alice")
// 输出会自动包含 trace_id 和 span_id
```

### HTTP 框架集成

```go
// 使用中间件
http.Handle("/api", ucontext.HTTPMiddleware(handler))
```

## 📝 最佳实践

1. **在程序启动时初始化雪花算法**
   ```go
   func main() {
       ucontext.InitSnowflake(getWorkerID())
       // ...
   }
   ```

2. **在 HTTP Handler 入口创建追踪上下文**
   ```go
   func handler(w http.ResponseWriter, r *http.Request) {
       ctx := ucontext.NewContext(r.Context())
       // ...
   }
   ```

3. **使用 Context 传递追踪信息**
   ```go
   func businessLogic(ctx context.Context) {
       ulogger.InfoCtx(ctx, "业务逻辑处理")
       // ...
   }
   ```

4. **为重要操作创建子 Span**
   ```go
   func importantOperation(ctx context.Context) {
       spanCtx := ucontext.NewSpan(ctx)
       // ...
   }
   ```

## 📄 License

MIT License
