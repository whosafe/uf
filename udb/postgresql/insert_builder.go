package postgresql

import (
	"context"
	"fmt"
	"strings"

	"github.com/whosafe/uf/uerror"
)

// InsertBuilder 插入构建器
type InsertBuilder struct {
	ctx        context.Context
	Connection *Connection
	table      string
	columns    []string
	values     []any
}

// Table 设置表名
func (b *InsertBuilder) Table(table string) *InsertBuilder {
	// 【安全修复】验证表名,防止 SQL 注入
	if err := ValidateIdentifier(table); err != nil {
		panic(err)
	}
	b.table = table
	return b
}

// Columns 设置列名
func (b *InsertBuilder) Columns(cols ...string) *InsertBuilder {
	// 【安全修复】验证列名,防止 SQL 注入
	if err := ValidateIdentifiers(cols...); err != nil {
		panic(err)
	}
	b.columns = cols
	return b
}

// Values 设置值
func (b *InsertBuilder) Values(vals ...any) *InsertBuilder {
	b.values = vals
	return b
}

// Exec 执行插入
func (b *InsertBuilder) Exec() (int64, error) {
	if b.table == "" {
		return 0, uerror.New("table name is required")
	}
	if len(b.columns) == 0 {
		return 0, uerror.New("columns are required")
	}
	if len(b.values) == 0 {
		return 0, uerror.New("values are required")
	}
	if len(b.columns) != len(b.values) {
		return 0, uerror.New("columns and values count mismatch")
	}

	// 构建 SQL
	var sql strings.Builder
	sql.WriteString("INSERT INTO ")
	sql.WriteString(b.table)
	sql.WriteString(" (")
	sql.WriteString(strings.Join(b.columns, ", "))
	sql.WriteString(") VALUES (")

	// 添加占位符
	placeholders := make([]string, len(b.values))
	for i := range b.values {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}
	sql.WriteString(strings.Join(placeholders, ", "))
	sql.WriteString(")")

	// 执行(Exec 方法内部已包含日志)
	return b.Connection.Exec(b.ctx, sql.String(), b.values...)
}

// ExecReturning 执行插入并返回数据
func (b *InsertBuilder) ExecReturning(dest Scanner) error {
	if b.table == "" {
		return uerror.New("table name is required")
	}
	if len(b.columns) == 0 {
		return uerror.New("columns are required")
	}
	if len(b.values) == 0 {
		return uerror.New("values are required")
	}
	if len(b.columns) != len(b.values) {
		return uerror.New("columns and values count mismatch")
	}

	// 构建 SQL
	var sql strings.Builder
	sql.WriteString("INSERT INTO ")
	sql.WriteString(b.table)
	sql.WriteString(" (")
	sql.WriteString(strings.Join(b.columns, ", "))
	sql.WriteString(") VALUES (")

	// 添加占位符
	placeholders := make([]string, len(b.values))
	for i := range b.values {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}
	sql.WriteString(strings.Join(placeholders, ", "))
	sql.WriteString(") RETURNING *")

	// 执行查询
	rows, err := b.Connection.pool.Query(b.ctx, sql.String(), b.values...)
	if err != nil {
		return uerror.Wrap(err, "插入数据失败")
	}
	defer rows.Close()

	if !rows.Next() {
		return ErrNoRows
	}

	// 获取列名和值
	fieldDescs := rows.FieldDescriptions()
	values, err := rows.Values()
	if err != nil {
		return uerror.Wrap(err, "获取插入结果失败")
	}

	// 逐个字段调用 Scan
	for i, fd := range fieldDescs {
		key := string(fd.Name)
		if err := dest.Scan(key, values[i]); err != nil {
			return uerror.Wrap(err, "扫描插入结果失败")
		}
	}

	return nil
}
