package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/whosafe/uf/ucontext"
	"github.com/whosafe/uf/uconv"
	"github.com/whosafe/uf/udb/postgresql"
)

// User 用户结构体
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

func main() {
	// 创建配置
	config := &postgresql.Config{
		Host:     "10.0.2.128",
		Port:     5432,
		Username: "postgres",
		Password: "111111",
		Database: "testdb",
		SSLMode:  "disable",
	}

	// 创建客户端
	conn, err := postgresql.New(config)
	if err != nil {
		log.Fatal("连接失败:", err)
	}
	defer conn.Close()

	// 创建追踪上下文
	ctx := ucontext.NewContext(context.Background())

	// 1. 测试 Ping
	fmt.Println("=== 测试 Ping ===")
	if err := conn.Ping(ctx); err != nil {
		log.Fatal("Ping 失败:", err)
	}
	fmt.Println(" Ping 成功")

	// 2. 查询单条记录
	fmt.Println("\n=== 查询单条记录 ===")
	var user User
	err = conn.Query(ctx).
		Table("users").
		Where("id = ?", 1).
		Scan(&user)

	if err != nil {
		if err == postgresql.ErrNoRows {
			fmt.Println("未找到记录")
		} else {
			log.Fatal("查询失败:", err)
		}
	} else {
		fmt.Printf("用户: %+v\n", user)
	}

	// 3. 查询多条记录
	fmt.Println("\n=== 查询多条记录 ===")
	results, err := conn.Query(ctx).
		Table("users").
		Where("age > ?", 18).
		OrderBy("id").
		Limit(10).
		ScanAll(func() postgresql.Scanner { return &User{} })

	if err != nil {
		log.Fatal("查询失败:", err)
	}

	fmt.Printf("找到 %d 条记录:\n", len(results))
	for _, r := range results {
		u := r.(*User)
		fmt.Printf("  - %s (%s)\n", u.Username, u.Email)
	}

	// 4. 插入数据
	fmt.Println("\n=== 插入数据 ===")
	affected, err := conn.Exec(ctx,
		"INSERT INTO users (username, email, age) VALUES ($1, $2, $3)",
		"test_user", "test@example.com", 25)
	if err != nil {
		log.Printf("插入失败: %v\n", err)
	} else {
		fmt.Printf(" 插入成功，影响行数: %d\n", affected)
	}

	// 5. 更新数据
	fmt.Println("\n=== 更新数据 ===")
	affected, err = conn.Exec(ctx,
		"UPDATE users SET age = $1 WHERE username = $2",
		26, "test_user")
	if err != nil {
		log.Printf("更新失败: %v\n", err)
	} else {
		fmt.Printf(" 更新成功，影响行数: %d\n", affected)
	}

	// 6. 删除数据
	fmt.Println("\n=== 删除数据 ===")
	affected, err = conn.Exec(ctx,
		"DELETE FROM users WHERE username = $1",
		"test_user")
	if err != nil {
		log.Printf("删除失败: %v\n", err)
	} else {
		fmt.Printf(" 删除成功，影响行数: %d\n", affected)
	}

	// 7. 连接池统计
	fmt.Println("\n=== 连接池统计 ===")
	stats := conn.Stats()
	fmt.Printf("总连接数: %d\n", stats.TotalConns())
	fmt.Printf("空闲连接数: %d\n", stats.IdleConns())
	fmt.Printf("获取连接数: %d\n", stats.AcquiredConns())
}
