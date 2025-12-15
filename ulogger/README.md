# ulogger - 高性能日志库

`ulogger` 是一个功能完善的 Go 日志库，基于标准库 `log/slog`，支持文件轮转、压缩、多输出等企业级特性。

## ✨ 特性

- 🎯 **基于 slog**: 完全兼容 Go 1.21+ 的标准日志接口
- 🔄 **文件轮转**: 支持按大小和时间自动轮转
- 📦 **备份管理**: 自动清理过期备份，支持数量和时间限制
- 🗜️ **压缩支持**: 可选的 gzip 压缩（0-9 级别）
- 📤 **多输出**: 同时输出到文件和终端
- 🎨 **多格式支持**: 支持文本、JSON 和自定义格式
- ⚡ **高性能**: 异步轮转和压缩，不阻塞日志写入
- 🔧 **易于使用**: 简洁的 API，开箱即用
- 🔌 **uconfig 集成**: 支持从配置文件自动加载

## 🎨 多格式日志支持

### 支持的格式

`ulogger` 支持多种日志格式：

1. **文本格式** (默认)
   - 标准格式：`2006-01-02 15:04:05 [INFO] file.go:10 message key=value`
   - 简洁格式：`15:04:05 message key=value`

2. **JSON 格式**
   - 结构化日志输出
   - 便于日志分析工具处理

3. **自定义格式**
   - 实现 `Formatter` 接口
   - 完全自定义日志格式

### JSON 格式

```go
config := &ulogger.Config{
    Path:   "./logs",
    File:   "app.json",
    Format: "json",
    Level:  slog.LevelInfo,
}

logger, _ := ulogger.New(config)
logger.Info("用户登录", "user", "alice", "ip", "192.168.1.1")
```

输出：
```json
{"time":"2025-12-12T14:40:59.720+08:00","level":"INFO","msg":"用户登录","source":"main.go:10","attrs":{"user":"alice","ip":"192.168.1.1"}}
```

### 自定义格式

```go
// 实现 Formatter 接口
type MyFormatter struct{}

func (f *MyFormatter) Format(r slog.Record, config *ulogger.Config) ([]byte, error) {
    return []byte(fmt.Sprintf("[%s] %s\n", r.Level, r.Message)), nil
}

// 使用自定义格式化器
config := &ulogger.Config{
    Format:    "custom",
    Formatter: &MyFormatter{},
}

logger, _ := ulogger.New(config)
logger.Info("自定义格式消息")
// 输出: [INFO] 自定义格式消息
```

### 配置文件方式

```yaml
logger:
  path: "./logs"
  file: "app.log"
  format: "json"  # 或 "text"
  level: "info"
```

```go
ulogger.Register()
uconfig.Load("config.yaml")
ulogger.Info("使用配置文件中的设置")
```

## 📦 安装

```bash
go get github.com/whosafe/uf/ulogger
```

## 🚀 快速开始

### 基本用法

```go
package main

import (
    "log/slog"
    "github.com/whosafe/uf/ulogger"
)

func main() {
    // 使用默认配置（仅输出到终端）
    logger, _ := ulogger.New(ulogger.DefaultConfig())
    defer logger.Close()

    logger.Info("Hello, ulogger!")
    logger.Warn("This is a warning", "key", "value")
}
```

### 输出到文件

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

### JSON 格式日志

```go
config := &ulogger.Config{
    Path:   "./logs",
    File:   "app.json",
    Format: "json", // 使用 JSON 格式
    Level:  slog.LevelInfo,
}

logger, _ := ulogger.New(config)
defer logger.Close()

logger.Info("API 请求", "method", "GET", "status", 200)
// 输出: {"time":"2025-12-12T14:40:59.720+08:00","level":"INFO","msg":"API 请求","source":"main.go:15","attrs":{"method":"GET","status":200}}
```

## 📚 配置说明

### Config 结构

```go
type Config struct {
    Path                 string     // 日志文件路径。默认为空，表示关闭，仅输出到终端
    File                 string     // 日志文件格式。默认为"2006-01-02.log"
    Prefix               string     // 日志内容输出前缀。默认为空
    Level                slog.Level // 日志输出级别
    UseStandardLogFormat bool       // 是否使用标准日志格式。默认true
    ShortFile            bool       // 日志文件是否只输出文件名。默认false
    Stdout               bool       // 日志是否同时输出到终端。默认true
    RotateSize           int        // 按照日志文件大小对文件进行滚动切分。默认为0，表示关闭滚动切分特性
    RotateExpire         int64      // 按照日志文件时间间隔对文件滚动切分。默认为0，表示关闭滚动切分特性
    RotateBackupLimit    int        // 按照切分的文件数量清理切分文件，当滚动切分特性开启时有效。默认为0，表示不备份，切分则删除
    RotateBackupExpire   int        // 按照切分的文件有效期清理切分文件，当滚动切分特性开启时有效。默认为0，表示不备份，切分则删除
    RotateBackupCompress uint16     // 滚动切分文件的压缩比（0-9）。默认为0，表示不压缩
}
```

### 配置详解

#### Path & File
- `Path`: 日志文件目录，为空则仅输出到终端
- `File`: 文件名，支持 Go 时间格式化（如 `2006-01-02.log`）

#### 日志级别
- `Level`: 使用 `slog.Level`
  - `slog.LevelDebug`: 调试
  - `slog.LevelInfo`: 信息（默认）
  - `slog.LevelWarn`: 警告
  - `slog.LevelError`: 错误

#### 格式化选项
- `UseStandardLogFormat`: 
  - `true`: `2006-01-02 15:04:05 [INFO] main.go:10 message key=value`
  - `false`: `15:04:05 message key=value`
