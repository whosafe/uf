package redis

import (
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/whosafe/uf/ucontext"
)

// 测试配置
func getTestConfig() *Config {
	config := DefaultConfig()
	config.Host = "10.0.2.128"
	config.Port = 6379
	config.DB = 0
	return config
}

// 测试连接创建
func TestNew(t *testing.T) {
	config := getTestConfig()
	conn, err := New(config)
	if err != nil {
		t.Fatalf("创建连接失败: %v", err)
	}
	defer conn.Close()

	if conn.client == nil {
		t.Error("客户端为空")
	}
}

// 测试 Ping
func TestPing(t *testing.T) {
	config := getTestConfig()
	conn, err := New(config)
	if err != nil {
		t.Fatalf("创建连接失败: %v", err)
	}
	defer conn.Close()

	ctx := ucontext.New()
	err = conn.Ping(ctx)
	if err != nil {
		t.Errorf("Ping 失败: %v", err)
	}
}

// 测试字符串操作
func TestStringOperations(t *testing.T) {
	config := getTestConfig()
	conn, err := New(config)
	if err != nil {
		t.Fatalf("创建连接失败: %v", err)
	}
	defer conn.Close()

	ctx := ucontext.New()
	key := "test:string"

	// 清理
	defer conn.Del(ctx, key)

	// 测试 SET
	err = conn.Set(ctx, key, "Hello, Redis!", 1*time.Minute)
	if err != nil {
		t.Errorf("SET 失败: %v", err)
	}

	// 测试 GET
	value, err := conn.Get(ctx, key)
	if err != nil {
		t.Errorf("GET 失败: %v", err)
	}
	if value != "Hello, Redis!" {
		t.Errorf("GET 返回值错误: 期望 'Hello, Redis!', 实际 '%s'", value)
	}

	// 测试 INCR
	counterKey := "test:counter"
	defer conn.Del(ctx, counterKey)

	conn.Set(ctx, counterKey, "0", 0)
	count, err := conn.Incr(ctx, counterKey)
	if err != nil {
		t.Errorf("INCR 失败: %v", err)
	}
	if count != 1 {
		t.Errorf("INCR 返回值错误: 期望 1, 实际 %d", count)
	}

	// 测试 INCRBY
	count, err = conn.IncrBy(ctx, counterKey, 10)
	if err != nil {
		t.Errorf("INCRBY 失败: %v", err)
	}
	if count != 11 {
		t.Errorf("INCRBY 返回值错误: 期望 11, 实际 %d", count)
	}
}

// 测试哈希操作
func TestHashOperations(t *testing.T) {
	config := getTestConfig()
	conn, err := New(config)
	if err != nil {
		t.Fatalf("创建连接失败: %v", err)
	}
	defer conn.Close()

	ctx := ucontext.New()
	key := "test:hash"
	defer conn.Del(ctx, key)

	// 测试 HSET
	_, err = conn.HSet(ctx, key, "name", "Alice", "age", 25)
	if err != nil {
		t.Errorf("HSET 失败: %v", err)
	}

	// 测试 HGET
	name, err := conn.HGet(ctx, key, "name")
	if err != nil {
		t.Errorf("HGET 失败: %v", err)
	}
	if name != "Alice" {
		t.Errorf("HGET 返回值错误: 期望 'Alice', 实际 '%s'", name)
	}

	// 测试 HGETALL
	fields, err := conn.HGetAll(ctx, key)
	if err != nil {
		t.Errorf("HGETALL 失败: %v", err)
	}
	if len(fields) != 2 {
		t.Errorf("HGETALL 返回字段数错误: 期望 2, 实际 %d", len(fields))
	}

	// 测试 HINCRBY
	newAge, err := conn.HIncrBy(ctx, key, "age", 1)
	if err != nil {
		t.Errorf("HINCRBY 失败: %v", err)
	}
	if newAge != 26 {
		t.Errorf("HINCRBY 返回值错误: 期望 26, 实际 %d", newAge)
	}
}

