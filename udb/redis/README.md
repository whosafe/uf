# udb/redis - Redis 客户端封装

基于 `github.com/redis/go-redis/v9` 的零反射、高性能 Redis 客户端封装。

## 📑 目录

- [核心特性](#-核心特性)
- [安装](#-安装)
- [快速开始](#-快速开始)
- [配置详解](#-配置详解)
- [基础操作](#-基础操作)
- [高级功能](#-高级功能)
- [API 参考](#-api-参考)
- [最佳实践](#-最佳实践)

## ✨ 核心特性

### 🚀 零反射设计

- 使用 `uconv` 进行类型安全转换
- 避免反射带来的性能损耗
- 类型安全，编译时检查

### 🔗 链路追踪

- 自动集成 `ucontext`，所有操作包含 `trace_id`
- 完整的日志记录，支持慢查询监控
- 可配置的日志级别和输出方式

### 📊 完整的 Redis 命令支持

- **字符串操作**: GET、SET、INCR、DECR 等
- **哈希操作**: HGET、HSET、HGETALL 等
- **列表操作**: LPUSH、RPUSH、LRANGE 等
- **集合操作**: SADD、SMEMBERS、SUNION 等
- **有序集合操作**: ZADD、ZRANGE、ZSCORE 等
- **键管理**: DEL、EXISTS、EXPIRE、TTL 等

### 🔄 高级功能

- **Pipeline**: 批量命令执行，提升性能
- **事务**: WATCH、MULTI、EXEC 支持
- **Pub/Sub**: 消息发布订阅

### ⚙️ 灵活的配置

- 支持 YAML 配置文件
- 连接池配置（大小、超时等）
- 查询配置（超时、慢查询阈值）
- 日志配置（级别、格式、输出方式）

## 📦 安装

```bash
go get github.com/whosafe/uf/udb/redis
go get github.com/redis/go-redis/v9
```

## 🚀 快速开始

### 1. 配置文件

创建 `config.yaml`:

```yaml
database:
  redis:
    # 连接配置
    host: "localhost"
    port: 6379
    password: ""
    db: 0
    
    # 连接池配置
    pool:
      pool_size: 10
      min_idle_conn: 5
      idle_timeout: "5m"
      max_lifetime: "1h"
    
    # 查询配置
    query:
      default_timeout: "30s"
      slow_query_threshold: "100ms"
    
    # 日志配置
    log:
      enabled: true
      level: "info"
      format: "text"
      output: "stdout"
      slow_query: true
      log_params: false
```

### 2. 基本使用

```go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/whosafe/uf/uconfig"
    "github.com/whosafe/uf/ucontext"
    "github.com/whosafe/uf/udb/redis"
)

func main() {
    // 加载配置
    uconfig.Load("config.yaml")
    
    // 创建连接
    conn, err := redis.New(redis.GetConfig())
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()
    
    // 创建追踪上下文
    ctx := ucontext.NewContext(context.Background())
    
    // 设置值
    err = conn.Set(ctx, "user:1:name", "Alice", 10*time.Minute)
    if err != nil {
        log.Fatal(err)
    }
    
    // 获取值
    name, err := conn.Get(ctx, "user:1:name")
    if err != nil {
        if err == redis.ErrNil {
            log.Println("键不存在")
        } else {
            log.Fatal(err)
        }
    }
    
    log.Printf("用户名: %s\n", name)
}
```

## 📖 配置详解

### 连接配置

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `host` | string | 是 | Redis 主机地址 |
| `port` | int | 是 | Redis 端口（默认 6379） |
| `password` | string | 否 | Redis 密码 |
| `db` | int | 否 | 数据库索引（默认 0） |

### 连接池配置

| 字段 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `pool_size` | int | 10 | 连接池大小 |
| `min_idle_conn` | int | 5 | 最小空闲连接数 |
| `idle_timeout` | duration | 5m | 空闲连接超时时间 |
| `max_lifetime` | duration | 1h | 连接最大生命周期 |

### 查询配置

| 字段 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `default_timeout` | duration | 30s | 默认查询超时时间 |
| `slow_query_threshold` | duration | 100ms | 慢查询阈值 |

### 日志配置

| 字段 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `enabled` | bool | true | 是否启用日志 |
| `level` | string | info | 日志级别：debug, info, warn, error |
| `format` | string | text | 日志格式：text, json |
| `output` | string | stdout | 输出方式：stdout, stderr, file |
| `slow_query` | bool | true | 是否记录慢查询 |
| `log_params` | bool | false | 是否记录查询参数 |

## 🔍 基础操作

### 字符串操作

```go
// 设置值
conn.Set(ctx, "key", "value", 0)

// 设置值（带过期时间）
conn.Set(ctx, "key", "value", 10*time.Minute)

// 仅在键不存在时设置
ok, _ := conn.SetNX(ctx, "key", "value", 0)

// 获取值
value, err := conn.Get(ctx, "key")
if err == redis.ErrNil {
    // 键不存在
}

// 批量获取
values, _ := conn.MGet(ctx, "key1", "key2", "key3")

// 批量设置
conn.MSet(ctx, "key1", "value1", "key2", "value2")

// 自增
count, _ := conn.Incr(ctx, "counter")

// 增加指定值
count, _ := conn.IncrBy(ctx, "counter", 10)

// 自减
count, _ := conn.Decr(ctx, "counter")

// 追加字符串
length, _ := conn.Append(ctx, "key", "suffix")
```

### 哈希操作

```go
// 设置字段
conn.HSet(ctx, "user:1", "name", "Alice", "age", 25)

// 获取字段
name, _ := conn.HGet(ctx, "user:1", "name")

// 获取所有字段
fields, _ := conn.HGetAll(ctx, "user:1")

// 删除字段
conn.HDel(ctx, "user:1", "age")

// 检查字段是否存在
exists, _ := conn.HExists(ctx, "user:1", "name")

// 获取字段数量
count, _ := conn.HLen(ctx, "user:1")

// 字段值自增
newAge, _ := conn.HIncrBy(ctx, "user:1", "age", 1)
```

### 列表操作

```go
// 从左侧插入
conn.LPush(ctx, "queue", "item1", "item2")

// 从右侧插入
conn.RPush(ctx, "queue", "item3", "item4")

// 从左侧弹出
item, _ := conn.LPop(ctx, "queue")

// 从右侧弹出
item, _ := conn.RPop(ctx, "queue")

// 获取范围内的元素
items, _ := conn.LRange(ctx, "queue", 0, -1)

// 获取列表长度
length, _ := conn.LLen(ctx, "queue")

// 获取指定索引的元素
item, _ := conn.LIndex(ctx, "queue", 0)

// 修剪列表
conn.LTrim(ctx, "queue", 0, 99)
```

### 集合操作

```go
// 添加成员
conn.SAdd(ctx, "tags", "go", "redis", "database")

// 移除成员
conn.SRem(ctx, "tags", "database")

// 获取所有成员
members, _ := conn.SMembers(ctx, "tags")

// 判断成员是否存在
exists, _ := conn.SIsMember(ctx, "tags", "go")

// 获取成员数量
count, _ := conn.SCard(ctx, "tags")

// 随机弹出成员
member, _ := conn.SPop(ctx, "tags")

// 集合运算
union, _ := conn.SUnion(ctx, "set1", "set2")
inter, _ := conn.SInter(ctx, "set1", "set2")
diff, _ := conn.SDiff(ctx, "set1", "set2")
```

### 有序集合操作

```go
import "github.com/redis/go-redis/v9"

// 添加成员
conn.ZAdd(ctx, "leaderboard",
    redis.Z{Score: 100, Member: "Alice"},
    redis.Z{Score: 95, Member: "Bob"},
    redis.Z{Score: 90, Member: "Charlie"})

// 移除成员
conn.ZRem(ctx, "leaderboard", "Charlie")

// 获取范围内的成员
members, _ := conn.ZRange(ctx, "leaderboard", 0, -1)

// 获取范围内的成员及分数
membersWithScores, _ := conn.ZRangeWithScores(ctx, "leaderboard", 0, -1)

// 倒序获取
members, _ := conn.ZRevRange(ctx, "leaderboard", 0, -1)

// 获取成员分数
score, _ := conn.ZScore(ctx, "leaderboard", "Alice")

// 获取成员排名
rank, _ := conn.ZRank(ctx, "leaderboard", "Alice")

// 增加成员分数
newScore, _ := conn.ZIncrBy(ctx, "leaderboard", 5, "Bob")

// 获取成员数量
count, _ := conn.ZCard(ctx, "leaderboard")
```

### 键管理

```go
// 删除键
conn.Del(ctx, "key1", "key2")

// 检查键是否存在
exists, _ := conn.Exists(ctx, "key")

// 设置过期时间
conn.Expire(ctx, "key", 10*time.Minute)

// 获取剩余生存时间
ttl, _ := conn.TTL(ctx, "key")

// 移除过期时间
conn.Persist(ctx, "key")

// 重命名键
conn.Rename(ctx, "oldkey", "newkey")

// 获取键的类型
keyType, _ := conn.Type(ctx, "key")

// 查找匹配的键
keys, _ := conn.Keys(ctx, "user:*")

// 迭代键
keys, cursor, _ := conn.Scan(ctx, 0, "user:*", 10)
```

## 🎯 高级功能

### Pipeline

Pipeline 可以批量执行命令，减少网络往返次数：

```go
// 方式1: 使用 Pipelined
cmds, err := conn.Pipelined(ctx, func(pipe redis.Pipeliner) error {
    pipe.Set(ctx, "key1", "value1", 0)
    pipe.Set(ctx, "key2", "value2", 0)
    pipe.Incr(ctx, "counter")
    return nil
})

// 方式2: 手动创建 Pipeline
pipe := conn.Pipeline()
setCmd := pipe.Set(ctx, "key1", "value1", 0)
incrCmd := pipe.Incr(ctx, "counter")
_, err := pipe.Exec(ctx)

// 获取结果
value := setCmd.Val()
count := incrCmd.Val()
```

### 事务

使用 WATCH 和事务 Pipeline 实现乐观锁：

```go
// 使用 Watch
err := conn.Watch(ctx, func(tx *redis.Tx) error {
    // 读取当前值
    val, err := tx.Get(ctx, "counter").Int64()
    if err != nil && err != redis.Nil {
        return err
    }
    
    // 在事务中更新
    _, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
        pipe.Set(ctx, "counter", val+1, 0)
        return nil
    })
    return err
}, "counter")

// 使用事务 Pipeline
pipe := conn.TxPipelineClient()
pipe.Set(ctx, "key1", "value1", 0)
pipe.Set(ctx, "key2", "value2", 0)
_, err := pipe.Exec(ctx)
```

### Pub/Sub

消息发布订阅：

```go
// 发布消息
subscribers, _ := conn.Publish(ctx, "news", "Hello, World!")

// 订阅频道
pubsub := conn.Subscribe(ctx, "news", "updates")
defer pubsub.Close()

// 接收消息
for {
    msg, err := pubsub.ReceiveMessage(ctx)
    if err != nil {
        break
    }
    fmt.Printf("收到消息: %s from %s\n", msg.Payload, msg.Channel)
}

// 模式订阅
pubsub := conn.PSubscribe(ctx, "news:*", "updates:*")
```

## 📚 API 参考

### Connection

| 方法 | 说明 |
|------|------|
| `Client() *redis.Client` | 获取原生 Redis 客户端 |
| `Close() error` | 关闭连接 |
| `Ping(ctx) error` | 健康检查 |

### 字符串操作

| 方法 | 说明 |
|------|------|
| `Get(ctx, key) (string, error)` | 获取值 |
| `Set(ctx, key, value, exp) error` | 设置值 |
| `SetNX(ctx, key, value, exp) (bool, error)` | 仅在不存在时设置 |
| `MGet(ctx, keys...) ([]string, error)` | 批量获取 |
| `MSet(ctx, pairs...) error` | 批量设置 |
| `Incr(ctx, key) (int64, error)` | 自增 |
| `IncrBy(ctx, key, value) (int64, error)` | 增加指定值 |
| `Decr(ctx, key) (int64, error)` | 自减 |
| `DecrBy(ctx, key, value) (int64, error)` | 减少指定值 |

### 哈希操作

| 方法 | 说明 |
|------|------|
| `HGet(ctx, key, field) (string, error)` | 获取字段值 |
| `HSet(ctx, key, values...) (int64, error)` | 设置字段值 |
| `HGetAll(ctx, key) (map[string]string, error)` | 获取所有字段 |
| `HDel(ctx, key, fields...) (int64, error)` | 删除字段 |
| `HExists(ctx, key, field) (bool, error)` | 检查字段是否存在 |
| `HLen(ctx, key) (int64, error)` | 获取字段数量 |
| `HIncrBy(ctx, key, field, incr) (int64, error)` | 字段值自增 |

### 列表操作

| 方法 | 说明 |
|------|------|
| `LPush(ctx, key, values...) (int64, error)` | 从左侧插入 |
| `RPush(ctx, key, values...) (int64, error)` | 从右侧插入 |
| `LPop(ctx, key) (string, error)` | 从左侧弹出 |
| `RPop(ctx, key) (string, error)` | 从右侧弹出 |
| `LRange(ctx, key, start, stop) ([]string, error)` | 获取范围元素 |
| `LLen(ctx, key) (int64, error)` | 获取列表长度 |
| `LIndex(ctx, key, index) (string, error)` | 获取指定索引元素 |

### 集合操作

| 方法 | 说明 |
|------|------|
| `SAdd(ctx, key, members...) (int64, error)` | 添加成员 |
| `SRem(ctx, key, members...) (int64, error)` | 移除成员 |
| `SMembers(ctx, key) ([]string, error)` | 获取所有成员 |
| `SIsMember(ctx, key, member) (bool, error)` | 判断成员是否存在 |
| `SCard(ctx, key) (int64, error)` | 获取成员数量 |
| `SUnion(ctx, keys...) ([]string, error)` | 并集 |
| `SInter(ctx, keys...) ([]string, error)` | 交集 |
| `SDiff(ctx, keys...) ([]string, error)` | 差集 |

### 有序集合操作

| 方法 | 说明 |
|------|------|
| `ZAdd(ctx, key, members...) (int64, error)` | 添加成员 |
| `ZRem(ctx, key, members...) (int64, error)` | 移除成员 |
| `ZRange(ctx, key, start, stop) ([]string, error)` | 获取范围成员 |
| `ZRangeWithScores(ctx, key, start, stop) ([]Z, error)` | 获取范围成员及分数 |
| `ZScore(ctx, key, member) (float64, error)` | 获取成员分数 |
| `ZRank(ctx, key, member) (int64, error)` | 获取成员排名 |
| `ZCard(ctx, key) (int64, error)` | 获取成员数量 |

### 键管理

| 方法 | 说明 |
|------|------|
| `Del(ctx, keys...) (int64, error)` | 删除键 |
| `Exists(ctx, keys...) (int64, error)` | 检查键是否存在 |
| `Expire(ctx, key, exp) (bool, error)` | 设置过期时间 |
| `TTL(ctx, key) (time.Duration, error)` | 获取剩余生存时间 |
| `Persist(ctx, key) (bool, error)` | 移除过期时间 |
| `Keys(ctx, pattern) ([]string, error)` | 查找匹配的键 |
| `Rename(ctx, key, newKey) error` | 重命名键 |
| `Type(ctx, key) (string, error)` | 获取键类型 |

## 🎯 最佳实践

### 1. 使用链路追踪

```go
// 使用 ucontext 创建追踪上下文
ctx := ucontext.NewContext(context.Background())

// 所有操作自动包含 trace_id
conn.Set(ctx, "key", "value", 0)
```

### 2. 错误处理

```go
value, err := conn.Get(ctx, "key")
if err != nil {
    // 判断是否为键不存在
    if err == redis.ErrNil {
        // 键不存在，使用默认值
        value = "default"
    } else {
        // 其他错误
        return err
    }
}
```

### 3. 使用 Pipeline 提升性能

```go
// 批量操作使用 Pipeline
conn.Pipelined(ctx, func(pipe redis.Pipeliner) error {
    for i := 0; i < 1000; i++ {
        pipe.Set(ctx, fmt.Sprintf("key:%d", i), i, 0)
    }
    return nil
})
```

### 4. 合理设置过期时间

```go
// 缓存数据设置合理的过期时间
conn.Set(ctx, "cache:user:1", userData, 10*time.Minute)
```

### 5. 使用类型转换

```go
import "github.com/whosafe/uf/uconv"

// 获取字符串并转换为整数
value, _ := conn.Get(ctx, "counter")
count := uconv.ToIntDef(value, 0)

// 获取字符串并转换为布尔值
value, _ := conn.Get(ctx, "flag")
flag := uconv.ToBoolDef(value, false)
```

## 📝 示例代码

完整示例请参考 `example/udb/redis/` 目录。

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License