- `ShortFile`: 
  - `true`: 只显示文件名 `main.go:10`
  - `false`: 显示完整路径 `/path/to/main.go:10`
- `Prefix`: 在每条日志前添加前缀，如 `[APP]`

#### 输出选项
- `Stdout`: 
  - `true`: 同时输出到文件和终端
  - `false`: 仅输出到文件

#### 轮转选项
- `RotateSize`: 文件大小阈值（字节），超过则轮转
- `RotateExpire`: 时间间隔（秒），到期则轮转

#### 备份管理
- `RotateBackupLimit`: 保留的备份文件数量
- `RotateBackupExpire`: 备份文件有效期（秒）
- `RotateBackupCompress`: gzip 压缩级别（0-9，0 表示不压缩）

## 🎯 使用示例

### 1. 全局 Logger

```go
// 设置全局 Logger
config := &ulogger.Config{
    Path:   "./logs",
    File:   "app.log",
    Prefix: "[APP]",
}
logger, _ := ulogger.New(config)
ulogger.SetDefault(logger)

// 使用全局函数
ulogger.Info("global message")
ulogger.Error("error occurred", "error", err)

// 标准库 slog 也会使用我们的 Logger
slog.Info("this uses our logger too")
```

### 2. 带属性的日志

```go
logger.Info("user login", 
    "user", "alice",
    "ip", "192.168.1.1",
    "timestamp", time.Now())
```

### 3. 子 Logger

```go
// 创建带有固定属性的子 Logger
requestLogger := logger.With("request_id", "12345")
requestLogger.Info("processing request")
requestLogger.Info("request completed")

// 创建带分组的子 Logger
dbLogger := logger.WithGroup("database")
dbLogger.Info("query executed", "duration", "100ms")
```

### 4. 按时间轮转

```go
config := &ulogger.Config{
    Path:         "./logs",
    File:         "2006-01-02.log", // 每天一个文件
    RotateExpire: 86400,             // 24 小时轮转
}
```

### 5. 完整配置示例

```go
config := &ulogger.Config{
    Path:                 "./logs",
    File:                 "app.log",
    Prefix:               "[MyApp]",
    Level:                slog.LevelDebug,
    UseStandardLogFormat: true,
    ShortFile:            true,
    Stdout:               true,
    RotateSize:           100 * 1024 * 1024, // 100MB
    RotateExpire:         0,
    RotateBackupLimit:    10,
    RotateBackupExpire:   7 * 24 * 3600, // 7 天
    RotateBackupCompress: 6,
}

logger, err := ulogger.New(config)
if err != nil {
    panic(err)
}
defer logger.Close()
```

## 📂 项目结构

```
ulogger/
├── config.go         # 配置定义
├── logger.go         # 核心 Logger 实现
├── handler.go        # 自定义 slog Handler
├── rotate.go         # 文件轮转实现
├── backup.go         # 备份管理
├── utils.go          # 工具函数
├── logger_test.go    # 基础测试
├── rotate_test.go    # 轮转测试
└── README.md         # 本文档
```

## 🧪 测试

```bash
# 运行所有测试
go test -v ./ulogger

# 运行特定测试
go test -v -run TestRotateBySize ./ulogger

# 运行示例
go run ./example/ulogger/main.go
```

## 📝 最佳实践

### 1. 生产环境配置

```go
config := &ulogger.Config{
    Path:                 "/var/log/myapp",
    File:                 "app.log",
    Level:                slog.LevelInfo, // 生产环境使用 Info
    UseStandardLogFormat: true,
    ShortFile:            false, // 完整路径便于定位
    Stdout:               false, // 生产环境不输出到终端
    RotateSize:           100 * 1024 * 1024, // 100MB
    RotateBackupLimit:    30,
    RotateBackupCompress: 6,
}
```

### 2. 开发环境配置

```go
config := &ulogger.Config{
    Level:  slog.LevelDebug, // 开发环境显示所有日志
    Stdout: true,            // 输出到终端便于调试
    UseStandardLogFormat: false, // 简洁格式
}
```

### 3. 优雅关闭

```go
logger, _ := ulogger.New(config)
defer func() {
    logger.Sync()  // 确保所有日志写入磁盘
    logger.Close() // 关闭文件和停止轮转
}()
```

### 4. 错误处理

```go
logger, err := ulogger.New(config)
if err != nil {
    // 降级到标准输出
    logger, _ = ulogger.New(&ulogger.Config{Stdout: true})
}
```

## 🔍 工作原理

### 文件轮转流程

1. **按大小轮转**: 每次写入前检查文件大小，超过阈值则轮转
2. **按时间轮转**: 后台定时检查（每 10 秒），到期则轮转
3. **轮转步骤**:
   - 关闭当前文件
   - 重命名为备份文件（添加时间戳）
   - 创建新文件
   - 异步压缩备份（如果启用）
   - 异步清理旧备份

### 备份文件命名

```
app.log                    # 当前日志文件
app.20231212-143025.log    # 备份文件（未压缩）
app.20231212-143025.log.gz # 备份文件（已压缩）
```

## ⚠️ 注意事项

1. **并发安全**: Logger 是并发安全的，可以在多个 goroutine 中使用
2. **资源清理**: 务必调用 `Close()` 以确保资源正确释放
3. **压缩性能**: 压缩是异步的，不会阻塞日志写入
4. **时间格式**: `File` 字段使用 Go 的时间格式化语法
5. **备份清理**: 清理操作是异步的，可能有短暂延迟

## 📄 许可证

MIT License

## 🔗 相关项目

- [uconfig](../uconfig) - 零依赖配置库
- [uconv](../uconv) - 类型转换工具库
- [uerror](../uerror) - 错误处理增强库