// 测试列表操作
func TestListOperations(t *testing.T) {
	config := getTestConfig()
	conn, err := New(config)
	if err != nil {
		t.Fatalf("创建连接失败: %v", err)
	}
	defer conn.Close()

	ctx := ucontext.New()
	key := "test:list"
	defer conn.Del(ctx, key)

	// 测试 LPUSH
	_, err = conn.LPush(ctx, key, "item1", "item2", "item3")
	if err != nil {
		t.Errorf("LPUSH 失败: %v", err)
	}

	// 测试 LLEN
	length, err := conn.LLen(ctx, key)
	if err != nil {
		t.Errorf("LLEN 失败: %v", err)
	}
	if length != 3 {
		t.Errorf("LLEN 返回值错误: 期望 3, 实际 %d", length)
	}

	// 测试 LRANGE
	items, err := conn.LRange(ctx, key, 0, -1)
	if err != nil {
		t.Errorf("LRANGE 失败: %v", err)
	}
	if len(items) != 3 {
		t.Errorf("LRANGE 返回元素数错误: 期望 3, 实际 %d", len(items))
	}

	// 测试 LPOP
	item, err := conn.LPop(ctx, key)
	if err != nil {
		t.Errorf("LPOP 失败: %v", err)
	}
	if item != "item3" {
		t.Errorf("LPOP 返回值错误: 期望 'item3', 实际 '%s'", item)
	}
}

// 测试集合操作
func TestSetOperations(t *testing.T) {
	config := getTestConfig()
	conn, err := New(config)
	if err != nil {
		t.Fatalf("创建连接失败: %v", err)
	}
	defer conn.Close()

	ctx := ucontext.New()
	key := "test:set"
	defer conn.Del(ctx, key)

	// 测试 SADD
	_, err = conn.SAdd(ctx, key, "member1", "member2", "member3")
	if err != nil {
		t.Errorf("SADD 失败: %v", err)
	}

	// 测试 SCARD
	count, err := conn.SCard(ctx, key)
	if err != nil {
		t.Errorf("SCARD 失败: %v", err)
	}
	if count != 3 {
		t.Errorf("SCARD 返回值错误: 期望 3, 实际 %d", count)
	}

	// 测试 SISMEMBER
	exists, err := conn.SIsMember(ctx, key, "member1")
	if err != nil {
		t.Errorf("SISMEMBER 失败: %v", err)
	}
	if !exists {
		t.Error("SISMEMBER 返回值错误: 期望 true, 实际 false")
	}

	// 测试 SMEMBERS
	members, err := conn.SMembers(ctx, key)
	if err != nil {
		t.Errorf("SMEMBERS 失败: %v", err)
	}
	if len(members) != 3 {
		t.Errorf("SMEMBERS 返回成员数错误: 期望 3, 实际 %d", len(members))
	}
}

// 测试有序集合操作
func TestZSetOperations(t *testing.T) {
	config := getTestConfig()
	conn, err := New(config)
	if err != nil {
		t.Fatalf("创建连接失败: %v", err)
	}
	defer conn.Close()

	ctx := ucontext.New()
	key := "test:zset"
	defer conn.Del(ctx, key)

	// 测试 ZADD
	_, err = conn.ZAdd(ctx, key,
		redis.Z{Score: 100, Member: "Alice"},
		redis.Z{Score: 95, Member: "Bob"},
		redis.Z{Score: 90, Member: "Charlie"})
	if err != nil {
		t.Errorf("ZADD 失败: %v", err)
	}

	// 测试 ZCARD
	count, err := conn.ZCard(ctx, key)
	if err != nil {
		t.Errorf("ZCARD 失败: %v", err)
	}
	if count != 3 {
		t.Errorf("ZCARD 返回值错误: 期望 3, 实际 %d", count)
	}

	// 测试 ZSCORE
	score, err := conn.ZScore(ctx, key, "Alice")
	if err != nil {
		t.Errorf("ZSCORE 失败: %v", err)
	}
	if score != 100 {
		t.Errorf("ZSCORE 返回值错误: 期望 100, 实际 %.0f", score)
	}

	// 测试 ZRANGE
	members, err := conn.ZRange(ctx, key, 0, -1)
	if err != nil {
		t.Errorf("ZRANGE 失败: %v", err)
	}
	if len(members) != 3 {
		t.Errorf("ZRANGE 返回成员数错误: 期望 3, 实际 %d", len(members))
	}
	if members[0] != "Charlie" || members[2] != "Alice" {
		t.Errorf("ZRANGE 返回顺序错误: %v", members)
	}
}

