package main

import (
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/whosafe/uf/ucontext"
	redisdb "github.com/whosafe/uf/udb/redis"
)

func main() {
	fmt.Println("=== Redis 模块示例程序 ===")

	// 加载配置
	// 注意：这里使用默认配置，实际使用时应该创建 config.yaml
	config := redisdb.DefaultConfig()
	config.Host = "10.0.2.128"
	config.Port = 6379
	config.DB = 0

	// 创建连接
	conn, err := redisdb.New(config)
	if err != nil {
		log.Fatal("连接 Redis 失败:", err)
	}
	defer conn.Close()

	// 创建追踪上下文
	ctx := ucontext.New()

	// 测试连接
	fmt.Println("1. 测试连接...")
	if err := conn.Ping(ctx); err != nil {
		log.Fatal("Ping 失败:", err)
	}
	fmt.Println("✓ 连接成功")

	// 测试字符串操作
	fmt.Println("2. 测试字符串操作...")
	testStringOps(ctx, conn)

	// 测试哈希操作
	fmt.Println("\n3. 测试哈希操作...")
	testHashOps(ctx, conn)

	// 测试列表操作
	fmt.Println("\n4. 测试列表操作...")
	testListOps(ctx, conn)

	// 测试集合操作
	fmt.Println("\n5. 测试集合操作...")
	testSetOps(ctx, conn)

	// 测试有序集合操作
	fmt.Println("\n6. 测试有序集合操作...")
	testZSetOps(ctx, conn)

	// 测试 Pipeline
	fmt.Println("\n7. 测试 Pipeline...")
	testPipeline(ctx, conn)

	// 测试 Pub/Sub
	fmt.Println("\n8. 测试 Pub/Sub...")
	testPubSub(ctx, conn)

	// 清理测试数据
	fmt.Println("\n9. 清理测试数据...")
	conn.Del(ctx, "test:string", "test:hash", "test:list", "test:set", "test:zset")
	fmt.Println("✓ 清理完成")

	fmt.Println("\n=== ✓ 所有测试完成! ===")
}

func testStringOps(ctx *ucontext.Context, conn *redisdb.Connection) {
	// SET
	err := conn.Set(ctx, "test:string", "Hello, Redis!", 10*time.Minute)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("✓ SET test:string")

	// GET
	value, err := conn.Get(ctx, "test:string")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("✓ GET test:string = %s\n", value)

	// INCR
	conn.Set(ctx, "test:counter", "0", 0)
	count, _ := conn.Incr(ctx, "test:counter")
	fmt.Printf("✓ INCR test:counter = %d\n", count)

	// INCRBY
	count, _ = conn.IncrBy(ctx, "test:counter", 10)
	fmt.Printf("✓ INCRBY test:counter 10 = %d\n", count)
}

func testHashOps(ctx *ucontext.Context, conn *redisdb.Connection) {
	// HSET
	conn.HSet(ctx, "test:hash", "name", "Alice", "age", 25, "city", "Beijing")
	fmt.Println("✓ HSET test:hash")

	// HGET
	name, _ := conn.HGet(ctx, "test:hash", "name")
	fmt.Printf("✓ HGET test:hash name = %s\n", name)

	// HGETALL
	fields, _ := conn.HGetAll(ctx, "test:hash")
	fmt.Printf("✓ HGETALL test:hash = %v\n", fields)

	// HINCRBY
	newAge, _ := conn.HIncrBy(ctx, "test:hash", "age", 1)
	fmt.Printf("✓ HINCRBY test:hash age 1 = %d\n", newAge)
}

func testListOps(ctx *ucontext.Context, conn *redisdb.Connection) {
	// LPUSH
	conn.LPush(ctx, "test:list", "item1", "item2", "item3")
	fmt.Println("✓ LPUSH test:list")

	// RPUSH
	conn.RPush(ctx, "test:list", "item4", "item5")
	fmt.Println("✓ RPUSH test:list")

	// LRANGE
	items, _ := conn.LRange(ctx, "test:list", 0, -1)
	fmt.Printf("✓ LRANGE test:list 0 -1 = %v\n", items)

	// LLEN
	length, _ := conn.LLen(ctx, "test:list")
	fmt.Printf("✓ LLEN test:list = %d\n", length)

	// LPOP
	item, _ := conn.LPop(ctx, "test:list")
	fmt.Printf("✓ LPOP test:list = %s\n", item)
}

func testSetOps(ctx *ucontext.Context, conn *redisdb.Connection) {
	// SADD
	conn.SAdd(ctx, "test:set", "member1", "member2", "member3")
	fmt.Println("✓ SADD test:set")

	// SMEMBERS
	members, _ := conn.SMembers(ctx, "test:set")
	fmt.Printf("✓ SMEMBERS test:set = %v\n", members)

	// SISMEMBER
	exists, _ := conn.SIsMember(ctx, "test:set", "member1")
	fmt.Printf("✓ SISMEMBER test:set member1 = %v\n", exists)

	// SCARD
	count, _ := conn.SCard(ctx, "test:set")
	fmt.Printf("✓ SCARD test:set = %d\n", count)
}

func testZSetOps(ctx *ucontext.Context, conn *redisdb.Connection) {
	// ZADD
	conn.ZAdd(ctx, "test:zset",
		redis.Z{Score: 100, Member: "Alice"},
		redis.Z{Score: 95, Member: "Bob"},
		redis.Z{Score: 90, Member: "Charlie"})
	fmt.Println("✓ ZADD test:zset")

	// ZRANGE
	members, _ := conn.ZRange(ctx, "test:zset", 0, -1)
	fmt.Printf("✓ ZRANGE test:zset 0 -1 = %v\n", members)

	// ZREVRANGE
	members, _ = conn.ZRevRange(ctx, "test:zset", 0, -1)
	fmt.Printf("✓ ZREVRANGE test:zset 0 -1 = %v\n", members)

	// ZSCORE
	score, _ := conn.ZScore(ctx, "test:zset", "Alice")
	fmt.Printf("✓ ZSCORE test:zset Alice = %.0f\n", score)

	// ZRANK
	rank, _ := conn.ZRank(ctx, "test:zset", "Alice")
	fmt.Printf("✓ ZRANK test:zset Alice = %d\n", rank)
}

func testPipeline(ctx *ucontext.Context, conn *redisdb.Connection) {
	_, err := conn.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Set(ctx.Context(), "test:pipe1", "value1", 0)
		pipe.Set(ctx.Context(), "test:pipe2", "value2", 0)
		pipe.Set(ctx.Context(), "test:pipe3", "value3", 0)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("✓ Pipeline 执行成功")

	// 清理
	conn.Del(ctx, "test:pipe1", "test:pipe2", "test:pipe3")
}

func testPubSub(ctx *ucontext.Context, conn *redisdb.Connection) {
	// 发布消息
	subscribers, _ := conn.Publish(ctx, "test:channel", "Hello, Pub/Sub!")
	fmt.Printf("✓ PUBLISH test:channel (订阅者数: %d)\n", subscribers)
}
