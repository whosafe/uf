package postgresql

import (
	"testing"
)

// TestQueryBuilder_Table 测试设置表名
func TestQueryBuilder_Table(t *testing.T) {
	qb := &QueryBuilder{}
	qb.Table("users")

	if qb.table != "users" {
		t.Errorf("期望表名 'users',实际 '%s'", qb.table)
	}
}

// TestQueryBuilder_Select 测试选择字段
func TestQueryBuilder_Select(t *testing.T) {
	qb := &QueryBuilder{}
	qb.Select("id", "name", "email")

	if len(qb.selectFields) != 3 {
		t.Errorf("期望 3 个字段,实际 %d 个", len(qb.selectFields))
	}

	if qb.selectFields[0] != "id" {
		t.Error("字段顺序不正确")
	}
}

// TestQueryBuilder_Distinct 测试去重
func TestQueryBuilder_Distinct(t *testing.T) {
	qb := &QueryBuilder{}
	qb.Distinct()

	if !qb.distinct {
		t.Error("应该启用 DISTINCT")
	}
}

// TestQueryBuilder_Where 测试 WHERE 条件
func TestQueryBuilder_Where(t *testing.T) {
	qb := &QueryBuilder{}
	qb.Where("age > ?", 18)

	if !qb.whereBuilder.HasConditions() {
		t.Error("应该有 WHERE 条件")
	}
}

// TestQueryBuilder_Join 测试 JOIN
func TestQueryBuilder_Join(t *testing.T) {
	qb := &QueryBuilder{}
	qb.Join("orders", "users.id = orders.user_id")

	if len(qb.joins) != 1 {
		t.Errorf("期望 1 个 JOIN,实际 %d 个", len(qb.joins))
	}

	if qb.joins[0].joinType != "INNER" {
		t.Errorf("期望 JOIN 类型 'INNER',实际 '%s'", qb.joins[0].joinType)
	}
}

// TestQueryBuilder_LeftJoin 测试 LEFT JOIN
func TestQueryBuilder_LeftJoin(t *testing.T) {
	qb := &QueryBuilder{}
	qb.LeftJoin("orders", "users.id = orders.user_id")

	if qb.joins[0].joinType != "LEFT" {
		t.Errorf("期望 JOIN 类型 'LEFT',实际 '%s'", qb.joins[0].joinType)
	}
}

// TestQueryBuilder_GroupBy 测试 GROUP BY
func TestQueryBuilder_GroupBy(t *testing.T) {
	qb := &QueryBuilder{}
	qb.GroupBy("category", "status")

	if len(qb.groupBy) != 2 {
		t.Errorf("期望 2 个 GROUP BY 字段,实际 %d 个", len(qb.groupBy))
	}
}

// TestQueryBuilder_Having 测试 HAVING
func TestQueryBuilder_Having(t *testing.T) {
	qb := &QueryBuilder{}
	qb.Having("COUNT(*) > ?", 5)

	if len(qb.having) != 1 {
		t.Errorf("期望 1 个 HAVING 条件,实际 %d 个", len(qb.having))
	}
}

// TestQueryBuilder_OrderBy 测试升序排序
func TestQueryBuilder_OrderBy(t *testing.T) {
	qb := &QueryBuilder{}
	qb.OrderBy("created_at")

	if len(qb.orderBy) != 1 {
		t.Errorf("期望 1 个 ORDER BY,实际 %d 个", len(qb.orderBy))
	}

	if qb.orderBy[0].desc {
		t.Error("应该是升序")
	}
}

// TestQueryBuilder_OrderByDesc 测试降序排序
func TestQueryBuilder_OrderByDesc(t *testing.T) {
	qb := &QueryBuilder{}
	qb.OrderByDesc("created_at")

	if !qb.orderBy[0].desc {
		t.Error("应该是降序")
	}
}

// TestQueryBuilder_Limit 测试 LIMIT
func TestQueryBuilder_Limit(t *testing.T) {
	qb := &QueryBuilder{}
	qb.Limit(10)

	if qb.limit != 10 {
		t.Errorf("期望 LIMIT 10,实际 %d", qb.limit)
	}
}

// TestQueryBuilder_Offset 测试 OFFSET
func TestQueryBuilder_Offset(t *testing.T) {
	qb := &QueryBuilder{}
	qb.Offset(20)

	if qb.offset != 20 {
		t.Errorf("期望 OFFSET 20,实际 %d", qb.offset)
	}
}

