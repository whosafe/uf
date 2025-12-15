package postgresql

import (
	"context"
	"testing"
)

// TestDeleteBuilder_Table 测试设置表名
func TestDeleteBuilder_Table(t *testing.T) {
	db := &DeleteBuilder{}
	db.Table("users")

	if db.table != "users" {
		t.Errorf("期望表名 'users',实际 '%s'", db.table)
	}
}

// TestDeleteBuilder_Where 测试 WHERE 条件
func TestDeleteBuilder_Where(t *testing.T) {
	db := &DeleteBuilder{}
	db.Where("id = ?", 1)

	if !db.whereBuilder.HasConditions() {
		t.Error("应该有 WHERE 条件")
	}
}

// TestDeleteBuilder_ChainedCalls 测试链式调用
func TestDeleteBuilder_ChainedCalls(t *testing.T) {
	db := &DeleteBuilder{}
	result := db.Table("users").
		Where("id = ?", 1).
		Where("status = ?", "deleted")

	if result != db {
		t.Error("链式调用应该返回同一个对象")
	}
}

// TestDeleteBuilder_Exec_NoTable 测试缺少表名的错误
func TestDeleteBuilder_Exec_NoTable(t *testing.T) {
	db := &DeleteBuilder{ctx: context.Background()}

	_, err := db.Exec()

	if err == nil {
		t.Error("缺少表名应该返回错误")
	}

	if err.Error() != "[0] table name is required" {
		t.Errorf("期望错误 '[0] table name is required',实际 '%s'", err.Error())
	}
}

// TestDeleteBuilder_MultipleWhere 测试多个 WHERE 条件
func TestDeleteBuilder_MultipleWhere(t *testing.T) {
	db := &DeleteBuilder{}
	db.Where("age > ?", 18).
		Where("status = ?", "inactive").
		OrWhere("deleted_at IS NOT NULL")

	conditions := db.whereBuilder.GetConditions()
	if len(conditions) != 3 {
		t.Errorf("期望 3 个条件,实际 %d 个", len(conditions))
	}
}
