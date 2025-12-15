package postgresql

import (
	"context"
	"strings"

	"github.com/whosafe/uf/uerror"
)

// DeleteBuilder 删除构建器
type DeleteBuilder struct {
	ctx          context.Context
	Connection   *Connection
	table        string
	whereBuilder WhereBuilder
}

// Table 设置表名
func (b *DeleteBuilder) Table(table string) *DeleteBuilder {
	// 【安全修复】验证表名,防止 SQL 注入
	if err := ValidateIdentifier(table); err != nil {
		panic(err)
	}
	b.table = table
	return b
}

// Where 添加 WHERE 条件
func (b *DeleteBuilder) Where(condition string, args ...any) *DeleteBuilder {
	b.whereBuilder.Where(condition, args...)
	return b
}

// OrWhere 添加 OR WHERE 条件
func (b *DeleteBuilder) OrWhere(condition string, args ...any) *DeleteBuilder {
	b.whereBuilder.OrWhere(condition, args...)
	return b
}

// WhereIn 添加 IN 条件
func (b *DeleteBuilder) WhereIn(field string, values []any) *DeleteBuilder {
	b.whereBuilder.WhereIn(field, values)
	return b
}

// WhereNotIn 添加 NOT IN 条件
func (b *DeleteBuilder) WhereNotIn(field string, values []any) *DeleteBuilder {
	b.whereBuilder.WhereNotIn(field, values)
	return b
}

// WhereBetween 添加 BETWEEN 条件
func (b *DeleteBuilder) WhereBetween(field string, min, max any) *DeleteBuilder {
	b.whereBuilder.WhereBetween(field, min, max)
	return b
}

// WhereNotBetween 添加 NOT BETWEEN 条件
func (b *DeleteBuilder) WhereNotBetween(field string, min, max any) *DeleteBuilder {
	b.whereBuilder.WhereNotBetween(field, min, max)
	return b
}

// WhereNull 添加 IS NULL 条件
func (b *DeleteBuilder) WhereNull(field string) *DeleteBuilder {
	b.whereBuilder.WhereNull(field)
	return b
}

// WhereNotNull 添加 IS NOT NULL 条件
func (b *DeleteBuilder) WhereNotNull(field string) *DeleteBuilder {
	b.whereBuilder.WhereNotNull(field)
	return b
}

// WhereLike 添加 LIKE 条件
func (b *DeleteBuilder) WhereLike(field string, pattern string) *DeleteBuilder {
	b.whereBuilder.WhereLike(field, pattern)
	return b
}

// Exec 执行删除
func (b *DeleteBuilder) Exec() (int64, error) {
	if b.table == "" {
		return 0, uerror.New("table name is required")
	}

	// 构建 SQL
	var sql strings.Builder
	var args []any

	sql.WriteString("DELETE FROM ")
	sql.WriteString(b.table)

	// WHERE 子句
	if b.whereBuilder.HasConditions() {
		whereSQL, whereArgs, _ := buildWhereClause(b.whereBuilder.GetConditions(), 1)
		sql.WriteString(" WHERE ")
		sql.WriteString(whereSQL)
		args = append(args, whereArgs...)
	}

	// 执行
	return b.Connection.Exec(b.ctx, sql.String(), args...)
}
