# uvalidator 使用示例

这个示例展示了 `uvalidator` 包的各种用法。

## 运行示例

```bash
cd d:\code\go\src\uf\example\uvalidator
go run main.go
```

## 示例内容

### 1. 基本验证 (全局语言)

演示如何使用全局语言设置进行验证:

```go
uvalidator.SetLanguage(uvalidator.LanguageZH)
if err := req.Validate(); err != nil {
    fmt.Println(err)
}
```

### 2. 有效数据验证

展示验证通过的情况。

### 3. 请求级语言选择

演示如何在同一个请求中使用不同的语言:

```go
// 使用英文
err := req.Validate(uvalidator.LanguageEN)

// 使用中文
err := req.Validate(uvalidator.LanguageZH)
```

### 4. 各种规则演示

展示以下规则的使用:

- 数字规则: `Between`, `Positive`
- 字符串规则: `UUID`, `Lowercase`
- 网络规则: `IP`
- 日期规则: `Date`
- 安全规则: `StrongPassword`, `NoHTML`
- 数组规则: `Unique`

### 5. HTTP 服务器示例

演示如何在 HTTP 处理器中使用请求级语言选择。

要启动 HTTP 服务器,取消 `main.go` 中最后一行的注释:

```go
http.ListenAndServe(":8080", nil)
```

然后测试:

```bash
# 使用中文
curl -X POST http://localhost:8080/users \
  -H 'Content-Type: application/json' \
  -H 'Accept-Language: zh-CN' \
  -d '{"username":"ab","email":"invalid"}'

# 使用英文
curl -X POST http://localhost:8080/users \
  -H 'Content-Type: application/json' \
  -H 'Accept-Language: en-US' \
  -d '{"username":"ab","email":"invalid"}'
```

## 输出示例

```
=== uvalidator 使用示例 ===

--- 示例 1: 基本验证 (全局语言) ---
验证失败:
Username长度不能少于3个字符; Email必须是有效的邮箱地址; Password不能为空; Age必须在18和100之间; Phone必须是有效的手机号

--- 示例 2: 有效数据 ---
验证通过! ✓

--- 示例 3: 请求级语言选择 ---
使用英文:
Username must be at least 3 characters; Email must be a valid email address; Password is required; Age must be between 18 and 100

使用中文:
Username长度不能少于3个字符; Email必须是有效的邮箱地址; Password不能为空; Age必须在18和100之间

--- 示例 4: 各种规则演示 ---
数字规则:
  Between(50): true
  Between(5): false
  Positive(10): true
  Positive(-5): false

字符串规则:
  UUID(valid): true
  UUID(invalid): false
  Lowercase(hello): true
  Lowercase(Hello): false

网络规则:
  IP(192.168.1.1): true
  IP(invalid): false

日期规则:
  Date(2024-01-01): true
  Date(2024/01/01): false

安全规则:
  StrongPassword(Abc123!@#): true
  StrongPassword(weak): false
  NoHTML(plain text): true
  NoHTML(<script>): false

数组规则:
  Unique([1,2,3]): true
  Unique([1,2,2]): false
```

## 关键特性

1. **请求级语言选择** - 支持在验证时动态指定语言,适合多语言并发场景
2. **全局语言设置** - 简单场景下可以使用全局语言设置
3. **丰富的验证规则** - 60+ 个内置验证规则
4. **零反射高性能** - 手写验证逻辑,性能接近原生代码
5. **类型安全** - 编译时检查,避免运行时错误
