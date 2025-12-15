package main

import (
	"context"
	"fmt"
	"log"

	"github.com/whosafe/uf/ucontext"
	"github.com/whosafe/uf/uconv"
	"github.com/whosafe/uf/udb/postgresql"
)

// Version 数据库版本
type Version struct {
	Version string
}

// Scan 实现 Scanner 接口
func (v *Version) Scan(key string, value any) error {
	switch key {
	case "version":
		v.Version = uconv.ToString(value)
	}
	return nil
}

// CurrentDB 当前数据库
type CurrentDB struct {
	Database string
}

// Scan 实现 Scanner 接口
func (d *CurrentDB) Scan(key string, value any) error {
	switch key {
	case "current_database":
		d.Database = uconv.ToString(value)
	}
	return nil
}

func main() {
	fmt.Println("=== PostgreSQL 连接测试 ===")

	// 创建配置
	config := &postgresql.Config{
		Host:     "10.0.2.128",
		Port:     5432,
		Username: "postgres",
		Password: "111111",
		Database: "testdb",
		SSLMode:  "disable",
	}

	fmt.Printf("连接配置:\n")
	fmt.Printf("  Host: %s:%d\n", config.Host, config.Port)
	fmt.Printf("  Database: %s\n", config.Database)
	fmt.Printf("  User: %s\n", config.Username)
	fmt.Printf("  SSL Mode: %s\n\n", config.SSLMode)

	// 创建连接
	fmt.Println("正在创建连接...")
	conn, err := postgresql.New(config)
	if err != nil {
		log.Fatal("连接失败:", err)
	}
	defer conn.Close()
	fmt.Println(" 连接创建成功")

	// 创建追踪上下文
	ctx := ucontext.NewContext(context.Background())

	// 测试 Ping
	fmt.Println("正在测试连接...")
	if err := conn.Ping(ctx); err != nil {
		log.Fatal("Ping 失败:", err)
	}
	fmt.Println(" Ping 成功")

	// 查询数据库版本
	fmt.Println("查询数据库版本...")
	var version Version
	err = conn.Query(ctx).
		Select("version()").
		Table("(SELECT version() as version) as v").
		Scan(&version)

	if err != nil {
		log.Printf("查询版本失败: %v\n", err)
	} else {
		fmt.Printf(" 数据库版本: %s\n\n", version.Version)
	}

	// 查询当前数据库
	fmt.Println("查询当前数据库...")
	var currentDB CurrentDB
	err = conn.Query(ctx).
		Select("current_database()").
		Table("(SELECT current_database()) as d").
		Scan(&currentDB)

	if err != nil {
		log.Printf("查询当前数据库失败: %v\n", err)
	} else {
		fmt.Printf(" 当前数据库: %s\n\n", currentDB.Database)
	}

	// 连接池统计
	fmt.Println("连接池统计:")
	stats := conn.Stats()
	fmt.Printf("  总连接数: %d\n", stats.TotalConns())
	fmt.Printf("  空闲连接数: %d\n", stats.IdleConns())
	fmt.Printf("  获取连接数: %d\n", stats.AcquiredConns())

	fmt.Println("\n 所有测试通过!")
}
