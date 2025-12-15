package postgresql

import (
	"context"
	"testing"
)

// TestInsertBuilder_Table 测试设置表名
func TestInsertBuilder_Table(t *testing.T) {
	ib := &InsertBuilder{}
	ib.Table("users")

	if ib.table != "users" {
		t.Errorf("期望表名 'users',实际 '%s'", ib.table)
	}
}

// TestInsertBuilder_Columns 测试设置列名
func TestInsertBuilder_Columns(t *testing.T) {
	ib := &InsertBuilder{}
	ib.Columns("name", "age", "email")

	if len(ib.columns) != 3 {
		t.Errorf("期望 3 个列,实际 %d 个", len(ib.columns))
	}

	if ib.columns[0] != "name" {
		t.Error("列顺序不正确")
	}
}

// TestInsertBuilder_Values 测试设置值
func TestInsertBuilder_Values(t *testing.T) {
	ib := &InsertBuilder{}
	ib.Values("John", 25, "john@example.com")

	if len(ib.values) != 3 {
		t.Errorf("期望 3 个值,实际 %d 个", len(ib.values))
	}

	if ib.values[0] != "John" {
		t.Error("值顺序不正确")
	}
}

// TestInsertBuilder_ChainedCalls 测试链式调用
func TestInsertBuilder_ChainedCalls(t *testing.T) {
	ib := &InsertBuilder{}
	result := ib.Table("users").
		Columns("name", "age").
		Values("John", 25)

	if result != ib {
		t.Error("链式调用应该返回同一个对象")
	}
}

// TestInsertBuilder_Exec_NoTable 测试缺少表名的错误
func TestInsertBuilder_Exec_NoTable(t *testing.T) {
	ib := &InsertBuilder{ctx: context.Background()}
	ib.Columns("name").Values("John")

	_, err := ib.Exec()

	if err == nil {
		t.Error("缺少表名应该返回错误")
	}

	if err.Error() != "[0] table name is required" {
		t.Errorf("期望错误 '[0] table name is required',实际 '%s'", err.Error())
	}
}

// TestInsertBuilder_Exec_NoColumns 测试缺少列名的错误
func TestInsertBuilder_Exec_NoColumns(t *testing.T) {
	ib := &InsertBuilder{ctx: context.Background()}
	ib.Table("users").Values("John")

	_, err := ib.Exec()

	if err == nil {
		t.Error("缺少列名应该返回错误")
	}

	if err.Error() != "[0] columns are required" {
		t.Errorf("期望错误 '[0] columns are required',实际 '%s'", err.Error())
	}
}

// TestInsertBuilder_Exec_NoValues 测试缺少值的错误
func TestInsertBuilder_Exec_NoValues(t *testing.T) {
	ib := &InsertBuilder{ctx: context.Background()}
	ib.Table("users").Columns("name")

	_, err := ib.Exec()

	if err == nil {
		t.Error("缺少值应该返回错误")
	}

	if err.Error() != "[0] values are required" {
		t.Errorf("期望错误 '[0] values are required',实际 '%s'", err.Error())
	}
}

// TestInsertBuilder_Exec_MismatchCount 测试列值数量不匹配的错误
func TestInsertBuilder_Exec_MismatchCount(t *testing.T) {
	ib := &InsertBuilder{ctx: context.Background()}
	ib.Table("users").
		Columns("name", "age").
		Values("John") // 只有 1 个值,但有 2 个列

	_, err := ib.Exec()

	if err == nil {
		t.Error("列值数量不匹配应该返回错误")
	}

	if err.Error() != "[0] columns and values count mismatch" {
		t.Errorf("期望错误 '[0] columns and values count mismatch',实际 '%s'", err.Error())
	}
}
