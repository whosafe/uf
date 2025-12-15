package postgresql

import (
	"testing"
)

// TestQueryBuilder_OrderBy_SQLInjection 测试 OrderBy 防注入
func TestQueryBuilder_OrderBy_SQLInjection(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("应该 panic,但没有")
		}
	}()
	qb := &QueryBuilder{}
	qb.OrderBy("id; DROP TABLE users;")
}

// TestQueryBuilder_OrderByDesc_SQLInjection 测试 OrderByDesc 防注入
func TestQueryBuilder_OrderByDesc_SQLInjection(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("应该 panic,但没有")
		}
	}()
	qb := &QueryBuilder{}
	qb.OrderByDesc("name' OR '1'='1")
}

// TestQueryBuilder_GroupBy_SQLInjection 测试 GroupBy 防注入
func TestQueryBuilder_GroupBy_SQLInjection(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("应该 panic,但没有")
		}
	}()
	qb := &QueryBuilder{}
	qb.GroupBy("id", "name; DELETE FROM users")
}

// TestQueryBuilder_OrderBy_ValidIdentifier 测试 OrderBy 接受有效标识符
func TestQueryBuilder_OrderBy_ValidIdentifier(t *testing.T) {
	qb := &QueryBuilder{}
	// 不应该 panic
	qb.OrderBy("user_id")
	qb.OrderBy("created_at")
	qb.OrderBy("_internal_field")
}

// TestQueryBuilder_GroupBy_ValidIdentifier 测试 GroupBy 接受有效标识符
func TestQueryBuilder_GroupBy_ValidIdentifier(t *testing.T) {
	qb := &QueryBuilder{}
	// 不应该 panic
	qb.GroupBy("category", "status", "user_id")
}
