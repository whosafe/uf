package postgresql

// WhereBuilder WHERE 条件构建器
// 通过组合的方式为各个 Builder 提供统一的 WHERE 方法
type WhereBuilder struct {
	conditions []whereClause
}

// Where 添加 WHERE 条件
func (w *WhereBuilder) Where(condition string, args ...any) *WhereBuilder {
	w.conditions = append(w.conditions, whereClause{
		condition: condition,
		args:      args,
		isOr:      false,
	})
	return w
}

// OrWhere 添加 OR WHERE 条件
func (w *WhereBuilder) OrWhere(condition string, args ...any) *WhereBuilder {
	w.conditions = append(w.conditions, whereClause{
		condition: condition,
		args:      args,
		isOr:      true,
	})
	return w
}

// WhereIn 添加 IN 条件
func (w *WhereBuilder) WhereIn(field string, values []any) *WhereBuilder {
	if len(values) == 0 {
		return w
	}
	w.conditions = append(w.conditions, buildInCondition(field, values, false))
	return w
}

// WhereNotIn 添加 NOT IN 条件
func (w *WhereBuilder) WhereNotIn(field string, values []any) *WhereBuilder {
	if len(values) == 0 {
		return w
	}
	w.conditions = append(w.conditions, buildInCondition(field, values, true))
	return w
}

// WhereBetween 添加 BETWEEN 条件
func (w *WhereBuilder) WhereBetween(field string, min, max any) *WhereBuilder {
	w.conditions = append(w.conditions, buildBetweenCondition(field, min, max, false))
	return w
}

// WhereNotBetween 添加 NOT BETWEEN 条件
func (w *WhereBuilder) WhereNotBetween(field string, min, max any) *WhereBuilder {
	w.conditions = append(w.conditions, buildBetweenCondition(field, min, max, true))
	return w
}

// WhereNull 添加 IS NULL 条件
func (w *WhereBuilder) WhereNull(field string) *WhereBuilder {
	w.conditions = append(w.conditions, buildNullCondition(field, false))
	return w
}

// WhereNotNull 添加 IS NOT NULL 条件
func (w *WhereBuilder) WhereNotNull(field string) *WhereBuilder {
	w.conditions = append(w.conditions, buildNullCondition(field, true))
	return w
}

// WhereLike 添加 LIKE 条件
func (w *WhereBuilder) WhereLike(field string, pattern string) *WhereBuilder {
	w.conditions = append(w.conditions, buildLikeCondition(field, pattern))
	return w
}

// GetConditions 获取所有条件(供内部使用)
func (w *WhereBuilder) GetConditions() []whereClause {
	return w.conditions
}

// HasConditions 检查是否有条件
func (w *WhereBuilder) HasConditions() bool {
	return len(w.conditions) > 0
}
