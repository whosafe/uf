package postgresql

import (
	"testing"
)

// TestWhereBuilder_Where 测试基础 WHERE 条件
func TestWhereBuilder_Where(t *testing.T) {
	wb := &WhereBuilder{}
	wb.Where("age > ?", 18)

	if !wb.HasConditions() {
		t.Error("应该有条件")
	}

	conditions := wb.GetConditions()
	if len(conditions) != 1 {
		t.Errorf("期望 1 个条件,实际 %d 个", len(conditions))
	}

	if conditions[0].condition != "age > ?" {
		t.Errorf("期望条件 'age > ?',实际 '%s'", conditions[0].condition)
	}

	if len(conditions[0].args) != 1 || conditions[0].args[0] != 18 {
		t.Errorf("期望参数 [18],实际 %v", conditions[0].args)
	}

	if conditions[0].isOr {
		t.Error("应该是 AND 条件")
	}
}

// TestWhereBuilder_OrWhere 测试 OR WHERE 条件
func TestWhereBuilder_OrWhere(t *testing.T) {
	wb := &WhereBuilder{}
	wb.Where("age > ?", 18)
	wb.OrWhere("status = ?", "active")

	conditions := wb.GetConditions()
	if len(conditions) != 2 {
		t.Errorf("期望 2 个条件,实际 %d 个", len(conditions))
	}

	if !conditions[1].isOr {
		t.Error("第二个条件应该是 OR")
	}
}

// TestWhereBuilder_WhereIn 测试 IN 条件
func TestWhereBuilder_WhereIn(t *testing.T) {
	wb := &WhereBuilder{}
	wb.WhereIn("id", []any{1, 2, 3})

	conditions := wb.GetConditions()
	if len(conditions) != 1 {
		t.Errorf("期望 1 个条件,实际 %d 个", len(conditions))
	}

	if conditions[0].condition != "id IN (?, ?, ?)" {
		t.Errorf("期望条件 'id IN (?, ?, ?)',实际 '%s'", conditions[0].condition)
	}

	if len(conditions[0].args) != 3 {
		t.Errorf("期望 3 个参数,实际 %d 个", len(conditions[0].args))
	}
}

// TestWhereBuilder_WhereIn_Empty 测试空数组的 IN 条件
func TestWhereBuilder_WhereIn_Empty(t *testing.T) {
	wb := &WhereBuilder{}
	wb.WhereIn("id", []any{})

	if wb.HasConditions() {
		t.Error("空数组不应该添加条件")
	}
}

// TestWhereBuilder_WhereNotIn 测试 NOT IN 条件
func TestWhereBuilder_WhereNotIn(t *testing.T) {
	wb := &WhereBuilder{}
	wb.WhereNotIn("status", []any{"deleted", "archived"})

	conditions := wb.GetConditions()
	if len(conditions) != 1 {
		t.Errorf("期望 1 个条件,实际 %d 个", len(conditions))
	}

	if conditions[0].condition != "status NOT IN (?, ?)" {
		t.Errorf("期望条件 'status NOT IN (?, ?)',实际 '%s'", conditions[0].condition)
	}
}

// TestWhereBuilder_WhereBetween 测试 BETWEEN 条件
func TestWhereBuilder_WhereBetween(t *testing.T) {
	wb := &WhereBuilder{}
	wb.WhereBetween("age", 18, 65)

	conditions := wb.GetConditions()
	if len(conditions) != 1 {
		t.Errorf("期望 1 个条件,实际 %d 个", len(conditions))
	}

	if conditions[0].condition != "age BETWEEN ? AND ?" {
		t.Errorf("期望条件 'age BETWEEN ? AND ?',实际 '%s'", conditions[0].condition)
	}

	if len(conditions[0].args) != 2 {
		t.Errorf("期望 2 个参数,实际 %d 个", len(conditions[0].args))
	}
}

