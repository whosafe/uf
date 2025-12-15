# uvalidator - 验证器包

基于结构化规则的高性能验证器，零反射。

## 包结构

``` sh
uvalidator/
├── validator.go      # Validator 接口和 Rule 接口
├── errors.go         # 错误类型
├── i18n.go          # 国际化支持
├── README.md
├── i18n/            # 国际化消息（按语言划分）
│   ├── en.go        # 英文消息
│   ├── zh.go        # 中文消息
│   └── helper.go    # 辅助函数
└── rule/            # 验证规则
    ├── required.go  # Required 规则结构体
    ├── min.go       # Min 规则结构体
    └── email.go     # Email 规则结构体
```

## 设计理念

### 规则结构

每个验证规则都是一个独立的结构体，实现 `Rule` 接口：

```go
type Rule interface {
    Validate(value any) bool
    GetMessage(field string, params map[string]string, lang ...Language) string
    Name() string
}
```

> **注意**: `GetMessage` 方法的 `lang` 参数是可选的。不传递时使用全局语言设置,传递时使用指定语言。

### 国际化结构（按语言划分）

```go
// i18n/en.go - 英文消息
var EN = map[string]string{
    "required": "{field} is required",
    "min": "{field} must be at least {param}",
    "email": "{field} must be a valid email address",
}

// i18n/zh.go - 中文消息
var ZH = map[string]string{
    "required": "{field}不能为空",
    "min": "{field}不能少于{param}",
    "email": "{field}必须是有效的邮箱地址",
}

// 规则中使用
import "github.com/whosafe/uf/uvalidator/i18n"

template := i18n.GetMessage("required")
```

## 使用示例

### 1. 定义结构体并实现 Validate 方法

```go
package model

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

### 2. 使用验证

```go
req := CreateUserRequest{
    Username: "ab",
    Email:    "invalid",
    Age:      15,
}

if err := req.Validate(); err != nil {
    fmt.Println(err)
}
```

### 3. 国际化

#### 方式 1: 全局语言设置

```go
// 设置全局语言为中文
uvalidator.SetLanguage(uvalidator.LanguageZH)

// 错误消息会显示中文
// "Username不能为空" 而不是 "Username is required"
```

#### 方式 2: 请求级语言选择 (推荐)

```go
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
    // 从 HTTP 请求头解析语言
    lang := uvalidator.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
    
    var req CreateUserRequest
    json.NewDecoder(r.Body).Decode(&req)
    
    // 验证时使用请求语言
    var errs uvalidator.ValidationErrors
    requiredRule := rule.NewRequired()
    if !requiredRule.Validate(req.Username) {
        errs = append(errs, uvalidator.NewFieldError(
            "Username",
            requiredRule.Name(),
            req.Username,
            requiredRule.GetMessage("Username", nil, lang), // 使用请求语言
        ))
    }
    
    if errs.HasErrors() {
        http.Error(w, errs.Error(), http.StatusBadRequest)
        return
    }
    
    // 处理请求...
}
```

**优势**:

- 支持多语言并发请求,互不干扰
- 自动解析 `Accept-Language` 请求头
- 向后兼容,不传语言参数时使用全局设置

## 添加新规则

### 1. 创建规则文件 `rule/phone.go`

```go
package rule

import (
    "github.com/whosafe/uf/uvalidator/i18n"
)

type Phone struct{}

func (p *Phone) Validate(value interface{}) bool {
    str, ok := value.(string)
    if !ok || str == "" {
        return true
    }
    return len(str) == 11 && str[0] == '1'
}

func (p *Phone) GetMessage(field string, lang ...uvalidator.Language) string  {
    template := i18n.GetMessage("phone", lang...)
    return replaceAll(template, "{field}", field)
}

func (p *Phone) Name() string {
    return "phone"
}

func NewPhone() *Phone {
    return &Phone{}
}
```

### 2. 添加国际化消息

在 `i18n/en.go` 中添加：

```go
var EN = map[string]string{
    // ... 其他消息
    "phone": "{field} must be a valid phone number",
}
```

在 `rule/i18n/zh.go` 中添加：

```go
var ZH = map[string]string{
    // ... 其他消息
    "phone": "{field}必须是有效的手机号",
}
```

## 添加新语言

创建 `rule/i18n/ja.go`（日语）：

```go
package i18n

