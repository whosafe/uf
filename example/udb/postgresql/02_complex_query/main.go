package main

import (
	"context"
	"fmt"
	"log"

	"github.com/whosafe/uf/ucontext"
	"github.com/whosafe/uf/udb/postgresql"
)

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
		if v, ok := value.(int64); ok {
			s.UserID = v
		}
	case "username":
		if v, ok := value.(string); ok {
			s.Username = v
		}
	case "order_count":
		if v, ok := value.(int64); ok {
			s.OrderCount = v
		}
	case "total_price":
		if v, ok := value.(float64); ok {
			s.TotalPrice = v
		}
	}
	return nil
}

// CategoryStats 分类统计
type CategoryStats struct {
	Category string
	Count    int64
	AvgPrice float64
}

// Scan 实现 Scanner 接口
func (c *CategoryStats) Scan(key string, value any) error {
	switch key {
	case "category":
		if v, ok := value.(string); ok {
			c.Category = v
		}
	case "count":
		if v, ok := value.(int64); ok {
			c.Count = v
		}
	case "avg_price":
		if v, ok := value.(float64); ok {
			c.AvgPrice = v
		}
	}
	return nil
}

// Product 产品
type Product struct {
	ID       int64
	Name     string
	Category string
	Price    float64
}

// Scan 实现 Scanner 接口
func (p *Product) Scan(key string, value any) error {
	switch key {
	case "id":
		if v, ok := value.(int64); ok {
			p.ID = v
		}
	case "name":
		if v, ok := value.(string); ok {
			p.Name = v
		}
	case "category":
		if v, ok := value.(string); ok {
			p.Category = v
		}
	case "price":
		if v, ok := value.(float64); ok {
			p.Price = v
		}
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

	// 1. LEFT JOIN + GROUP BY + HAVING
	fmt.Println("=== LEFT JOIN + GROUP BY + HAVING ===")
	statsResults, err := conn.Query(ctx).
		Select("u.id as user_id", "u.username", "COUNT(o.id) as order_count", "COALESCE(SUM(o.total_price), 0) as total_price").
		Table("users u").
		LeftJoin("orders o", "u.id = o.user_id").
		GroupBy("u.id", "u.username").
		Having("COUNT(o.id) > ?", 0).
		OrderByDesc("order_count").
		Limit(10).
		ScanAll(func() postgresql.Scanner { return &UserOrderStats{} })

	if err != nil {
		log.Printf("查询失败: %v\n", err)
	} else {
		fmt.Printf("找到 %d 个用户:\n", len(statsResults))
		for _, r := range statsResults {
			s := r.(*UserOrderStats)
			fmt.Printf("  - %s: %d 个订单, 总价 %.2f\n", s.Username, s.OrderCount, s.TotalPrice)
		}
	}

	// 2. GROUP BY + HAVING + 聚合函数
	fmt.Println("\n=== GROUP BY + HAVING + 聚合函数 ===")
	categoryResults, err := conn.Query(ctx).
		Select("category", "COUNT(*) as count", "AVG(price) as avg_price").
		Table("products").
		GroupBy("category").
		Having("AVG(price) > ?", 100).
		OrderByDesc("avg_price").
		ScanAll(func() postgresql.Scanner { return &CategoryStats{} })

	if err != nil {
		log.Printf("查询失败: %v\n", err)
	} else {
		fmt.Printf("找到 %d 个分类:\n", len(categoryResults))
		for _, r := range categoryResults {
			c := r.(*CategoryStats)
			fmt.Printf("  - %s: %d 个产品, 平均价格 %.2f\n", c.Category, c.Count, c.AvgPrice)
		}
	}

	// 3. DISTINCT
	fmt.Println("\n=== DISTINCT ===")
	distinctResults, err := conn.Query(ctx).
		Select("category").
		Table("products").
		Distinct().
		OrderBy("category").
		ScanAll(func() postgresql.Scanner { return &Product{} })

	if err != nil {
		log.Printf("查询失败: %v\n", err)
	} else {
		fmt.Printf("找到 %d 个不同的分类:\n", len(distinctResults))
		for _, r := range distinctResults {
			p := r.(*Product)
			fmt.Printf("  - %s\n", p.Category)
		}
	}

	// 4. 多条件 WHERE
	fmt.Println("\n=== 多条件 WHERE ===")
	productResults, err := conn.Query(ctx).
		Table("products").
		Where("category = ?", "Electronics").
		Where("price > ?", 500).
		Where("price < ?", 2000).
		OrderBy("price").
		Limit(5).
		ScanAll(func() postgresql.Scanner { return &Product{} })

	if err != nil {
		log.Printf("查询失败: %v\n", err)
	} else {
		fmt.Printf("找到 %d 个产品:\n", len(productResults))
		for _, r := range productResults {
			p := r.(*Product)
			fmt.Printf("  - %s (%.2f)\n", p.Name, p.Price)
		}
	}
}
