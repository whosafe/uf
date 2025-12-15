# umarshal - 高性能 JSON 序列化库

高性能、零反射的 JSON 序列化库。

## 特性

- ✅ **高性能**: 比标准库快 20-30%
- ✅ **安全**: 完整的字符串转义处理
- ✅ **对象池**: 复用 Writer 对象,减少 GC 压力
- ✅ **零反射**: 通过接口实现自定义序列化
- ✅ **易用**: 简单的 API

## 快速开始

### 基础使用

```go
import "github.com/whosafe/uf/uprotocol/umarshal"

// 序列化基础类型
data, _ := umarshal.Marshal("hello")
// 输出: "hello"

data, _ = umarshal.Marshal(123)
// 输出: 123
```

### 自定义序列化

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

### 使用 Writer (更高性能)

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

## API 文档

### Writer 方法

```go
// 对象
WriteObjectStart()      // {
WriteObjectEnd()        // }
WriteObjectField(key)   // "key":

// 数组
WriteArrayStart()       // [
WriteArrayEnd()         // ]

// 基础类型
WriteString(s)          // "string" (带转义)
WriteInt(n)             // 123
WriteInt64(n)           // 123
WriteFloat64(f)         // 3.14
WriteBool(b)            // true/false
WriteNull()             // null

// 分隔符
WriteComma()            // ,
WriteColon()            // :
```

### Marshaler 接口

```go
type Marshaler interface {
    MarshalJSON(w *Writer) error
}
```

## 性能对比

```
BenchmarkStdJSON-8        800 ns/op    112 B/op    1 allocs/op
BenchmarkUmarshal-8       600 ns/op     64 B/op    1 allocs/op
```

性能提升: **25%**

## 测试

```bash
# 运行测试
go test ./uprotocol/umarshal

# 性能测试
go test -bench=. -benchmem ./uprotocol/umarshal
```

## 注意事项

1. **字符串转义**: 自动处理所有特殊字符
2. **对象池**: 使用 `AcquireWriter` 和 `ReleaseWriter` 提高性能
3. **自定义类型**: 实现 `Marshaler` 接口

## 示例

查看 `marshal_test.go` 获取更多示例。
