# uconfig - 零依赖高性能配置库

`uconfig` 是一个轻量级、高性能的 Go 配置加载库，**完全零依赖**，采用自研 YAML 解析器，性能比标准库快 **4-5 倍**。

## ✨ 特性

- 🚀 **极致性能**: 比 `yaml.v3` 快 4-5 倍，内存占用减少 40%
- 📦 **零依赖**: 无需任何第三方库，完全自包含
- 🎯 **类型安全**: 强类型 API，编译期检查
- 🔧 **灵活解析**: 支持自定义解析逻辑
- 📝 **简洁 API**: 易于使用，学习成本低
- 🌳 **嵌套支持**: 完美支持嵌套结构和数组
- 🔌 **回调机制**: 支持未知配置项的被动回调

## 📊 性能对比

| 方式 | 耗时 (ns/op) | 内存 (B/op) | 分配次数 | 性能提升 |
|:---|:---|:---|:---|:---|
| **uconfig** | **12,211** | **7,912** | **60** | **基准** |
| yaml.v3 | 53,414 | 13,640 | 200 | 慢 4.4x |

## 📦 安装

```bash
go get github.com/whosafe/uf/uconfig
```

## 🚀 快速开始

### 1. 准备配置文件 (config.yaml)

```yaml
server:
  host: "0.0.0.1"
  port:
    - 8080
    - 8081
    - 8082

database:
  dsn: "user:pass@tcp(localhost:3306)/dbname"
  max_open: 100
  logger:
    level: "debug"
    path: "/var/log/app.log"
```

### 2. 定义配置结构

```go
package main

import (
    "github.com/whosafe/uf/uconfig"
    "github.com/whosafe/uf/uconv"
)

// ServerConfig 服务器配置
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

// DatabaseConfig 数据库配置
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

type LoggerConfig struct {
    Level string
    Path  string
}

func (l *LoggerConfig) UnmarshalYAML(key string, value *uconfig.Node) error {
    switch key {
    case "level":
        l.Level = value.String()
    case "path":
        l.Path = value.String()
    }
    return nil
}
```

### 3. 加载配置

```go
func main() {
    var srvCfg ServerConfig
    var dbCfg DatabaseConfig

    // 注册配置解析器
    uconfig.Register("server", srvCfg.UnmarshalYAML)
    uconfig.Register("database", dbCfg.UnmarshalYAML)

    // 加载配置文件
    if err := uconfig.Load("config.yaml"); err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Server: %+v\n", srvCfg)
    fmt.Printf("Database: %+v\n", dbCfg)
}
```

## 📚 核心 API

### 注册与加载

```go
// Register 注册配置项解析器
func Register(key string, cb ICallback)

// RegisterUnknown 注册未知配置项的回调
func RegisterUnknown(cb ICallback)

// Load 加载配置文件
func Load(path string) error

// ParseConfig 解析配置内容（从字节流）
func ParseConfig(data []byte) error

// Callback 手动触发配置回调
func Callback(key string, cb ICallback) error
```

### Node 方法

```go
// String 获取节点的字符串值
func (n *Node) String() string

// Iter 遍历数组节点
func (n *Node) Iter(cb func(i int, v *Node) error) error

// Decode 解码到实现了 Unmarshaler 的结构体
func (n *Node) Decode(v any) error
```

### Unmarshaler 接口

```go
type Unmarshaler interface {
    UnmarshalYAML(key string, value *Node) error
}
```

## 🎯 高级用法

### 处理未知配置项

```go
uconfig.RegisterUnknown(func(key string, value *uconfig.Node) error {
    fmt.Printf("Unknown config: %s\n", key)
    return nil
})
```

### 手动触发回调

```go
err := uconfig.Callback("server", func(key string, value *uconfig.Node) error {
    fmt.Printf("Key: %s, Value: %s\n", key, value.String())
    return nil
})
```

### 解析数组

```go
func (c *Config) UnmarshalYAML(key string, value *uconfig.Node) error {
    if key == "items" {
        return value.Iter(func(i int, v *uconfig.Node) error {
            c.Items = append(c.Items, v.String())
            return nil
        })
    }
    return nil
}
```

### 嵌套结构解析

```go
func (c *Config) UnmarshalYAML(key string, value *uconfig.Node) error {
    if key == "nested" {
        // 自动递归解析
        return value.Decode(&c.Nested)
    }
    return nil
}
```

## 📂 项目结构

```
uconfig/
├── config.go       # 核心配置加载逻辑
├── node.go         # Node 结构和方法
├── parser.go       # YAML 解析器实现
├── registry.go     # 注册表管理
├── config_test.go  # 测试用例
└── README.md       # 本文档
```

## 🔍 设计理念

### 为什么不用 yaml.v3？

1. **性能**: `yaml.v3` 为了支持完整的 YAML 1.2 规范，包含了大量不常用的特性（锚点、别名、流式风格等），导致性能开销大
2. **依赖**: 引入第三方依赖增加了项目复杂度
3. **灵活性**: 自研解析器可以针对配置场景优化，提供更灵活的 API

### 支持的 YAML 特性

`uconfig` 专注于配置文件场景，支持：

- ✅ 键值对 (Map)
- ✅ 数组 (Sequence)
- ✅ 标量值 (Scalar)
- ✅ 嵌套结构
- ✅ 注释
- ✅ 引号字符串

**不支持**（配置场景不常用）：

- ❌ 锚点和别名
- ❌ 流式风格
- ❌ 多文档
- ❌ 复杂键

## 🧪 测试

```bash
# 运行测试
go test -v ./uconfig

# 运行示例
go run ./example/uconfig/main.go

# 性能测试
go test -bench=. -benchmem ./uconfig
```

## 📝 最佳实践

1. **使用类型转换工具**: 配合 `uconv` 包进行安全的类型转换
2. **错误处理**: 在 `UnmarshalYAML` 中妥善处理错误
3. **嵌套结构**: 对于复杂嵌套，使用 `Decode` 方法递归解析
4. **数组遍历**: 使用 `Iter` 方法处理数组，避免手动索引

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License

## 🔗 相关项目

- [uconv](../uconv) - 类型转换工具库
- [uerror](../uerror) - 错误处理增强库