var JA = map[string]string{
    "required": "{field}は必須です",
    "min": "{field}は{param}以上でなければなりません",
    "email": "{field}は有効なメールアドレスでなければなりません",
}
```

更新 `rule/i18n/helper.go`：

```go
func GetMessage(key string) string {
    lang := uvalidator.GetLanguage()
    
    var messages map[string]string
    switch lang {
    case uvalidator.LanguageZH:
        messages = ZH
    case uvalidator.LanguageJA:
        messages = JA
    default:
        messages = EN
    }
    
    // ...
}
```

## 优势

- ✅ **结构化**：每个规则都是独立的结构体
- ✅ **零反射**：手动实现验证逻辑
- ✅ **高性能**：接近手写代码
- ✅ **易扩展**：添加新规则只需创建新文件
- ✅ **国际化**：按语言组织，易于管理
- ✅ **类型安全**：编译时检查

## 内置规则

### 基础规则

- `Required` - 必填验证
- `Min` / `Max` - 最小值/最大值验证
- `Len` - 长度验证
- `Between` - 数值范围验证

### 比较规则

- `Gt` / `Gte` / `Lt` / `Lte` - 大于/大于等于/小于/小于等于验证

### 数字规则

- `Positive` / `Negative` - 正数/负数验证
- `Integer` - 整数验证
- `Decimal` - 小数验证

### 字符串规则

- `Email` - 邮箱验证
- `URL` - URL验证
- `Phone` - 手机号验证(中国)
- `Alpha` / `AlphaNum` / `Numeric` - 字母/字母数字/纯数字验证
- `Contains` / `StartsWith` / `EndsWith` - 字符串匹配
- `Regex` - 正则表达式验证
- `UUID` - UUID格式验证
- `JSON` - JSON格式验证
- `Base64` - Base64编码验证
- `Lowercase` / `Uppercase` - 大小写验证
- `ASCII` - ASCII字符验证
- `NotBlank` - 非空白验证

### 网络规则

- `IP` / `IPv4` / `IPv6` - IP地址验证
- `MAC` - MAC地址验证
- `Domain` - 域名验证
- `Port` - 端口号验证

### 日期时间规则

- `Date` / `DateTime` - 日期/日期时间格式验证
- `DateBefore` / `DateAfter` - 日期比较验证
- `DateBetween` - 日期范围验证

### 中国特色规则

- `IDCard` - 身份证号验证(15位或18位)
- `BankCard` - 银行卡号验证(Luhn算法)
- `UnifiedSocialCreditCode` - 统一社会信用代码验证
- `PostalCode` - 邮政编码验证
- `ChineseName` - 中文姓名验证

### 文件规则

- `FileExtension` - 文件扩展名验证
- `MimeType` - MIME类型验证
- `FileSize` - 文件大小验证

### 集合/数组规则

- `In` / `InInt` / `InFloat` - 值在指定列表中验证（支持字符串、整数、浮点数）
- `OneOf` / `OneOfInt` - 枚举值验证
- `Unique` - 数组元素唯一性验证
- `ArrayMin` / `ArrayMax` - 数组长度验证
- `ArrayContains` - 数组包含验证

### 安全规则

- `StrongPassword` - 强密码验证
- `NoHTML` - 不包含HTML标签验证
- `NoSQL` - 不包含SQL注入字符验证
- `NoXSS` - 不包含XSS攻击字符验证

### 其他实用规则

- `Confirmed` - 确认字段验证
- `Distinct` - 不同于指定值验证
- `NotIn` - 不在指定列表中验证
- `Nullable` - 允许为null验证

## 测试

包含完整的测试套件,覆盖所有验证规则:

```bash
# 运行所有测试
go test ./...

# 运行规则测试
go test ./rule/...

# 运行带详细输出的测试
go test -v ./rule/...
```

**测试覆盖**:

- 100+ 个测试用例
- 覆盖所有 60+ 个验证规则
- 包括边界测试和正负测试
- 表驱动测试模式

## 未来计划

- [ ] 字段关系验证规则（EqualField, RequiredIf 等）
- [ ] 代码生成工具自动生成 Validate 方法
- [ ] 规则组合和链式调用
- [ ] 自定义错误消息模板
