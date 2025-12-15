package main

import (
	"context"
	"fmt"
	"log"

	"github.com/whosafe/uf/ucontext"
	"github.com/whosafe/uf/uconv"
	"github.com/whosafe/uf/udb/postgresql"
)

// User 用户结构体
type User struct {
	ID       int64
	Username string
	Email    string
	Age      int
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

	fmt.Println("=== 事务示例 ===")

	// 1. 基本事务
	fmt.Println("\n--- 基本事务 ---")
	tx, err := conn.Begin(ctx)
	if err != nil {
		log.Fatal("开始事务失败:", err)
	}

	// 执行多个操作
	_, err = tx.Exec(ctx, "INSERT INTO users (username, email, age) VALUES ($1, $2, $3)",
		"tx_user1", "tx1@example.com", 30)
	if err != nil {
		tx.Rollback()
		log.Fatal("插入失败:", err)
	}

	_, err = tx.Exec(ctx, "INSERT INTO users (username, email, age) VALUES ($1, $2, $3)",
		"tx_user2", "tx2@example.com", 25)
	if err != nil {
		tx.Rollback()
		log.Fatal("插入失败:", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		log.Fatal("提交事务失败:", err)
	}
	fmt.Println(" 事务提交成功")

	// 2. 事务回滚
	fmt.Println("\n--- 事务回滚 ---")
	tx2, err := conn.Begin(ctx)
	if err != nil {
		log.Fatal("开始事务失败:", err)
	}

	_, err = tx2.Exec(ctx, "INSERT INTO users (username, email, age) VALUES ($1, $2, $3)",
		"rollback_user", "rollback@example.com", 28)
	if err != nil {
		tx2.Rollback()
		log.Fatal("插入失败:", err)
	}

	// 故意回滚
	if err := tx2.Rollback(); err != nil {
		log.Fatal("回滚失败:", err)
	}
	fmt.Println(" 事务回滚成功")

	// 3. 事务中查询
	fmt.Println("\n--- 事务中查询 ---")
	tx3, err := conn.Begin(ctx)
	if err != nil {
		log.Fatal("开始事务失败:", err)
	}

	var user User
	err = tx3.Query(ctx).
		Table("users").
		Where("username = ?", "tx_user1").
		Scan(&user)

	if err != nil {
		tx3.Rollback()
		if err == postgresql.ErrNoRows {
			fmt.Println("未找到记录")
		} else {
			log.Fatal("查询失败:", err)
		}
	} else {
		fmt.Printf("查询到用户: %+v\n", user)
	}

	tx3.Commit()

	// 清理测试数据
	fmt.Println("\n--- 清理测试数据 ---")
	conn.Exec(ctx, "DELETE FROM users WHERE username LIKE $1", "tx_%")
	fmt.Println(" 清理完成")
}
