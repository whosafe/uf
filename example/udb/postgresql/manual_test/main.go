package main

import (
	"context"
	"fmt"
	"log"
	"strings"
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

// Scan 实现 Scanner 接口
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

// Product 产品结构体
type Product struct {
	ID       int64
	Name     string
	Category string
	Price    float64
	Stock    int
}

// Scan 实现 Scanner 接口
func (p *Product) Scan(key string, value any) error {
	switch key {
	case "id":
		p.ID = uconv.ToInt64Def(value, 0)
	case "name":
		p.Name = uconv.ToString(value)
	case "category":
		p.Category = uconv.ToString(value)
	case "price":
		p.Price = uconv.ToFloat64Def(value, 0)
	case "stock":
		p.Stock = uconv.ToIntDef(value, 0)
	}
	return nil
}

// Order 订单结构体
type Order struct {
	ID         int64
	UserID     int64
	TotalPrice float64
	Status     string
	CreatedAt  time.Time
}

// Scan 实现 Scanner 接口
func (o *Order) Scan(key string, value any) error {
	switch key {
	case "id":
		o.ID = uconv.ToInt64Def(value, 0)
	case "user_id":
		o.UserID = uconv.ToInt64Def(value, 0)
	case "total_price":
		o.TotalPrice = uconv.ToFloat64Def(value, 0)
	case "status":
		o.Status = uconv.ToString(value)
	case "created_at":
		o.CreatedAt = uconv.ToTimeDef(value, time.Time{})
	}
	return nil
}

// UserOrderStats 用户订单统计
type UserOrderStats struct {
	UserID     int64
	Username   string
	OrderCount int64
	TotalPrice float64
}

// Scan 实现 Scanner 接口
func (s *UserOrderStats) Scan(key string, value any) error {
	switch key {
	case "user_id":
		s.UserID = uconv.ToInt64Def(value, 0)
	case "username":
		s.Username = uconv.ToString(value)
	case "order_count":
		s.OrderCount = uconv.ToInt64Def(value, 0)
	case "total_price":
		s.TotalPrice = uconv.ToFloat64Def(value, 0)
	}
	return nil
}

func main() {
	fmt.Println("=== PostgreSQL 完整手动测试 ===")

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
	fmt.Printf("  User: %s\n\n", config.Username)

	// 第一步:连接到 postgres 数据库,检查并创建 testdb
	fmt.Println("=== 步骤 1: 数据库初始化 ===")
	postgresConfig := &postgresql.Config{
		Host:     config.Host,
		Port:     config.Port,
		Username: config.Username,
		Password: config.Password,
		Database: "postgres", // 连接到默认数据库
		SSLMode:  config.SSLMode,
	}

	postgresConn, err := postgresql.New(postgresConfig)
	if err != nil {
		log.Fatal("连接 postgres 数据库失败:", err)
	}

	ctx := ucontext.NewContext(context.Background())

	// 尝试创建 testdb 数据库(如果已存在会报错,但我们可以忽略)
	fmt.Println("初始化数据库...")
	_, err = postgresConn.Exec(ctx, "CREATE DATABASE testdb")
	if err != nil {
		// 如果数据库已存在,错误信息会包含 "already exists"
		// 我们可以忽略这个错误
		if !strings.Contains(err.Error(), "already exists") {
			postgresConn.Close()
			log.Fatal("创建数据库失败:", err)
		}
		fmt.Println("✓ 数据库已存在")
	} else {
		fmt.Println("✓ 数据库创建成功")
	}

	postgresConn.Close()

	// 第二步:连接到 testdb 并创建表结构
	fmt.Println("\n=== 步骤 2: 创建表结构 ===")
	conn, err := postgresql.New(config)
	if err != nil {
		log.Fatal("连接 testdb 失败:", err)
	}
	defer conn.Close()

	// 创建用户表
	fmt.Println("创建 users 表...")
	_, err = conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) NOT NULL UNIQUE,
			email VARCHAR(100) NOT NULL UNIQUE,
			age INTEGER,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal("创建 users 表失败:", err)
	}
	fmt.Println("✓ users 表创建成功")

	// 创建产品表
	fmt.Println("创建 products 表...")
	_, err = conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS products (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			category VARCHAR(50) NOT NULL,
			price DECIMAL(10, 2) NOT NULL,
			stock INTEGER DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal("创建 products 表失败:", err)
	}
	fmt.Println("✓ products 表创建成功")

	// 创建订单表
	fmt.Println("创建 orders 表...")
	_, err = conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			total_price DECIMAL(10, 2) NOT NULL,
			status VARCHAR(20) DEFAULT 'pending',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal("创建 orders 表失败:", err)
	}
	fmt.Println("✓ orders 表创建成功")

	// 第三步:插入测试数据
	fmt.Println("\n=== 步骤 3: 插入测试数据 ===")

	// 插入用户数据
	fmt.Println("插入用户数据...")
	users := []struct {
		username string
		email    string
		age      int
	}{
		{"alice", "alice@example.com", 25},
		{"bob", "bob@example.com", 30},
		{"charlie", "charlie@example.com", 22},
		{"david", "david@example.com", 28},
		{"eve", "eve@example.com", 35},
	}

	for _, u := range users {
		_, err = conn.Insert(ctx).
			Table("users").
			Columns("username", "email", "age").
			Values(u.username, u.email, u.age).
			Exec()
		if err != nil {
			// 忽略重复键错误
			if !strings.Contains(err.Error(), "duplicate key") {
				log.Printf("插入用户 %s 失败: %v\n", u.username, err)
			}
		}
	}
	fmt.Println("✓ 用户数据插入完成")

	// 插入产品数据
	fmt.Println("插入产品数据...")
	products := []struct {
		name     string
		category string
		price    float64
		stock    int
	}{
		{"iPhone 15", "Electronics", 999.99, 50},
		{"MacBook Pro", "Electronics", 2499.99, 30},
		{"AirPods Pro", "Electronics", 249.99, 100},
		{"iPad Air", "Electronics", 599.99, 40},
		{"Apple Watch", "Electronics", 399.99, 60},
		{"Desk Chair", "Furniture", 299.99, 20},
		{"Standing Desk", "Furniture", 599.99, 15},
		{"Monitor", "Electronics", 449.99, 35},
		{"Keyboard", "Electronics", 129.99, 80},
		{"Mouse", "Electronics", 79.99, 100},
	}

	for _, p := range products {
		_, err = conn.Insert(ctx).
			Table("products").
			Columns("name", "category", "price", "stock").
			Values(p.name, p.category, p.price, p.stock).
			Exec()
		if err != nil {
			log.Printf("插入产品 %s 失败: %v\n", p.name, err)
		}
	}
	fmt.Println("✓ 产品数据插入完成")

	// 插入订单数据
	fmt.Println("插入订单数据...")
	orders := []struct {
		userID     int64
		totalPrice float64
		status     string
	}{
		{1, 1249.98, "completed"},
		{2, 2499.99, "completed"},
		{3, 329.98, "pending"},
		{1, 999.99, "completed"},
		{4, 899.98, "shipped"},
	}

	for _, o := range orders {
		_, err = conn.Insert(ctx).
			Table("orders").
			Columns("user_id", "total_price", "status").
			Values(o.userID, o.totalPrice, o.status).
			Exec()
		if err != nil {
			log.Printf("插入订单失败: %v\n", err)
		}
	}
	fmt.Println("✓ 订单数据插入完成")

	// 第四步:基础查询测试
	fmt.Println("\n=== 步骤 4: 基础查询测试 ===")

	// 测试 Ping
	fmt.Println("\n1. 测试连接 Ping")
	if err := conn.Ping(ctx); err != nil {
		log.Fatal("Ping 失败:", err)
	}
	fmt.Println("✓ Ping 成功")

	// 查询单条记录
	fmt.Println("\n2. 查询单条用户记录")
	var user User
	err = conn.Query(ctx).
		Table("users").
		Where("username = ?", "alice").
		Scan(&user)

	if err != nil {
		if err == postgresql.ErrNoRows {
			fmt.Println("✗ 未找到用户")
		} else {
			log.Printf("查询失败: %v\n", err)
		}
	} else {
		fmt.Printf("✓ 找到用户: %s (%s), 年龄: %d\n", user.Username, user.Email, user.Age)
	}

	// 查询多条记录
	fmt.Println("\n3. 查询多条产品记录")
	productResults, err := conn.Query(ctx).
		Table("products").
		Where("category = ?", "Electronics").
		Where("price < ?", 1000).
		OrderBy("price DESC").
		Limit(5).
		ScanAll(func() postgresql.Scanner { return &Product{} })

	if err != nil {
		log.Printf("查询失败: %v\n", err)
	} else {
		fmt.Printf("✓ 找到 %d 个产品:\n", len(productResults))
		for _, r := range productResults {
			p := r.(*Product)
			fmt.Printf("  - %s: $%.2f (库存: %d)\n", p.Name, p.Price, p.Stock)
		}
	}

	// 第五步:查询构建器测试
	fmt.Println("\n=== 步骤 5: 查询构建器测试 ===")

	// JOIN 查询
	fmt.Println("\n1. JOIN 查询 - 用户订单统计")
	statsResults, err := conn.Query(ctx).
		Select("u.id as user_id", "u.username", "COUNT(o.id) as order_count", "COALESCE(SUM(o.total_price), 0) as total_price").
		Table("users u").
		LeftJoin("orders o", "u.id = o.user_id").
		GroupBy("u.id", "u.username").
		Having("COUNT(o.id) > ?", 0).
		OrderByDesc("order_count").
		ScanAll(func() postgresql.Scanner { return &UserOrderStats{} })

	if err != nil {
		log.Printf("查询失败: %v\n", err)
	} else {
		fmt.Printf("✓ 找到 %d 个有订单的用户:\n", len(statsResults))
		for _, r := range statsResults {
			s := r.(*UserOrderStats)
			fmt.Printf("  - %s: %d 个订单, 总金额 $%.2f\n", s.Username, s.OrderCount, s.TotalPrice)
		}
	}

	// DISTINCT 查询
	fmt.Println("\n2. DISTINCT 查询 - 产品分类")
	categoryResults, err := conn.Query(ctx).
		Select("category").
		Table("products").
		Distinct().
		OrderBy("category").
		ScanAll(func() postgresql.Scanner { return &Product{} })

	if err != nil {
		log.Printf("查询失败: %v\n", err)
	} else {
		fmt.Printf("✓ 找到 %d 个产品分类:\n", len(categoryResults))
		for _, r := range categoryResults {
			p := r.(*Product)
			fmt.Printf("  - %s\n", p.Category)
		}
	}

	// 第六步:插入、更新、删除测试
	fmt.Println("\n=== 步骤 6: CRUD 操作测试 ===")

	// 插入测试
	fmt.Println("\n1. 插入新用户")
	affected, err := conn.Insert(ctx).
		Table("users").
		Columns("username", "email", "age").
		Values("test_user", "test@example.com", 27).
		Exec()
	if err != nil {
		log.Printf("✗ 插入失败: %v\n", err)
	} else {
		fmt.Printf("✓ 插入成功, 影响行数: %d\n", affected)
	}

	// 更新测试
	fmt.Println("\n2. 更新用户年龄")
	affected, err = conn.Update(ctx).
		Table("users").
		Set("age", 28).
		Where("username = ?", "test_user").
		Exec()
	if err != nil {
		log.Printf("✗ 更新失败: %v\n", err)
	} else {
		fmt.Printf("✓ 更新成功, 影响行数: %d\n", affected)
	}

	// 查询验证更新
	var testUser User
	err = conn.Query(ctx).
		Table("users").
		Where("username = ?", "test_user").
		Scan(&testUser)
	if err == nil {
		fmt.Printf("✓ 验证更新: 用户年龄为 %d\n", testUser.Age)
	}

	// 删除测试
	fmt.Println("\n3. 删除测试用户")
	affected, err = conn.Delete(ctx).
		Table("users").
		Where("username = ?", "test_user").
		Exec()
	if err != nil {
		log.Printf("✗ 删除失败: %v\n", err)
	} else {
		fmt.Printf("✓ 删除成功, 影响行数: %d\n", affected)
	}

	// 第七步:事务测试
	fmt.Println("\n=== 步骤 7: 事务测试 ===")

	// 事务提交测试
	fmt.Println("\n1. 事务提交测试")
	tx, err := conn.Begin(ctx)
	if err != nil {
		log.Fatal("开始事务失败:", err)
	}

	_, err = tx.Insert(ctx).
		Table("users").
		Columns("username", "email", "age").
		Values("tx_user1", "tx1@example.com", 30).
		Exec()
	if err != nil {
		tx.Rollback()
		log.Printf("✗ 事务插入失败: %v\n", err)
	} else {
		if err := tx.Commit(); err != nil {
			log.Printf("✗ 事务提交失败: %v\n", err)
		} else {
			fmt.Println("✓ 事务提交成功")
		}
	}

	// 事务回滚测试
	fmt.Println("\n2. 事务回滚测试")
	tx2, err := conn.Begin(ctx)
	if err != nil {
		log.Fatal("开始事务失败:", err)
	}

	_, err = tx2.Insert(ctx).
		Table("users").
		Columns("username", "email", "age").
		Values("rollback_user", "rollback@example.com", 25).
		Exec()
	if err != nil {
		tx2.Rollback()
		log.Printf("✗ 事务插入失败: %v\n", err)
	} else {
		// 故意回滚
		if err := tx2.Rollback(); err != nil {
			log.Printf("✗ 事务回滚失败: %v\n", err)
		} else {
			fmt.Println("✓ 事务回滚成功")
		}
	}

	// 验证回滚
	var rollbackUser User
	err = conn.Query(ctx).
		Table("users").
		Where("username = ?", "rollback_user").
		Scan(&rollbackUser)
	if err == postgresql.ErrNoRows {
		fmt.Println("✓ 验证回滚: 用户不存在(符合预期)")
	} else if err != nil {
		log.Printf("查询失败: %v\n", err)
	} else {
		fmt.Println("✗ 验证回滚失败: 用户仍然存在")
	}

	// 事务中查询测试
	fmt.Println("\n3. 事务中查询测试")
	tx3, err := conn.Begin(ctx)
	if err != nil {
		log.Fatal("开始事务失败:", err)
	}

	var txUser User
	err = tx3.Query(ctx).
		Table("users").
		Where("username = ?", "tx_user1").
		Scan(&txUser)

	if err != nil {
		tx3.Rollback()
		if err == postgresql.ErrNoRows {
			fmt.Println("✗ 未找到用户")
		} else {
			log.Printf("✗ 查询失败: %v\n", err)
		}
	} else {
		fmt.Printf("✓ 在事务中查询到用户: %s (%s)\n", txUser.Username, txUser.Email)
		tx3.Commit()
	}

	// 第八步:连接池统计
	fmt.Println("\n=== 步骤 8: 连接池统计 ===")
	stats := conn.Stats()
	fmt.Printf("总连接数: %d\n", stats.TotalConns())
	fmt.Printf("空闲连接数: %d\n", stats.IdleConns())
	fmt.Printf("获取连接数: %d\n", stats.AcquiredConns())

	// 清理测试数据
	fmt.Println("\n=== 清理测试数据 ===")
	fmt.Println("删除测试用户...")
	conn.Delete(ctx).
		Table("users").
		WhereLike("username", "tx_%").
		Exec()
	fmt.Println("✓ 清理完成")

	fmt.Println("\n=== ✓ 所有测试完成! ===")
}
