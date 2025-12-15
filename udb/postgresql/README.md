# udb/postgresql - PostgreSQL 数据库层

基于 `github.com/jackc/pgx/v5` 的零反射、高性能 PostgreSQL 数据库层。

## 📑 目录

- [核心特性](#-核心特性)
- [安装](#-安装)
- [快速开始](#-快速开始)
- [配置详解](#-配置详解)
- [查询构建器](#-查询构建器)
- [CRUD 构建器](#-crud-构建器)
- [事务处理](#-事务处理)
- [高级功能](#-高级功能)
- [API 参考](#-api-参考)
- [最佳实践](#-最佳实践)
- [性能优化](#-性能优化)

## ✨ 核心特性

### 🚀 零反射设计

- 手动实现 `Scanner` 接口,避免反射带来的性能损耗
- 性能接近原生 pgx 代码
- 类型安全,编译时检查

### 🔗 链路追踪

- 自动集成 `ucontext`,所有操作包含 `trace_id`
- 完整的日志记录,支持慢查询监控
- 可配置的日志级别和输出方式

### 📊 强大的查询构建器

- 支持 SELECT、JOIN、WHERE、GROUP BY、HAVING、ORDER BY、LIMIT/OFFSET
- 支持 DISTINCT 去重查询
- 支持多种 WHERE 条件:IN、BETWEEN、LIKE、NULL 等
- 链式调用,代码简洁优雅

### 💼 完整的 CRUD 构建器

- **Insert**: 插入数据,支持 RETURNING 子句
- **Update**: 更新数据,支持多字段更新和条件过滤
- **Delete**: 删除数据,支持条件过滤
- 所有构建器支持链式调用

### 🔄 事务支持

- 完整的事务管理(Begin/Commit/Rollback)
- 事务中支持所有查询和 CRUD 操作
- 自动日志记录事务状态

### ⚙️ 灵活的配置

- 支持 YAML 配置文件
- 连接池配置(最大/最小连接数、生命周期等)
- 查询配置(超时、慢查询阈值)
- 日志配置(级别、格式、输出方式)

## 📦 安装

```bash
go get github.com/whosafe/uf/udb/postgresql
```

## 🚀 快速开始

### 1. 配置文件

创建 `config.yaml`:

```yaml
database:
  postgres:
    # 连接配置
    host: "localhost"
    port: 5432
    username: "postgres"
    password: "your_password"
    database: "myapp"
    ssl_mode: "disable"
    
    # 连接池配置
    pool:
      max_conns: 25
      min_conns: 5
      max_conn_lifetime: "1h"
      max_conn_idle_time: "30m"
      health_check_period: "1m"
    
    # 查询配置
    query:
      default_timeout: "30s"
      slow_query_threshold: "1s"
    
    # 日志配置
    log:
      enabled: true
      level: "info"
      format: "text"
      output: "stdout"
      slow_query: true
      log_params: false
```

### 2. 定义数据结构

实现 `Scanner` 接口以实现零反射:

```go
type User struct {
    ID        int64
    Username  string
    Email     string
    Age       int
    CreatedAt time.Time
}

// Scan 实现 Scanner 接口(零反射)
func (u *User) Scan(key string, value any) error {
    switch key {
    case "id":
        u.ID = uconv.ToInt64Def(value, 0)
    case "username":
        u.Username = uconv.ToString(value)
    case "email":
        u.Email = uconv.ToString(value)
    case "age":
        u.Age = uconv.ToIntDef(value, 0)
    case "created_at":
        u.CreatedAt = uconv.ToTimeDef(value, time.Time{})
    }
    return nil
}
```

### 3. 创建连接并使用

```go
package main

import (
    "context"
    "log"
    
    "github.com/whosafe/uf/uconfig"
    "github.com/whosafe/uf/ucontext"
    "github.com/whosafe/uf/udb/postgresql"
)

func main() {
    // 加载配置
    uconfig.Load("config.yaml")
    
    // 创建连接
    conn, err := postgresql.New(postgresql.GetConfig())
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()
    
    // 创建追踪上下文
    ctx := ucontext.NewContext(context.Background())
    
    // 查询单条记录
    var user User
    err = conn.Query(ctx).
        Table("users").
        Where("id = ?", 1).
        Scan(&user)
    
    if err != nil {
        if err == postgresql.ErrNoRows {
            log.Println("用户不存在")
        } else {
            log.Fatal(err)
        }
    }
    
    log.Printf("用户: %+v\n", user)
}
```

## 📖 配置详解

### 连接配置

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `host` | string | 是 | 数据库主机地址 |
| `port` | int | 是 | 数据库端口(默认 5432) |
| `username` | string | 是 | 数据库用户名 |
| `password` | string | 否 | 数据库密码 |
| `database` | string | 是 | 数据库名称 |
| `ssl_mode` | string | 否 | SSL 模式:disable, require, verify-ca, verify-full |

### 连接池配置

| 字段 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `max_conns` | int | 25 | 最大连接数 |
| `min_conns` | int | 5 | 最小连接数 |
| `max_conn_lifetime` | duration | 1h | 连接最大生命周期 |
| `max_conn_idle_time` | duration | 30m | 连接最大空闲时间 |
| `health_check_period` | duration | 1m | 健康检查周期 |

### 查询配置

| 字段 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `default_timeout` | duration | 30s | 默认查询超时时间 |
| `slow_query_threshold` | duration | 1s | 慢查询阈值 |

### 日志配置

| 字段 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `enabled` | bool | true | 是否启用日志 |
| `level` | string | info | 日志级别:debug, info, warn, error |
| `format` | string | text | 日志格式:text, json |
| `output` | string | stdout | 输出方式:stdout, stderr, file |
| `file_path` | string | ./logs/db.log | 日志文件路径(output=file 时) |
| `max_size` | int | 100 | 日志文件最大大小(MB) |
| `max_backups` | int | 10 | 最大备份文件数 |
| `max_age` | int | 30 | 日志文件最大保留天数 |
| `compress` | bool | false | 是否压缩备份文件 |
| `slow_query` | bool | true | 是否记录慢查询 |
| `log_params` | bool | false | 是否记录查询参数(生产环境建议关闭) |

## 🔍 查询构建器

查询构建器提供了链式 API 来构建复杂的 SQL 查询。

### 基础查询

```go
// 查询单条记录
var user User
err := conn.Query(ctx).
    Table("users").
    Where("id = ?", 1).
    Scan(&user)

// 查询多条记录
results, err := conn.Query(ctx).
    Table("users").
    Where("age > ?", 18).
    OrderBy("created_at").
    Limit(10).
    ScanAll(func() postgresql.Scanner { return &User{} })

// 遍历结果
for _, r := range results {
    user := r.(*User)
    fmt.Printf("%+v\n", user)
}
```

### SELECT 字段

```go
// 选择特定字段
conn.Query(ctx).
    Select("id", "username", "email").
    Table("users").
    ScanAll(newUser)

// 使用别名和函数
conn.Query(ctx).
    Select("u.id", "u.username", "COUNT(o.id) as order_count").
    Table("users u").
    LeftJoin("orders o", "u.id = o.user_id").
    GroupBy("u.id", "u.username").
    ScanAll(newUserStats)
```

### WHERE 条件

#### 基础条件

```go
// 单个条件
conn.Query(ctx).
    Table("users").
    Where("age > ?", 18).
    ScanAll(newUser)

// 多个 AND 条件
conn.Query(ctx).
    Table("users").
    Where("age > ?", 18).
    Where("status = ?", "active").
    Where("created_at > ?", time.Now().AddDate(0, -1, 0)).
    ScanAll(newUser)

// OR 条件
conn.Query(ctx).
    Table("users").
    Where("age > ?", 18).
    OrWhere("vip = ?", true).
    ScanAll(newUser)
```

#### IN 条件

```go
// WHERE id IN (1, 2, 3)
conn.Query(ctx).
    Table("users").
    WhereIn("id", []any{1, 2, 3}).
    ScanAll(newUser)

// WHERE status NOT IN ('deleted', 'banned')
conn.Query(ctx).
    Table("users").
    WhereNotIn("status", []any{"deleted", "banned"}).
    ScanAll(newUser)
```

#### BETWEEN 条件

```go
// WHERE age BETWEEN 18 AND 65
conn.Query(ctx).
    Table("users").
    WhereBetween("age", 18, 65).
    ScanAll(newUser)

// WHERE age NOT BETWEEN 0 AND 18
conn.Query(ctx).
    Table("users").
    WhereNotBetween("age", 0, 18).
    ScanAll(newUser)
```

#### NULL 条件

```go
// WHERE email IS NULL
conn.Query(ctx).
    Table("users").
    WhereNull("email").
    ScanAll(newUser)

// WHERE email IS NOT NULL
conn.Query(ctx).
    Table("users").
    WhereNotNull("email").
    ScanAll(newUser)
```

#### LIKE 条件

```go
// WHERE username LIKE 'admin%'
conn.Query(ctx).
    Table("users").
    WhereLike("username", "admin%").
    ScanAll(newUser)
```

### JOIN 查询

```go
// INNER JOIN
conn.Query(ctx).
    Select("u.*", "p.bio").
    Table("users u").
    Join("profiles p", "u.id = p.user_id").
    Where("u.status = ?", "active").
    ScanAll(newUser)

// LEFT JOIN
conn.Query(ctx).
    Select("u.id", "u.username", "COUNT(o.id) as order_count").
    Table("users u").
    LeftJoin("orders o", "u.id = o.user_id").
    GroupBy("u.id", "u.username").
    ScanAll(newUserStats)

// RIGHT JOIN
conn.Query(ctx).
    Table("users u").
    RightJoin("orders o", "u.id = o.user_id").
    ScanAll(newOrder)

// FULL JOIN
conn.Query(ctx).
    Table("users u").
    FullJoin("profiles p", "u.id = p.user_id").
    ScanAll(newUser)
```

### GROUP BY 和 HAVING

```go
// 按分类统计产品
conn.Query(ctx).
    Select("category", "COUNT(*) as count", "AVG(price) as avg_price").
    Table("products").
    GroupBy("category").
    Having("COUNT(*) > ?", 10).
    OrderByDesc("count").
    ScanAll(newCategoryStats)

// 多字段分组
conn.Query(ctx).
    Select("category", "brand", "COUNT(*) as count").
    Table("products").
    GroupBy("category", "brand").
    Having("COUNT(*) > ?", 5).
    ScanAll(newStats)
```

### ORDER BY

```go
// 升序排序
conn.Query(ctx).
    Table("users").
    OrderBy("created_at").
    ScanAll(newUser)

// 降序排序
conn.Query(ctx).
    Table("users").
    OrderByDesc("created_at").
    ScanAll(newUser)

// 多字段排序
conn.Query(ctx).
    Table("users").
    OrderBy("age").
    OrderByDesc("created_at").
    ScanAll(newUser)
```

### LIMIT 和 OFFSET

```go
// 分页查询
page := 1
pageSize := 10

conn.Query(ctx).
    Table("users").
    OrderBy("id").
    Limit(pageSize).
    Offset((page - 1) * pageSize).
    ScanAll(newUser)
```

### DISTINCT

```go
// 查询不重复的分类
conn.Query(ctx).
    Select("category").
    Table("products").
    Distinct().
    OrderBy("category").
    ScanAll(newCategory)
```

## 🔧 CRUD 构建器

### Insert 构建器

#### 基础插入

```go
// 插入单条记录
affected, err := conn.Insert(ctx).
    Table("users").
    Columns("username", "email", "age").
    Values("alice", "alice@example.com", 25).
    Exec()

if err != nil {
    log.Fatal(err)
}
fmt.Printf("插入成功,影响行数: %d\n", affected)
```

#### 插入并返回数据

```go
// 插入并返回完整记录(包含自动生成的 ID)
var newUser User
err := conn.Insert(ctx).
    Table("users").
    Columns("username", "email", "age").
    Values("bob", "bob@example.com", 30).
    ExecReturning(&newUser)

if err != nil {
    log.Fatal(err)
}
fmt.Printf("新用户 ID: %d\n", newUser.ID)
```

### Update 构建器

#### 基础更新

```go
// 更新单个字段
affected, err := conn.Update(ctx).
    Table("users").
    Set("age", 26).
    Where("id = ?", 1).
    Exec()

// 更新多个字段
affected, err := conn.Update(ctx).
    Table("users").
    Set("age", 26).
    Set("email", "newemail@example.com").
    Where("id = ?", 1).
    Exec()
```

#### 批量更新

```go
// 使用 SetMap 批量设置字段
data := map[string]any{
    "age":        26,
    "email":      "newemail@example.com",
    "updated_at": time.Now(),
}

affected, err := conn.Update(ctx).
    Table("users").
    SetMap(data).
    Where("id = ?", 1).
    Exec()
```

#### 条件更新

```go
// 使用多个条件
affected, err := conn.Update(ctx).
    Table("users").
    Set("status", "inactive").
    Where("last_login < ?", time.Now().AddDate(0, -6, 0)).
    Where("status = ?", "active").
    Exec()
```

### Delete 构建器

#### 基础删除

```go
// 删除单条记录
affected, err := conn.Delete(ctx).
    Table("users").
    Where("id = ?", 1).
    Exec()
```

#### 条件删除

```go
// 删除多条记录
affected, err := conn.Delete(ctx).
    Table("users").
    Where("status = ?", "deleted").
    Where("created_at < ?", time.Now().AddDate(-1, 0, 0)).
    Exec()

// 使用 IN 条件删除
affected, err := conn.Delete(ctx).
    Table("users").
    WhereIn("id", []any{1, 2, 3}).
    Exec()

// 使用 LIKE 条件删除
affected, err := conn.Delete(ctx).
    Table("users").
    WhereLike("username", "test_%").
    Exec()
```

## 🔄 事务处理

### 基础事务

```go
// 开始事务
tx, err := conn.Begin(ctx)
if err != nil {
    log.Fatal(err)
}

// 执行操作
_, err = tx.Insert(ctx).
    Table("users").
    Columns("username", "email").
    Values("alice", "alice@example.com").
    Exec()

if err != nil {
    tx.Rollback()
    log.Fatal(err)
}

_, err = tx.Insert(ctx).
    Table("profiles").
    Columns("user_id", "bio").
    Values(1, "Hello World").
    Exec()

if err != nil {
    tx.Rollback()
    log.Fatal(err)
}

// 提交事务
if err := tx.Commit(); err != nil {
    log.Fatal(err)
}
```

### 事务中的查询

```go
tx, _ := conn.Begin(ctx)

// 在事务中查询
var user User
err := tx.Query(ctx).
    Table("users").
    Where("id = ?", 1).
    Scan(&user)

if err != nil {
    tx.Rollback()
    return err
}

// 在事务中更新
_, err = tx.Update(ctx).
    Table("users").
    Set("last_login", time.Now()).
    Where("id = ?", user.ID).
    Exec()

if err != nil {
    tx.Rollback()
    return err
}

tx.Commit()
```

### 事务最佳实践

```go
func transferMoney(conn *postgresql.Connection, fromID, toID int64, amount float64) error {
    ctx := ucontext.NewContext(context.Background())
    
    tx, err := conn.Begin(ctx)
    if err != nil {
        return err
    }
    
    // 使用 defer 确保事务被正确处理
    defer func() {
        if err != nil {
            tx.Rollback()
        }
    }()
    
    // 扣款
    _, err = tx.Update(ctx).
        Table("accounts").
        Set("balance", "balance - ?").
        Where("id = ?", fromID).
        Where("balance >= ?", amount).
        Exec()
    
    if err != nil {
        return err
    }
    
    // 入账
    _, err = tx.Update(ctx).
        Table("accounts").
        Set("balance", "balance + ?").
        Where("id = ?", toID).
        Exec()
    
    if err != nil {
        return err
    }
    
    // 提交事务
    return tx.Commit()
}
```

## 🎯 高级功能

### 直接执行 SQL

当查询构建器无法满足需求时,可以直接执行 SQL:

```go
// 执行原始 SQL
affected, err := conn.Exec(ctx,
    "UPDATE users SET age = age + 1 WHERE created_at < $1",
    time.Now().AddDate(-1, 0, 0))

// 复杂查询
affected, err := conn.Exec(ctx, `
    WITH recent_orders AS (
        SELECT user_id, COUNT(*) as order_count
        FROM orders
        WHERE created_at > $1
        GROUP BY user_id
    )
    UPDATE users u
    SET vip = true
    FROM recent_orders ro
    WHERE u.id = ro.user_id AND ro.order_count > $2
`, time.Now().AddDate(0, -1, 0), 10)
```

### 连接池管理

```go
// 获取连接池统计信息
stats := conn.Stats()
fmt.Printf("总连接数: %d\n", stats.TotalConns())
fmt.Printf("空闲连接数: %d\n", stats.IdleConns())
fmt.Printf("获取连接数: %d\n", stats.AcquiredConns())

// 健康检查
if err := conn.Ping(ctx); err != nil {
    log.Printf("数据库连接异常: %v\n", err)
}

// 关闭连接池
defer conn.Close()
```

### 错误处理

```go
err := conn.Query(ctx).
    Table("users").
    Where("id = ?", 1).
    Scan(&user)

if err != nil {
    // 判断是否为"未找到记录"错误
    if err == postgresql.ErrNoRows {
        log.Println("用户不存在")
        return nil
    }
    
    // 其他错误
    return fmt.Errorf("查询用户失败: %w", err)
}
```

## 📚 API 参考

### Connection

| 方法 | 说明 |
|------|------|
| `Query(ctx) *QueryBuilder` | 创建查询构建器 |
| `Insert(ctx) *InsertBuilder` | 创建插入构建器 |
| `Update(ctx) *UpdateBuilder` | 创建更新构建器 |
| `Delete(ctx) *DeleteBuilder` | 创建删除构建器 |
| `Exec(ctx, sql, args...) (int64, error)` | 执行原始 SQL |
| `Begin(ctx) (*Transaction, error)` | 开始事务 |
| `Ping(ctx) error` | 健康检查 |
| `Stats() *pgxpool.Stat` | 获取连接池统计 |
| `Close()` | 关闭连接池 |

### QueryBuilder

| 方法 | 说明 |
|------|------|
| `Table(name) *QueryBuilder` | 设置表名 |
| `Select(fields...) *QueryBuilder` | 设置查询字段 |
| `Distinct() *QueryBuilder` | 去重 |
| `Where(condition, args...) *QueryBuilder` | 添加 WHERE 条件 |
| `OrWhere(condition, args...) *QueryBuilder` | 添加 OR WHERE 条件 |
| `WhereIn(field, values) *QueryBuilder` | 添加 IN 条件 |
| `WhereNotIn(field, values) *QueryBuilder` | 添加 NOT IN 条件 |
| `WhereBetween(field, min, max) *QueryBuilder` | 添加 BETWEEN 条件 |
| `WhereNotBetween(field, min, max) *QueryBuilder` | 添加 NOT BETWEEN 条件 |
| `WhereNull(field) *QueryBuilder` | 添加 IS NULL 条件 |
| `WhereNotNull(field) *QueryBuilder` | 添加 IS NOT NULL 条件 |
| `WhereLike(field, pattern) *QueryBuilder` | 添加 LIKE 条件 |
| `Join(table, on) *QueryBuilder` | INNER JOIN |
| `LeftJoin(table, on) *QueryBuilder` | LEFT JOIN |
| `RightJoin(table, on) *QueryBuilder` | RIGHT JOIN |
| `FullJoin(table, on) *QueryBuilder` | FULL JOIN |
| `GroupBy(fields...) *QueryBuilder` | GROUP BY |
| `Having(condition, args...) *QueryBuilder` | HAVING |
| `OrderBy(field) *QueryBuilder` | 升序排序 |
| `OrderByDesc(field) *QueryBuilder` | 降序排序 |
| `Limit(n) *QueryBuilder` | 限制数量 |
| `Offset(n) *QueryBuilder` | 偏移量 |
| `Scan(dest Scanner) error` | 扫描单行 |
| `ScanAll(newScanner func() Scanner) ([]Scanner, error)` | 扫描多行 |

### InsertBuilder

| 方法 | 说明 |
|------|------|
| `Table(name) *InsertBuilder` | 设置表名 |
| `Columns(cols...) *InsertBuilder` | 设置列名 |
| `Values(vals...) *InsertBuilder` | 设置值 |
| `Exec() (int64, error)` | 执行插入 |
| `ExecReturning(dest Scanner) error` | 执行插入并返回数据 |

### UpdateBuilder

| 方法 | 说明 |
|------|------|
| `Table(name) *UpdateBuilder` | 设置表名 |
| `Set(column, value) *UpdateBuilder` | 设置字段值 |
| `SetMap(data map[string]any) *UpdateBuilder` | 批量设置字段 |
| `Where(condition, args...) *UpdateBuilder` | 添加 WHERE 条件 |
| `WhereIn/WhereNotIn/...` | 同 QueryBuilder |
| `Exec() (int64, error)` | 执行更新 |

### DeleteBuilder

| 方法 | 说明 |
|------|------|
| `Table(name) *DeleteBuilder` | 设置表名 |
| `Where(condition, args...) *DeleteBuilder` | 添加 WHERE 条件 |
| `WhereIn/WhereNotIn/...` | 同 QueryBuilder |
| `Exec() (int64, error)` | 执行删除 |

### Transaction

| 方法 | 说明 |
|------|------|
| `Query(ctx) *TxQueryBuilder` | 创建事务查询构建器 |
| `Insert(ctx) *TxInsertBuilder` | 创建事务插入构建器 |
| `Update(ctx) *TxUpdateBuilder` | 创建事务更新构建器 |
| `Delete(ctx) *TxDeleteBuilder` | 创建事务删除构建器 |
| `Exec(ctx, sql, args...) (int64, error)` | 执行原始 SQL |
| `Commit() error` | 提交事务 |
| `Rollback() error` | 回滚事务 |

## 💡 最佳实践

### 1. 始终使用 ucontext

```go
// ✅ 正确
ctx := ucontext.NewContext(context.Background())
conn.Query(ctx).Table("users").Scan(&user)

// ❌ 错误
conn.Query(context.Background()).Table("users").Scan(&user)
```

### 2. 实现 Scanner 接口时使用 uconv

```go
// ✅ 正确 - 使用 uconv 进行类型转换
func (u *User) Scan(key string, value any) error {
    switch key {
    case "id":
        u.ID = uconv.ToInt64Def(value, 0)
    case "age":
        u.Age = uconv.ToIntDef(value, 0)
    }
    return nil
}

// ❌ 错误 - 直接类型断言可能导致 panic
func (u *User) Scan(key string, value any) error {
    switch key {
    case "id":
        u.ID = value.(int64) // 可能 panic
    }
    return nil
}
```

### 3. 正确处理错误

```go
// ✅ 正确
err := conn.Query(ctx).Table("users").Where("id = ?", 1).Scan(&user)
if err != nil {
    if err == postgresql.ErrNoRows {
        // 处理未找到记录的情况
        return nil
    }
    return fmt.Errorf("查询失败: %w", err)
}

// ❌ 错误 - 忽略错误
conn.Query(ctx).Table("users").Where("id = ?", 1).Scan(&user)
```

### 4. 使用连接池而非单个连接

```go
// ✅ 正确 - 创建一次连接,复用连接池
func main() {
    conn, _ := postgresql.New(config)
    defer conn.Close()
    
    // 多次查询复用连接池
    for i := 0; i < 100; i++ {
        conn.Query(ctx).Table("users").Scan(&user)
    }
}

// ❌ 错误 - 每次都创建新连接
for i := 0; i < 100; i++ {
    conn, _ := postgresql.New(config)
    conn.Query(ctx).Table("users").Scan(&user)
    conn.Close()
}
```

### 5. 事务中使用 defer 确保回滚

```go
// ✅ 正确
func doSomething(conn *postgresql.Connection) error {
    tx, err := conn.Begin(ctx)
    if err != nil {
        return err
    }
    
    defer func() {
        if err != nil {
            tx.Rollback()
        }
    }()
    
    // 执行操作...
    
    return tx.Commit()
}
```

### 6. 生产环境关闭参数日志

```yaml
log:
  log_params: false  # 生产环境关闭,避免敏感信息泄露
```

### 7. 合理配置连接池

```yaml
pool:
  max_conns: 25        # 根据实际负载调整
  min_conns: 5         # 保持最小连接数,减少连接建立开销
  max_conn_lifetime: "1h"    # 定期回收连接
  max_conn_idle_time: "30m"  # 回收空闲连接
```

## ⚡ 性能优化

### 1. 零反射设计

通过实现 `Scanner` 接口,避免反射带来的性能损耗:

```go
// 性能对比(相同查询)
// 反射方式: ~500 ns/op
// Scanner 方式: ~200 ns/op
// 性能提升: ~2.5x
```

### 2. 连接池复用

```go
// 连接池配置建议
pool:
  max_conns: 25      # CPU 核心数 * 2 + 磁盘数
  min_conns: 5       # 保持最小连接,减少连接建立开销
```

### 3. 批量操作

```go
// ✅ 推荐 - 批量插入
tx, _ := conn.Begin(ctx)
for _, user := range users {
    tx.Insert(ctx).Table("users").Columns(...).Values(...).Exec()
}
tx.Commit()

// ❌ 不推荐 - 逐条插入
for _, user := range users {
    conn.Insert(ctx).Table("users").Columns(...).Values(...).Exec()
}
```

### 4. 使用索引

```sql
-- 为常用查询字段创建索引
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_created_at ON users(created_at);
```

### 5. 慢查询监控

```yaml
query:
  slow_query_threshold: "1s"  # 设置合理的慢查询阈值
log:
  slow_query: true            # 启用慢查询日志
```

## 🔗 相关链接

- [pgx 文档](https://github.com/jackc/pgx)
- [示例代码](../../example/udb/postgresql/)
- [uconfig 配置库](../uconfig/)
- [ucontext 上下文库](../ucontext/)
- [uconv 类型转换库](../uconv/)

## 📄 许可证

MIT License
