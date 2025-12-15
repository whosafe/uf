package redis

import (
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/whosafe/uf/ucontext"
)

// 性能测试配置
func getBenchConfig() *Config {
	config := DefaultConfig()
	config.Host = "10.0.2.128"
	config.Port = 6379
	config.DB = 0
	config.Pool.PoolSize = 100 // 增加连接池大小以支持并发测试
	return config
}

// 基准测试：字符串 SET 操作
func BenchmarkSet(b *testing.B) {
	conn, err := New(getBenchConfig())
	if err != nil {
		b.Fatalf("创建连接失败: %v", err)
	}
	defer conn.Close()

	ctx := ucontext.New()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		conn.Set(ctx, "bench:key", "value", 1*time.Minute)
	}
}

// 基准测试：字符串 GET 操作
func BenchmarkGet(b *testing.B) {
	conn, err := New(getBenchConfig())
	if err != nil {
		b.Fatalf("创建连接失败: %v", err)
	}
	defer conn.Close()

	ctx := ucontext.New()
	conn.Set(ctx, "bench:key", "value", 1*time.Minute)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		conn.Get(ctx, "bench:key")
	}
}

// 基准测试：哈希 HSET 操作
func BenchmarkHSet(b *testing.B) {
	conn, err := New(getBenchConfig())
	if err != nil {
		b.Fatalf("创建连接失败: %v", err)
	}
	defer conn.Close()

	ctx := ucontext.New()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		conn.HSet(ctx, "bench:hash", "field", "value")
	}
}

// 基准测试：哈希 HGET 操作
func BenchmarkHGet(b *testing.B) {
	conn, err := New(getBenchConfig())
	if err != nil {
		b.Fatalf("创建连接失败: %v", err)
	}
	defer conn.Close()

	ctx := ucontext.New()
	conn.HSet(ctx, "bench:hash", "field", "value")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		conn.HGet(ctx, "bench:hash", "field")
	}
}

// 基准测试：列表 LPUSH 操作
func BenchmarkLPush(b *testing.B) {
	conn, err := New(getBenchConfig())
	if err != nil {
		b.Fatalf("创建连接失败: %v", err)
	}
	defer conn.Close()

	ctx := ucontext.New()
	defer conn.Del(ctx, "bench:list")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		conn.LPush(ctx, "bench:list", "item")
	}
}

// 基准测试：列表 LPOP 操作
func BenchmarkLPop(b *testing.B) {
	conn, err := New(getBenchConfig())
	if err != nil {
		b.Fatalf("创建连接失败: %v", err)
	}
	defer conn.Close()

	ctx := ucontext.New()

	// 预填充数据
	for i := 0; i < b.N; i++ {
		conn.LPush(ctx, "bench:list", "item")
	}
	defer conn.Del(ctx, "bench:list")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		conn.LPop(ctx, "bench:list")
	}
}

// 基准测试：集合 SADD 操作
func BenchmarkSAdd(b *testing.B) {
	conn, err := New(getBenchConfig())
	if err != nil {
		b.Fatalf("创建连接失败: %v", err)
	}
	defer conn.Close()

	ctx := ucontext.New()
	defer conn.Del(ctx, "bench:set")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		conn.SAdd(ctx, "bench:set", "member")
	}
}

// 基准测试：有序集合 ZADD 操作
func BenchmarkZAdd(b *testing.B) {
	conn, err := New(getBenchConfig())
	if err != nil {
		b.Fatalf("创建连接失败: %v", err)
	}
	defer conn.Close()

	ctx := ucontext.New()
	defer conn.Del(ctx, "bench:zset")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		conn.ZAdd(ctx, "bench:zset", redis.Z{Score: float64(i), Member: "member"})
	}
}

// 基准测试：Pipeline 批量操作
func BenchmarkPipeline(b *testing.B) {
	conn, err := New(getBenchConfig())
	if err != nil {
		b.Fatalf("创建连接失败: %v", err)
	}
	defer conn.Close()

	ctx := ucontext.New()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		conn.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Set(ctx.Context(), "bench:pipe1", "value1", 0)
			pipe.Set(ctx.Context(), "bench:pipe2", "value2", 0)
			pipe.Set(ctx.Context(), "bench:pipe3", "value3", 0)
			return nil
		})
	}
}

// 基准测试：并发 SET 操作
func BenchmarkSetParallel(b *testing.B) {
	conn, err := New(getBenchConfig())
	if err != nil {
		b.Fatalf("创建连接失败: %v", err)
	}
	defer conn.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		ctx := ucontext.New()
		i := 0
		for pb.Next() {
			conn.Set(ctx, "bench:parallel", "value", 1*time.Minute)
			i++
		}
	})
}

// 基准测试：并发 GET 操作
func BenchmarkGetParallel(b *testing.B) {
	conn, err := New(getBenchConfig())
	if err != nil {
		b.Fatalf("创建连接失败: %v", err)
	}
	defer conn.Close()

	ctx := ucontext.New()
	conn.Set(ctx, "bench:parallel", "value", 1*time.Minute)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		ctx := ucontext.New()
		for pb.Next() {
			conn.Get(ctx, "bench:parallel")
		}
	})
}

// 基准测试：混合操作
func BenchmarkMixedOperations(b *testing.B) {
	conn, err := New(getBenchConfig())
	if err != nil {
		b.Fatalf("创建连接失败: %v", err)
	}
	defer conn.Close()

	ctx := ucontext.New()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// 字符串操作
		conn.Set(ctx, "bench:mixed:str", "value", 1*time.Minute)
		conn.Get(ctx, "bench:mixed:str")

		// 哈希操作
		conn.HSet(ctx, "bench:mixed:hash", "field", "value")
		conn.HGet(ctx, "bench:mixed:hash", "field")

		// 列表操作
		conn.LPush(ctx, "bench:mixed:list", "item")
		conn.LPop(ctx, "bench:mixed:list")
	}
}
