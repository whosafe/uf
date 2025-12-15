package postgresql

import (
	"context"
	"testing"
)

// TestUpdateBuilder_Table 测试设置表名
func TestUpdateBuilder_Table(t *testing.T) {
	ub := &UpdateBuilder{}
	ub.Table("users")

	if ub.table != "users" {
		t.Errorf("期望表名 'users',实际 '%s'", ub.table)
	}
}

// TestUpdateBuilder_Set 测试设置单个字段
func TestUpdateBuilder_Set(t *testing.T) {
	ub := &UpdateBuilder{}
	ub.Set("name", "John")

	if len(ub.sets) != 1 {
		t.Errorf("期望 1 个字段,实际 %d 个", len(ub.sets))
	}

	if ub.sets["name"] != "John" {
		t.Errorf("期望值 'John',实际 '%v'", ub.sets["name"])
	}
}

// TestUpdateBuilder_Set_Multiple 测试设置多个字段
func TestUpdateBuilder_Set_Multiple(t *testing.T) {
	ub := &UpdateBuilder{}
	ub.Set("name", "John").Set("age", 25)

	if len(ub.sets) != 2 {
		t.Errorf("期望 2 个字段,实际 %d 个", len(ub.sets))
	}
}

// TestUpdateBuilder_SetMap 测试批量设置字段
func TestUpdateBuilder_SetMap(t *testing.T) {
	ub := &UpdateBuilder{}
	data := map[string]any{
		"name":  "John",
		"age":   25,
		"email": "john@example.com",
	}
	ub.SetMap(data)

	if len(ub.sets) != 3 {
		t.Errorf("期望 3 个字段,实际 %d 个", len(ub.sets))
	}

	if ub.sets["name"] != "John" {
		t.Error("SetMap 应该正确设置字段")
	}
}

// TestUpdateBuilder_Where 测试 WHERE 条件
func TestUpdateBuilder_Where(t *testing.T) {
	ub := &UpdateBuilder{}
	ub.Where("id = ?", 1)

	if !ub.whereBuilder.HasConditions() {
		t.Error("应该有 WHERE 条件")
	}
}

// TestUpdateBuilder_ChainedCalls 测试链式调用
func TestUpdateBuilder_ChainedCalls(t *testing.T) {
	ub := &UpdateBuilder{}
	result := ub.Table("users").
		Set("name", "John").
		Set("age", 25).
		Where("id = ?", 1)

	if result != ub {
		t.Error("链式调用应该返回同一个对象")
	}
}

// TestUpdateBuilder_Exec_NoTable 测试缺少表名的错误
func TestUpdateBuilder_Exec_NoTable(t *testing.T) {
	ub := &UpdateBuilder{ctx: context.Background()}
	ub.Set("name", "John")

	_, err := ub.Exec()

	if err == nil {
		t.Error("缺少表名应该返回错误")
	}

	if err.Error() != "[0] table name is required" {
		t.Errorf("期望错误 '[0] table name is required',实际 '%s'", err.Error())
	}
}

// TestUpdateBuilder_Exec_NoFields 测试缺少字段的错误
func TestUpdateBuilder_Exec_NoFields(t *testing.T) {
	ub := &UpdateBuilder{ctx: context.Background()}
	ub.Table("users")

	_, err := ub.Exec()

	if err == nil {
		t.Error("缺少字段应该返回错误")
	}

	if err.Error() != "[0] no fields to update" {
		t.Errorf("期望错误 '[0] no fields to update',实际 '%s'", err.Error())
	}
}
