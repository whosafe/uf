# uhttp 配置示例

本目录包含多个配置文件示例,用于不同场景的 HTTP 服务器配置。

## 配置文件

- `config.yaml` - 标准 HTTP 服务器配置
- `config.https.yaml` - HTTPS 服务器配置
- `config.minimal.yaml` - 最小化配置示例

## 配置项说明

### 基础配置

- `name` - 服务名称
- `protocol` - 协议类型 (http/https/tcp/quic)
- `address` - 监听地址 (如 `:8080`)
- `unix_path` - UNIX 套接字路径 (可选)

### 超时配置

- `read_timeout` - 读取超时
- `write_timeout` - 写入超时
- `idle_timeout` - 空闲超时

### 大小限制

- `max_header_bytes` - 最大请求头字节数
- `max_body_bytes` - 最大请求体字节数
- `max_form_bytes` - 最大表单字节数

### TLS 配置

仅在 `protocol: https` 时需要:

```yaml
tls:
  cert_file: "./cert.pem"
  key_file: "./key.pem"
```

### HTTP/2 配置

```yaml
http2:
  enabled: true
  max_concurrent_streams: 250
```

### 静态文件服务

```yaml
static:
  enabled: true
  root: "./public"
  prefix: "/static"
  index: ["index.html"]
```

### Cookie 配置

```yaml
cookie:
  domain: ""
  path: "/"
  max_age: 86400
  secure: false
  http_only: true
  same_site: "lax"
```

### Session 配置

```yaml
session:
  enabled: true
  provider: "memory"
  cookie_name: "session_id"
  max_age: 3600
```

### 日志配置

日志分为两类:

**访问日志** - 记录所有请求:

```yaml
access_log:
  enabled: true
  level: "info"
  format: "json"
  output: "stdout"
  file_path: "./logs/access.log"
```

**错误日志** - 仅记录错误请求 (4xx, 5xx):

```yaml
error_log:
  enabled: true
  level: "error"
  format: "json"
  output: "stderr"
  file_path: "./logs/error.log"
```

## 使用方式

```go
import (
    "github.com/whosafe/uf/uconfig"
    "github.com/whosafe/uf/uprotocol/uhttp"
)

func main() {
    // 加载配置
    cfg := uconfig.Load("config.yaml")
    
    // 创建服务器
    app := uhttp.NewWithConfig(cfg.Server)
    
    // 启动服务器
    app.Start()
}
```
