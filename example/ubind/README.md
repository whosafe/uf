# ubind 示例程序

演示 ubind 数据绑定包的各种功能。

## 运行

```bash
cd d:\code\go\src\uf\example\ubind
go run main.go
```

## 示例内容

1. **JSON 简单对象** - 解析基本的 JSON 对象
2. **JSON 嵌套对象** - 解析包含嵌套对象的 JSON
3. **JSON 数组** - 解析包含数组的 JSON
4. **Form 数据** - 解析 application/x-www-form-urlencoded 格式
5. **Form URL 编码** - 演示 URL 解码功能 (+ 转空格)
6. **自动格式识别** - Parse() 自动识别 JSON 或 Form

## 预期输出

```
=== ubind 数据绑定示例 ===

--- 示例 1: JSON 简单对象 ---
JSON: {"id":1,"name":"Alice","age":25}
解析结果: ID=1, Name=Alice, Age=25

--- 示例 2: JSON 嵌套对象 ---
JSON: {"id":2,"name":"Bob","address":{"city":"Beijing","street":"Main St"}}
解析结果: ID=2, Name=Bob, City=Beijing, Street=Main St

--- 示例 3: JSON 数组 ---
JSON: {"name":"Dev Team","members":[{"id":1,"name":"Alice","age":25},{"id":2,"name":"Bob","age":30}]}
解析结果: Team=Dev Team, Members=2
  Member 1: ID=1, Name=Alice, Age=25
  Member 2: ID=2, Name=Bob, Age=30

--- 示例 4: Form 数据 ---
Form: name=Charlie&age=28&id=3
解析结果: ID=3, Name=Charlie, Age=28

--- 示例 5: Form URL 编码 ---
Form: name=Hello+World&id=4&age=35
解析结果: ID=4, Name=Hello World, Age=35
(注意: 'Hello+World' 被解码为 'Hello World')

--- 示例 6: 自动格式识别 ---
Parse() 会自动识别 JSON 或 Form 格式:
  - JSON: 以 { 或 [ 开头
  - Form: 包含 = 字符

所有示例完成! ✓
```
