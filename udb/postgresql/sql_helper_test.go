package postgresql

import (
	"testing"
)

// TestBuildWhereClause_SingleCondition 测试单个条件
func TestBuildWhereClause_SingleCondition(t *testing.T) {
	conditions := []whereClause{
		{condition: "age > ?", args: []any{18}, isOr: false},
	}

	sql, args, paramIndex := buildWhereClause(conditions, 1)

	expectedSQL := "age > $1"
	if sql != expectedSQL {
		t.Errorf("期望 SQL '%s',实际 '%s'", expectedSQL, sql)
	}

	if len(args) != 1 || args[0] != 18 {
		t.Errorf("期望参数 [18],实际 %v", args)
	}

	if paramIndex != 2 {
		t.Errorf("期望 paramIndex 2,实际 %d", paramIndex)
	}
}

// TestBuildWhereClause_MultipleAND 测试多个 AND 条件
func TestBuildWhereClause_MultipleAND(t *testing.T) {
	conditions := []whereClause{
		{condition: "age > ?", args: []any{18}, isOr: false},
		{condition: "status = ?", args: []any{"active"}, isOr: false},
		{condition: "role = ?", args: []any{"user"}, isOr: false},
	}

	sql, args, paramIndex := buildWhereClause(conditions, 1)

	expectedSQL := "age > $1 AND status = $2 AND role = $3"
	if sql != expectedSQL {
		t.Errorf("期望 SQL '%s',实际 '%s'", expectedSQL, sql)
	}

	if len(args) != 3 {
		t.Errorf("期望 3 个参数,实际 %d 个", len(args))
	}

	if paramIndex != 4 {
		t.Errorf("期望 paramIndex 4,实际 %d", paramIndex)
	}
}

// TestBuildWhereClause_MixedANDOR 测试混合 AND/OR 条件
func TestBuildWhereClause_MixedANDOR(t *testing.T) {
	conditions := []whereClause{
		{condition: "age > ?", args: []any{18}, isOr: false},
		{condition: "status = ?", args: []any{"active"}, isOr: false},
		{condition: "role = ?", args: []any{"admin"}, isOr: true},
	}

	sql, args, _ := buildWhereClause(conditions, 1)

	expectedSQL := "age > $1 AND status = $2 OR role = $3"
	if sql != expectedSQL {
		t.Errorf("期望 SQL '%s',实际 '%s'", expectedSQL, sql)
	}

	if len(args) != 3 {
		t.Errorf("期望 3 个参数,实际 %d 个", len(args))
	}
}

// TestBuildWhereClause_EmptyConditions 测试空条件
func TestBuildWhereClause_EmptyConditions(t *testing.T) {
	conditions := []whereClause{}

	sql, args, paramIndex := buildWhereClause(conditions, 1)

	if sql != "" {
		t.Errorf("期望空 SQL,实际 '%s'", sql)
	}

	if len(args) != 0 {
		t.Errorf("期望 0 个参数,实际 %d 个", len(args))
	}

	if paramIndex != 1 {
		t.Errorf("期望 paramIndex 1,实际 %d", paramIndex)
	}
}

// TestBuildWhereClause_MultipleArgs 测试多参数条件
func TestBuildWhereClause_MultipleArgs(t *testing.T) {
	conditions := []whereClause{
		{condition: "age BETWEEN ? AND ?", args: []any{18, 65}, isOr: false},
	}

	sql, args, paramIndex := buildWhereClause(conditions, 1)

	expectedSQL := "age BETWEEN $1 AND $2"
	if sql != expectedSQL {
		t.Errorf("期望 SQL '%s',实际 '%s'", expectedSQL, sql)
	}

	if len(args) != 2 {
		t.Errorf("期望 2 个参数,实际 %d 个", len(args))
	}

	if paramIndex != 3 {
		t.Errorf("期望 paramIndex 3,实际 %d", paramIndex)
	}
}

