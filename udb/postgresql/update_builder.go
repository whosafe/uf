package postgresql

import (
	"context"
	"fmt"
	"strings"

	"github.com/whosafe/uf/uerror"
)

// UpdateBuilder 更新构建器
type UpdateBuilder struct {
	ctx          context.Context
	Connection   *Connection
	table        string
	sets         map[string]any
	whereBuilder WhereBuilder
}

// Table 设置表名
func (b *UpdateBuilder) Table(table string) *UpdateBuilder {
	// 【安全修复】验证表名,防止 SQL 注入
	if err := ValidateIdentifier(table); err != nil {
		panic(err)
	}
	b.table = table
	return b
}

// Set 设置字段值
func (b *UpdateBuilder) Set(column string, value any) *UpdateBuilder {
	// 【安全修复】验证列名,防止 SQL 注入
	if err := ValidateIdentifier(column); err != nil {
		panic(err)
	}
	if b.sets == nil {
		b.sets = make(map[string]any)
	}
	b.sets[column] = value
	return b
}

// SetMap 批量设置字段
func (b *UpdateBuilder) SetMap(data map[string]any) *UpdateBuilder {
	if b.sets == nil {
		b.sets = make(map[string]any)
	}
	for k, v := range data {
		b.sets[k] = v
	}
	return b
}

// Where 添加 WHERE 条件
func (b *UpdateBuilder) Where(condition string, args ...any) *UpdateBuilder {
	b.whereBuilder.Where(condition, args...)
	return b
}

// OrWhere 添加 OR WHERE 条件
func (b *UpdateBuilder) OrWhere(condition string, args ...any) *UpdateBuilder {
	b.whereBuilder.OrWhere(condition, args...)
	return b
}

// WhereIn 添加 IN 条件
func (b *UpdateBuilder) WhereIn(field string, values []any) *UpdateBuilder {
	b.whereBuilder.WhereIn(field, values)
	return b
}

// WhereNotIn 添加 NOT IN 条件
func (b *UpdateBuilder) WhereNotIn(field string, values []any) *UpdateBuilder {
	b.whereBuilder.WhereNotIn(field, values)
	return b
}

// WhereBetween 添加 BETWEEN 条件
func (b *UpdateBuilder) WhereBetween(field string, min, max any) *UpdateBuilder {
	b.whereBuilder.WhereBetween(field, min, max)
	return b
}

// WhereNotBetween 添加 NOT BETWEEN 条件
func (b *UpdateBuilder) WhereNotBetween(field string, min, max any) *UpdateBuilder {
	b.whereBuilder.WhereNotBetween(field, min, max)
	return b
}

// WhereNull 添加 IS NULL 条件
func (b *UpdateBuilder) WhereNull(field string) *UpdateBuilder {
	b.whereBuilder.WhereNull(field)
	return b
}

// WhereNotNull 添加 IS NOT NULL 条件
func (b *UpdateBuilder) WhereNotNull(field string) *UpdateBuilder {
	b.whereBuilder.WhereNotNull(field)
	return b
}

// WhereLike 添加 LIKE 条件
func (b *UpdateBuilder) WhereLike(field string, pattern string) *UpdateBuilder {
	b.whereBuilder.WhereLike(field, pattern)
	return b
}

// Exec 执行更新
func (b *UpdateBuilder) Exec() (int64, error) {
	if b.table == "" {
		return 0, uerror.New("table name is required")
	}
	if len(b.sets) == 0 {
		return 0, uerror.New("no fields to update")
	}

	// 构建 SQL
	var sql strings.Builder
	var args []any
	paramIndex := 1

	sql.WriteString("UPDATE ")
	sql.WriteString(b.table)
	sql.WriteString(" SET ")

	// SET 子句
	setClauses := make([]string, 0, len(b.sets))
	for col, val := range b.sets {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", col, paramIndex))
		args = append(args, val)
		paramIndex++
	}
	sql.WriteString(strings.Join(setClauses, ", "))

	// WHERE 子句
	if b.whereBuilder.HasConditions() {
		whereSQL, whereArgs, _ := buildWhereClause(b.whereBuilder.GetConditions(), paramIndex)
		sql.WriteString(" WHERE ")
		sql.WriteString(whereSQL)
		args = append(args, whereArgs...)
	}

	// 执行
	return b.Connection.Exec(b.ctx, sql.String(), args...)
}
