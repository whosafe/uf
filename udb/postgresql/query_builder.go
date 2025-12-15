package postgresql

import "context"

// QueryBuilder 查询构建器
type QueryBuilder struct {
	ctx        context.Context
	Connection *Connection

	// SELECT 子句
	selectFields []string
	distinct     bool

	// FROM 子句
	table string

	// JOIN 子句
	joins []joinClause

	// WHERE 子句
	whereBuilder WhereBuilder

	// GROUP BY 子句
	groupBy []string

	// HAVING 子句
	having []whereClause

	// ORDER BY 子句
	orderBy []orderClause

	// LIMIT/OFFSET
	limit  int
	offset int
}

// Table 设置表名
func (q *QueryBuilder) Table(name string) *QueryBuilder {
	// 【安全修复】验证表名,防止 SQL 注入
	if err := ValidateTableExpression(name); err != nil {
		panic(err) // 在构建阶段发现错误,直接 panic
	}
	q.table = name
	return q
}

// Select 设置查询字段
func (q *QueryBuilder) Select(fields ...string) *QueryBuilder {
	// 【安全修复】验证字段名,防止 SQL 注入
	if err := ValidateFieldExpressions(fields...); err != nil {
		panic(err)
	}
	q.selectFields = fields
	return q
}

// Distinct 去重
func (q *QueryBuilder) Distinct() *QueryBuilder {
	q.distinct = true
	return q
}

// Where 添加 WHERE 条件
func (q *QueryBuilder) Where(condition string, args ...any) *QueryBuilder {
	q.whereBuilder.Where(condition, args...)
	return q
}

// OrWhere 添加 OR WHERE 条件
func (q *QueryBuilder) OrWhere(condition string, args ...any) *QueryBuilder {
	q.whereBuilder.OrWhere(condition, args...)
	return q
}

// WhereIn 添加 IN 条件
func (q *QueryBuilder) WhereIn(field string, values []any) *QueryBuilder {
	q.whereBuilder.WhereIn(field, values)
	return q
}

// WhereNotIn 添加 NOT IN 条件
func (q *QueryBuilder) WhereNotIn(field string, values []any) *QueryBuilder {
	q.whereBuilder.WhereNotIn(field, values)
	return q
}

// WhereBetween 添加 BETWEEN 条件
func (q *QueryBuilder) WhereBetween(field string, min, max any) *QueryBuilder {
	q.whereBuilder.WhereBetween(field, min, max)
	return q
}

// WhereNotBetween 添加 NOT BETWEEN 条件
func (q *QueryBuilder) WhereNotBetween(field string, min, max any) *QueryBuilder {
	q.whereBuilder.WhereNotBetween(field, min, max)
	return q
}

// WhereNull 添加 IS NULL 条件
func (q *QueryBuilder) WhereNull(field string) *QueryBuilder {
	q.whereBuilder.WhereNull(field)
	return q
}

// WhereNotNull 添加 IS NOT NULL 条件
func (q *QueryBuilder) WhereNotNull(field string) *QueryBuilder {
	q.whereBuilder.WhereNotNull(field)
	return q
}

// WhereLike 添加 LIKE 条件
func (q *QueryBuilder) WhereLike(field string, pattern string) *QueryBuilder {
	q.whereBuilder.WhereLike(field, pattern)
	return q
}

// joinStrings 连接字符串切片(辅助函数)
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}

// Join 添加 INNER JOIN
func (q *QueryBuilder) Join(table, on string) *QueryBuilder {
	q.joins = append(q.joins, joinClause{
		joinType: "INNER",
		table:    table,
		on:       on,
	})
	return q
}

// LeftJoin 左连接
func (q *QueryBuilder) LeftJoin(table, on string) *QueryBuilder {
	q.joins = append(q.joins, joinClause{
		joinType: "LEFT",
		table:    table,
		on:       on,
	})
	return q
}

// RightJoin 右连接
func (q *QueryBuilder) RightJoin(table, on string) *QueryBuilder {
	q.joins = append(q.joins, joinClause{
		joinType: "RIGHT",
		table:    table,
		on:       on,
	})
	return q
}

// FullJoin 全连接
func (q *QueryBuilder) FullJoin(table, on string) *QueryBuilder {
	q.joins = append(q.joins, joinClause{
		joinType: "FULL",
		table:    table,
		on:       on,
	})
	return q
}

// GroupBy 分组
func (q *QueryBuilder) GroupBy(fields ...string) *QueryBuilder {
	// 【安全修复】验证字段名,防止 SQL 注入
	if err := ValidateFieldExpressions(fields...); err != nil {
		panic(err)
	}
	q.groupBy = fields
	return q
}

// Having HAVING 条件
func (q *QueryBuilder) Having(condition string, args ...any) *QueryBuilder {
	q.having = append(q.having, whereClause{
		condition: condition,
		args:      args,
	})
	return q
}

// OrderBy 升序排序
func (q *QueryBuilder) OrderBy(field string) *QueryBuilder {
	// 【安全修复】验证字段名,防止 SQL 注入
	if err := ValidateFieldExpression(field); err != nil {
		panic(err)
	}
	q.orderBy = append(q.orderBy, orderClause{
		field: field,
		desc:  false,
	})
	return q
}

// OrderByDesc 降序排序
func (q *QueryBuilder) OrderByDesc(field string) *QueryBuilder {
	// 【安全修复】验证字段名,防止 SQL 注入
	if err := ValidateFieldExpression(field); err != nil {
		panic(err)
	}
	q.orderBy = append(q.orderBy, orderClause{
		field: field,
		desc:  true,
	})
	return q
}

// Limit 限制数量
func (q *QueryBuilder) Limit(n int) *QueryBuilder {
	q.limit = n
	return q
}

// Offset 偏移量
func (q *QueryBuilder) Offset(n int) *QueryBuilder {
	q.offset = n
	return q
}