// TestBuildWhereClause_StartParamIndex 测试自定义起始参数索引
func TestBuildWhereClause_StartParamIndex(t *testing.T) {
	conditions := []whereClause{
		{condition: "age > ?", args: []any{18}, isOr: false},
	}

	sql, _, paramIndex := buildWhereClause(conditions, 5)

	expectedSQL := "age > $5"
	if sql != expectedSQL {
		t.Errorf("期望 SQL '%s',实际 '%s'", expectedSQL, sql)
	}

	if paramIndex != 6 {
		t.Errorf("期望 paramIndex 6,实际 %d", paramIndex)
	}
}

// TestBuildInCondition_Normal 测试正常的 IN 条件
func TestBuildInCondition_Normal(t *testing.T) {
	clause := buildInCondition("id", []any{1, 2, 3}, false)

	if clause.condition != "id IN (?, ?, ?)" {
		t.Errorf("期望条件 'id IN (?, ?, ?)',实际 '%s'", clause.condition)
	}

	if len(clause.args) != 3 {
		t.Errorf("期望 3 个参数,实际 %d 个", len(clause.args))
	}

	if clause.isOr {
		t.Error("应该是 AND 条件")
	}
}

// TestBuildInCondition_NotIn 测试 NOT IN 条件
func TestBuildInCondition_NotIn(t *testing.T) {
	clause := buildInCondition("status", []any{"deleted", "archived"}, true)

	if clause.condition != "status NOT IN (?, ?)" {
		t.Errorf("期望条件 'status NOT IN (?, ?)',实际 '%s'", clause.condition)
	}
}

// TestBuildInCondition_Empty 测试空数组
func TestBuildInCondition_Empty(t *testing.T) {
	clause := buildInCondition("id", []any{}, false)

	if clause.condition != "" {
		t.Error("空数组应该返回空条件")
	}

	if clause.args != nil {
		t.Error("空数组应该返回 nil 参数")
	}
}

// TestBuildBetweenCondition_Normal 测试 BETWEEN 条件
func TestBuildBetweenCondition_Normal(t *testing.T) {
	clause := buildBetweenCondition("age", 18, 65, false)

	if clause.condition != "age BETWEEN ? AND ?" {
		t.Errorf("期望条件 'age BETWEEN ? AND ?',实际 '%s'", clause.condition)
	}

	if len(clause.args) != 2 {
		t.Errorf("期望 2 个参数,实际 %d 个", len(clause.args))
	}

	if clause.args[0] != 18 || clause.args[1] != 65 {
		t.Errorf("期望参数 [18, 65],实际 %v", clause.args)
	}
}

// TestBuildBetweenCondition_NotBetween 测试 NOT BETWEEN 条件
func TestBuildBetweenCondition_NotBetween(t *testing.T) {
	clause := buildBetweenCondition("price", 100, 500, true)

	if clause.condition != "price NOT BETWEEN ? AND ?" {
		t.Errorf("期望条件 'price NOT BETWEEN ? AND ?',实际 '%s'", clause.condition)
	}
}

// TestBuildNullCondition_IsNull 测试 IS NULL 条件
func TestBuildNullCondition_IsNull(t *testing.T) {
	clause := buildNullCondition("deleted_at", false)

	if clause.condition != "deleted_at IS NULL" {
		t.Errorf("期望条件 'deleted_at IS NULL',实际 '%s'", clause.condition)
	}

	if clause.args != nil {
		t.Error("IS NULL 不应该有参数")
	}
}

// TestBuildNullCondition_IsNotNull 测试 IS NOT NULL 条件
func TestBuildNullCondition_IsNotNull(t *testing.T) {
	clause := buildNullCondition("email", true)

	if clause.condition != "email IS NOT NULL" {
		t.Errorf("期望条件 'email IS NOT NULL',实际 '%s'", clause.condition)
	}
}

// TestBuildLikeCondition 测试 LIKE 条件
func TestBuildLikeCondition(t *testing.T) {
	clause := buildLikeCondition("name", "%John%")

	if clause.condition != "name LIKE ?" {
		t.Errorf("期望条件 'name LIKE ?',实际 '%s'", clause.condition)
	}

	if len(clause.args) != 1 || clause.args[0] != "%John%" {
		t.Errorf("期望参数 ['%%John%%'],实际 %v", clause.args)
	}
}
