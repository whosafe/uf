# unet - 网络层抽象

协议无关的网络层抽象接口,支持 HTTP、TCP、QUIC 等多种协议。

## 核心设计

**两参数设计** - `HandlerFunc(ctx context.Context, req Request)` 实现职责分离:

- `context.Context` - 链路追踪上下文
- `Request` - 请求对象(包含请求数据和响应能力)

## 核心接口

### Protocol - 协议类型

```go
type Protocol string

const (
    ProtocolHTTP  Protocol = "http"
    ProtocolHTTPS Protocol = "https"
    ProtocolTCP   Protocol = "tcp"
    ProtocolQUIC  Protocol = "quic"
)
```

### HandlerFunc - 处理器函数

```go
type HandlerFunc func(ctx context.Context, req Request) error
```

### Request - 请求接口

```go
type Request interface {
    Protocol() Protocol
    RemoteAddr() net.Addr
    LocalAddr() net.Addr
    Get(key string) (any, bool)
    Set(key string, value any)
    Bind(obj ubind.Binder) error
    Response() Response
}
```

### Response - 响应接口

```go
type Response interface {
    JSON(code int, data any) error
    String(code int, text string) error
    Bytes(code int, data []byte) error
}
```

### Server - 服务器接口

```go
type Server interface {
    Start(addr string) error
    Stop(ctx context.Context) error
    Serve(listener net.Listener) error
    Use(middleware ...MiddlewareFunc)
    Handle(pattern string, handler HandlerFunc)
}
```

## 使用示例

```go
package main

import (
    "github.com/whosafe/uf/ucontext"
    "github.com/whosafe/uf/uprotocol/unet"
    "github.com/whosafe/uf/uprotocol/ubind"
)

type User struct {
    ID   int
    Name string
}

func (u *User) Bind(key string, value *ubind.Value) error {
    switch key {
    case "id":
        u.ID = value.Int()
    case "name":
        u.Name = value.Str()
    }
    return nil
}

// 处理器函数 - 协议无关
func CreateUser(ctx *ucontext.Context, req unet.Request) error {
    var user User
    if err := req.Bind(&user); err != nil {
        return err
    }
    
    // 可以访问标准 context
    stdCtx := ctx.Context()
    
    // 可以访问追踪信息
    trace := ctx.Trace()
    trace.SetMetadata("user_id", strconv.Itoa(user.ID))
    
    // 业务逻辑
    saveUser(stdCtx, &user)
    
    return req.Response().JSON(200, user)
}
```

## 协议无缝切换

同一个 `CreateUser` 处理器可以在不同协议中使用:

```go
// HTTP
httpServer.POST("/users", CreateUser)

// TCP (未来)
tcpServer.Handle(MSG_CREATE_USER, CreateUser)

// QUIC (未来)
quicServer.Handle("/users", CreateUser)
```

## 优势

1. **协议无关** - 业务逻辑不依赖具体协议
2. **职责分离** - context 用于追踪,Request 用于数据
3. **易于测试** - 接口清晰,便于 Mock
4. **可扩展** - 新协议只需实现接口
