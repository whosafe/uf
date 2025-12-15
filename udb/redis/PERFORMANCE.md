# Redis 性能优化最佳实践

本文档提供 Redis 模块的性能优化建议和最佳实践。

## 配置优化

### 1. 连接池配置

根据不同的应用场景选择合适的连接池配置：

#### 默认配置（中小型应用，QPS < 1000）

```go
config := redis.DefaultConfig()
// PoolSize: 10
// MinIdleConn: 5
// MaxActive: 100
```

#### 高性能配置（高并发场景，QPS > 5000）

```go
config := redis.HighPerformanceConfig()
// PoolSize: 50
// MinIdleConn: 20
// MaxActive: 500
// SlowQueryThreshold: 50ms
```

#### 低延迟配置（延迟敏感场景）

```go
config := redis.LowLatencyConfig()
// PoolSize: 30
// MinIdleConn: 15
// SlowQueryThreshold: 20ms
```

### 2. 连接池大小计算公式

```
PoolSize = (并发请求数 × 平均响应时间) / 1000
```

示例：

- 并发请求数：1000 QPS
- 平均响应时间：10ms
- 推荐 PoolSize：(1000 × 10) / 1000 = 10

### 3. 最小空闲连接数建议

```
MinIdleConn = PoolSize × 0.5
```

保持一定数量的预热连接可以减少连接建立的延迟。

## 性能优化技巧

### 1. 使用 Pipeline 批量操作

**不推荐**（单个操作）：

```go
for i := 0; i < 100; i++ {
    conn.Set(ctx, fmt.Sprintf("key%d", i), "value", 0)
}
```

**推荐**（Pipeline 批量）：

```go
conn.Pipelined(ctx, func(pipe redis.Pipeliner) error {
    for i := 0; i < 100; i++ {
        pipe.Set(ctx.Context(), fmt.Sprintf("key%d", i), "value", 0)
    }
    return nil
})
```

**性能提升**：3-5 倍

### 2. 合理设置超时时间

根据业务场景设置合适的超时时间：

```go
// 快速查询（缓存读取）
config.Query.DefaultTimeout = 5 * time.Second

// 复杂操作（批量写入）
config.Query.DefaultTimeout = 30 * time.Second

// 后台任务（非关键路径）
config.Query.DefaultTimeout = 60 * time.Second
```

### 3. 慢查询监控

启用慢查询日志并设置合理的阈值：

```go
config.Log.SlowQuery = true
config.Query.SlowQueryThreshold = 100 * time.Millisecond
```

**阈值建议**：

- 缓存场景：20-50ms
- 一般场景：100ms
- 后台任务：500ms

### 4. 避免大 Key

**不推荐**：

```go
// 单个哈希存储大量字段
conn.HSet(ctx, "user:all", field1, value1, field2, value2, ...) // 10000+ 字段
```

**推荐**：

```go
// 分片存储
conn.HSet(ctx, "user:1000", field1, value1, ...)
conn.HSet(ctx, "user:2000", field2, value2, ...)
```

### 5. 使用合适的数据结构

| 场景 | 推荐数据结构 | 原因 |
|------|------------|------|
| 简单键值 | String | 最快 |
| 对象存储 | Hash | 节省内存 |
| 列表/队列 | List | 支持阻塞操作 |
| 去重/集合运算 | Set | O(1) 查找 |
| 排行榜 | ZSet | 自动排序 |

## 内存优化

### 1. 设置过期时间

```go
// 推荐：为所有 key 设置过期时间
conn.Set(ctx, "cache:user:1", data, 10*time.Minute)

// 不推荐：永久存储
conn.Set(ctx, "cache:user:1", data, 0)
```

### 2. 使用压缩

对于大数据，考虑压缩后存储：

```go
import "compress/gzip"

// 压缩数据
compressed := compress(data)
conn.Set(ctx, "large:data", compressed, 1*time.Hour)

// 读取并解压
compressed, _ := conn.Get(ctx, "large:data")
data := decompress(compressed)
```

## 并发优化

### 1. 连接复用

```go
// 推荐：复用连接
var conn *redis.Connection
conn, _ = redis.New(config)
defer conn.Close()

// 多个 goroutine 共享同一个连接
for i := 0; i < 100; i++ {
    go func() {
        ctx := ucontext.New()
        conn.Get(ctx, "key")
    }()
}
```

### 2. 避免连接泄漏

```go
// 确保连接关闭
conn, err := redis.New(config)
if err != nil {
    return err
}
defer conn.Close() // 重要！
```

## 监控指标

### 关键指标

1. **QPS**（每秒查询数）
2. **平均响应时间**
3. **慢查询数量**
4. **连接池使用率**
5. **错误率**

### 监控示例

```go
// 启用详细日志
config.Log.Enabled = true
config.Log.SlowQuery = true
config.Log.LogParams = true // 开发环境

// 生产环境建议
config.Log.LogParams = false // 避免敏感信息泄露
```

## 故障排查

### 常见问题

#### 1. 连接池耗尽

**症状**：大量超时错误

**解决**：

```go
// 增加连接池大小
config.Pool.PoolSize = 50
config.Pool.MaxActive = 500
```

#### 2. 慢查询过多

**症状**：响应时间长

**解决**：

- 使用 Pipeline 批量操作
- 优化数据结构
- 添加索引（使用 ZSet）

#### 3. 内存占用高

**症状**：Redis 内存持续增长

**解决**：

- 设置过期时间
- 删除不用的 key
- 使用 Hash 代替多个 String

## 生产环境检查清单

- [ ] 连接池大小根据 QPS 合理配置
- [ ] 设置了合适的超时时间
- [ ] 启用慢查询日志
- [ ] 关闭参数日志（避免敏感信息）
- [ ] 所有 key 都设置了过期时间
- [ ] 使用 Pipeline 进行批量操作
- [ ] 添加了监控和告警
- [ ] 定期检查慢查询日志
- [ ] 连接正确关闭（defer conn.Close()）

## 性能基准

参考 `redis_bench_test.go` 中的基准测试：

```bash
# 运行性能测试
go test -bench=. -benchmem ./udb/redis

# 对比优化前后性能
go test -bench=BenchmarkSet -benchmem ./udb/redis
```

## 更多资源

- [Redis 官方文档](https://redis.io/docs/)
- [Redis 最佳实践](https://redis.io/docs/manual/patterns/)
- [性能测试报告](performance_report.md)