// TestQueryBuilder_BuildSQL_Simple 测试简单查询的 SQL 构建
func TestQueryBuilder_BuildSQL_Simple(t *testing.T) {
	qb := &QueryBuilder{}
	qb.Table("users")

	sql, args := qb.BuildSQL()

	expectedSQL := "SELECT * FROM users"
	if sql != expectedSQL {
		t.Errorf("期望 SQL '%s',实际 '%s'", expectedSQL, sql)
	}

	if len(args) != 0 {
		t.Errorf("期望 0 个参数,实际 %d 个", len(args))
	}
}

// TestQueryBuilder_BuildSQL_WithSelect 测试带 SELECT 字段的查询
func TestQueryBuilder_BuildSQL_WithSelect(t *testing.T) {
	qb := &QueryBuilder{}
	qb.Select("id", "name").Table("users")

	sql, _ := qb.BuildSQL()

	expectedSQL := "SELECT id, name FROM users"
	if sql != expectedSQL {
		t.Errorf("期望 SQL '%s',实际 '%s'", expectedSQL, sql)
	}
}

// TestQueryBuilder_BuildSQL_WithDistinct 测试带 DISTINCT 的查询
func TestQueryBuilder_BuildSQL_WithDistinct(t *testing.T) {
	qb := &QueryBuilder{}
	qb.Select("category").Table("products").Distinct()

	sql, _ := qb.BuildSQL()

	expectedSQL := "SELECT DISTINCT category FROM products"
	if sql != expectedSQL {
		t.Errorf("期望 SQL '%s',实际 '%s'", expectedSQL, sql)
	}
}

// TestQueryBuilder_BuildSQL_WithWhere 测试带 WHERE 的查询
func TestQueryBuilder_BuildSQL_WithWhere(t *testing.T) {
	qb := &QueryBuilder{}
	qb.Table("users").Where("age > ?", 18)

	sql, args := qb.BuildSQL()

	expectedSQL := "SELECT * FROM users WHERE age > $1"
	if sql != expectedSQL {
		t.Errorf("期望 SQL '%s',实际 '%s'", expectedSQL, sql)
	}

	if len(args) != 1 || args[0] != 18 {
		t.Errorf("期望参数 [18],实际 %v", args)
	}
}

// TestQueryBuilder_BuildSQL_Complex 测试复杂查询
func TestQueryBuilder_BuildSQL_Complex(t *testing.T) {
	qb := &QueryBuilder{}
	qb.Select("u.id", "u.name", "COUNT(o.id) as order_count").
		Table("users u").
		LeftJoin("orders o", "u.id = o.user_id").
		Where("u.status = ?", "active").
		GroupBy("u.id", "u.name").
		Having("COUNT(o.id) > ?", 0).
		OrderByDesc("order_count").
		Limit(10).
		Offset(20)

	sql, args := qb.BuildSQL()

	// 验证 SQL 包含所有关键部分
	if !contains(sql, "SELECT u.id, u.name, COUNT(o.id) as order_count") {
		t.Error("SQL 应该包含 SELECT 字段")
	}
	if !contains(sql, "FROM users u") {
		t.Error("SQL 应该包含 FROM")
	}
	if !contains(sql, "LEFT JOIN orders o ON u.id = o.user_id") {
		t.Error("SQL 应该包含 LEFT JOIN")
	}
	if !contains(sql, "WHERE u.status = $1") {
		t.Error("SQL 应该包含 WHERE")
	}
	if !contains(sql, "GROUP BY u.id, u.name") {
		t.Error("SQL 应该包含 GROUP BY")
	}
	if !contains(sql, "HAVING COUNT(o.id) > $2") {
		t.Error("SQL 应该包含 HAVING")
	}
	if !contains(sql, "ORDER BY order_count DESC") {
		t.Error("SQL 应该包含 ORDER BY")
	}
	if !contains(sql, "LIMIT 10") {
		t.Error("SQL 应该包含 LIMIT")
	}
	if !contains(sql, "OFFSET 20") {
		t.Error("SQL 应该包含 OFFSET")
	}

	if len(args) != 2 {
		t.Errorf("期望 2 个参数,实际 %d 个", len(args))
	}
}

// TestQueryBuilder_ChainedCalls 测试链式调用
func TestQueryBuilder_ChainedCalls(t *testing.T) {
	qb := &QueryBuilder{}
	result := qb.Table("users").
		Select("id", "name").
		Where("age > ?", 18).
		OrderBy("name").
		Limit(10)

	if result != qb {
		t.Error("链式调用应该返回同一个对象")
	}
}

// 辅助函数:检查字符串是否包含子串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
