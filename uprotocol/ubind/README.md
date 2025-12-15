# ubind - 数据绑定包

零反射的数据绑定解决方案,用于协议无关的请求数据解析。

## 特性

- ✅ **零反射** - 手动实现,高性能
- ✅ **类型安全** - 编译时检查
- ✅ **协议无关** - 支持 HTTP、TCP、QUIC 等多种协议
- ✅ **自动识别** - 自动识别 JSON/Form/Binary 格式
- ✅ **支持嵌套** - 完整支持嵌套对象和数组

## 快速开始

### 简单对象

```go
package main

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

### 嵌套对象

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

### 数组

```go
type Team struct {
    Name    string
    Members []User
}

func (t *Team) Bind(key string, value *ubind.Value) error {
    switch key {
    case "name":
        t.Name = value.Str()
    case "members":
        if value.IsArray() {
            t.Members = make([]User, value.Len())
            for i := 0; i < value.Len(); i++ {
                ubind.Bind(value.Index(i), &t.Members[i])
            }
        }
    }
    return nil
}
```

### Form 数据解析

```go
// Form 数据: name=Alice&age=25&email=test%40example.com
formData := []byte("name=Alice&age=25&email=test%40example.com")

val := ubind.Parse(formData)  // 自动识别 Form 格式
// 自动 URL 解码: email = "test@example.com"

var user User
ubind.Bind(val, &user)
```

## 在协议中使用

所有协议使用相同的 `Bind` 方法:

```go
// HTTP
func (c *HTTPContext) Bind(obj ubind.Binder) error {
    val := ubind.Parse(c.Request.Body)
    return ubind.Bind(val, obj)
}

// TCP
func (c *TCPContext) Bind(obj ubind.Binder) error {
    val := ubind.Parse(c.Message.Data)
    return ubind.Bind(val, obj)
}

// QUIC
func (c *QUICContext) Bind(obj ubind.Binder) error {
    val := ubind.Parse(c.Stream.Data)
    return ubind.Bind(val, obj)
}
```

## API

### Value 类型

```go
type Value struct {
    Type   ValueType
    Bool   bool
    Number float64
    String string
    Array  []*Value
    Object map[string]*Value
}
```

### 辅助方法

- `value.Int()` - 获取整数值
- `value.Float()` - 获取浮点数值
- `value.Str()` - 获取字符串值
- `value.Get(key)` - 获取对象字段
- `value.Index(i)` - 获取数组元素
- `value.Len()` - 获取数组长度
- `value.IsObject()` / `IsArray()` / `IsString()` 等 - 类型判断

### Binder 接口

```go
type Binder interface {
    Bind(key string, value *Value) error
}
```

### 函数

- `Parse(data []byte) *Value` - 自动识别格式并解析 (JSON/Form)
- `ParseJSON(data []byte) *Value` - 强制 JSON 解析
- `ParseForm(data []byte) *Value` - 强制 Form 解析
- `Bind(val *Value, v Binder) error` - 绑定数据到对象

## 测试

```bash
go test -v
```

所有测试通过 ✅

## 优势

1. **零反射** - 完全手动实现,性能高
2. **类型安全** - 编译时检查,避免运行时错误
3. **协议无关** - 可用于任何协议的数据解析
4. **易于使用** - 简洁的 API,链式调用
5. **完整支持** - 支持嵌套对象、数组、所有 JSON 类型

## 包结构

```text
ubind/
├── value.go         # Value 结构体定义
├── binder.go        # Binder 接口
├── bind.go          # Bind 函数
├── parse.go         # Parse 自动识别
├── parse_json.go    # JSON 解析器
├── parse_form.go    # Form 解析器 (application/x-www-form-urlencoded)
├── parse_binary.go  # Binary 解析器 (TODO)
└── README.md        # 文档
```

## 支持的格式

- ✅ **JSON** - 完整支持对象、数组、字符串、数字、布尔值、null
- ✅ **Form** - application/x-www-form-urlencoded,支持 URL 解码
- ⏳ **Binary** - 自定义二进制格式 (计划中)