// 测试键管理操作
func TestKeyOperations(t *testing.T) {
	config := getTestConfig()
	conn, err := New(config)
	if err != nil {
		t.Fatalf("创建连接失败: %v", err)
	}
	defer conn.Close()

	ctx := ucontext.New()
	key := "test:key"

	// 设置测试键
	conn.Set(ctx, key, "value", 0)
	defer conn.Del(ctx, key)

	// 测试 EXISTS
	count, err := conn.Exists(ctx, key)
	if err != nil {
		t.Errorf("EXISTS 失败: %v", err)
	}
	if count != 1 {
		t.Errorf("EXISTS 返回值错误: 期望 1, 实际 %d", count)
	}

	// 测试 EXPIRE
	success, err := conn.Expire(ctx, key, 10*time.Second)
	if err != nil {
		t.Errorf("EXPIRE 失败: %v", err)
	}
	if !success {
		t.Error("EXPIRE 返回值错误: 期望 true, 实际 false")
	}

	// 测试 TTL
	ttl, err := conn.TTL(ctx, key)
	if err != nil {
		t.Errorf("TTL 失败: %v", err)
	}
	if ttl <= 0 {
		t.Errorf("TTL 返回值错误: 期望 > 0, 实际 %v", ttl)
	}

	// 测试 DEL
	count, err = conn.Del(ctx, key)
	if err != nil {
		t.Errorf("DEL 失败: %v", err)
	}
	if count != 1 {
		t.Errorf("DEL 返回值错误: 期望 1, 实际 %d", count)
	}
}

// 测试 Pipeline
func TestPipeline(t *testing.T) {
	config := getTestConfig()
	conn, err := New(config)
	if err != nil {
		t.Fatalf("创建连接失败: %v", err)
	}
	defer conn.Close()

	ctx := ucontext.New()
	keys := []string{"test:pipe1", "test:pipe2", "test:pipe3"}
	defer conn.Del(ctx, keys...)

	// 测试 Pipelined
	cmds, err := conn.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Set(ctx.Context(), "test:pipe1", "value1", 0)
		pipe.Set(ctx.Context(), "test:pipe2", "value2", 0)
		pipe.Set(ctx.Context(), "test:pipe3", "value3", 0)
		return nil
	})
	if err != nil {
		t.Errorf("Pipelined 失败: %v", err)
	}
	if len(cmds) != 3 {
		t.Errorf("Pipelined 返回命令数错误: 期望 3, 实际 %d", len(cmds))
	}

	// 验证数据
	value, err := conn.Get(ctx, "test:pipe1")
	if err != nil {
		t.Errorf("验证 Pipeline 数据失败: %v", err)
	}
	if value != "value1" {
		t.Errorf("Pipeline 数据错误: 期望 'value1', 实际 '%s'", value)
	}
}

// 测试 Pub/Sub
func TestPubSub(t *testing.T) {
	config := getTestConfig()
	conn, err := New(config)
	if err != nil {
		t.Fatalf("创建连接失败: %v", err)
	}
	defer conn.Close()

	ctx := ucontext.New()

	// 测试 Publish（没有订阅者）
	subscribers, err := conn.Publish(ctx, "test:channel", "Hello")
	if err != nil {
		t.Errorf("Publish 失败: %v", err)
	}
	if subscribers != 0 {
		t.Errorf("Publish 返回订阅者数错误: 期望 0, 实际 %d", subscribers)
	}
}