// TestWhereBuilder_WhereNotBetween 测试 NOT BETWEEN 条件
func TestWhereBuilder_WhereNotBetween(t *testing.T) {
	wb := &WhereBuilder{}
	wb.WhereNotBetween("price", 100, 500)

	conditions := wb.GetConditions()
	if conditions[0].condition != "price NOT BETWEEN ? AND ?" {
		t.Errorf("期望条件 'price NOT BETWEEN ? AND ?',实际 '%s'", conditions[0].condition)
	}
}

// TestWhereBuilder_WhereNull 测试 IS NULL 条件
func TestWhereBuilder_WhereNull(t *testing.T) {
	wb := &WhereBuilder{}
	wb.WhereNull("deleted_at")

	conditions := wb.GetConditions()
	if conditions[0].condition != "deleted_at IS NULL" {
		t.Errorf("期望条件 'deleted_at IS NULL',实际 '%s'", conditions[0].condition)
	}

	if conditions[0].args != nil {
		t.Error("IS NULL 不应该有参数")
	}
}

// TestWhereBuilder_WhereNotNull 测试 IS NOT NULL 条件
func TestWhereBuilder_WhereNotNull(t *testing.T) {
	wb := &WhereBuilder{}
	wb.WhereNotNull("email")

	conditions := wb.GetConditions()
	if conditions[0].condition != "email IS NOT NULL" {
		t.Errorf("期望条件 'email IS NOT NULL',实际 '%s'", conditions[0].condition)
	}
}

// TestWhereBuilder_WhereLike 测试 LIKE 条件
func TestWhereBuilder_WhereLike(t *testing.T) {
	wb := &WhereBuilder{}
	wb.WhereLike("name", "%John%")

	conditions := wb.GetConditions()
	if conditions[0].condition != "name LIKE ?" {
		t.Errorf("期望条件 'name LIKE ?',实际 '%s'", conditions[0].condition)
	}

	if len(conditions[0].args) != 1 || conditions[0].args[0] != "%John%" {
		t.Errorf("期望参数 ['%%John%%'],实际 %v", conditions[0].args)
	}
}

// TestWhereBuilder_ChainedConditions 测试链式调用多个条件
func TestWhereBuilder_ChainedConditions(t *testing.T) {
	wb := &WhereBuilder{}
	wb.Where("age > ?", 18).
		Where("status = ?", "active").
		OrWhere("role = ?", "admin")

	conditions := wb.GetConditions()
	if len(conditions) != 3 {
		t.Errorf("期望 3 个条件,实际 %d 个", len(conditions))
	}

	// 检查第一个条件
	if conditions[0].isOr {
		t.Error("第一个条件应该是 AND")
	}

	// 检查第二个条件
	if conditions[1].isOr {
		t.Error("第二个条件应该是 AND")
	}

	// 检查第三个条件
	if !conditions[2].isOr {
		t.Error("第三个条件应该是 OR")
	}
}

// TestWhereBuilder_HasConditions 测试 HasConditions 方法
func TestWhereBuilder_HasConditions(t *testing.T) {
	wb := &WhereBuilder{}

	if wb.HasConditions() {
		t.Error("新建的 WhereBuilder 不应该有条件")
	}

	wb.Where("id = ?", 1)

	if !wb.HasConditions() {
		t.Error("添加条件后应该返回 true")
	}
}

// TestWhereBuilder_GetConditions 测试 GetConditions 返回正确的条件列表
func TestWhereBuilder_GetConditions(t *testing.T) {
	wb := &WhereBuilder{}
	wb.Where("a = ?", 1)
	wb.Where("b = ?", 2)

	conditions := wb.GetConditions()

	if len(conditions) != 2 {
		t.Errorf("期望 2 个条件,实际 %d 个", len(conditions))
	}

	// 验证条件顺序
	if conditions[0].condition != "a = ?" {
		t.Error("条件顺序不正确")
	}

	if conditions[1].condition != "b = ?" {
		t.Error("条件顺序不正确")
	}
}
