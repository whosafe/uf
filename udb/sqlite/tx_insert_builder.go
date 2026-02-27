package sqlite

import (
	"context"
	"database/sql"
	"strings"

	"github.com/whosafe/uf/uerror"
)

// TxInsertBuilder 事务插入构建器
type TxInsertBuilder struct {
	ctx     context.Context
	tx      *sql.Tx
	table   string
	columns []string
	values  []any
}

// Table 设置表名
func (b *TxInsertBuilder) Table(table string) *TxInsertBuilder {
	b.table = table
	return b
}

// Columns 设置列名
func (b *TxInsertBuilder) Columns(cols ...string) *TxInsertBuilder {
	b.columns = cols
	return b
}

// Values 设置值
func (b *TxInsertBuilder) Values(vals ...any) *TxInsertBuilder {
	b.values = vals
	return b
}

// Exec 执行插入
func (b *TxInsertBuilder) Exec() (int64, error) {
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
	var sqlStr strings.Builder
	sqlStr.WriteString("INSERT INTO ")
	sqlStr.WriteString(b.table)
	sqlStr.WriteString(" (")
	sqlStr.WriteString(strings.Join(b.columns, ", "))
	sqlStr.WriteString(") VALUES (")

	// 添加占位符
	placeholders := make([]string, len(b.values))
	for i := range b.values {
		placeholders[i] = "?"
	}
	sqlStr.WriteString(strings.Join(placeholders, ", "))
	sqlStr.WriteString(")")

	// 执行
	result, err := b.tx.ExecContext(b.ctx, sqlStr.String(), b.values...)
	if err != nil {
		return 0, uerror.Wrap(err, "事务插入失败")
	}
	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, nil
}

// ExecReturning 执行插入并返回数据
func (b *TxInsertBuilder) ExecReturning(dest Scanner) error {
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
	var sqlStr strings.Builder
	sqlStr.WriteString("INSERT INTO ")
	sqlStr.WriteString(b.table)
	sqlStr.WriteString(" (")
	sqlStr.WriteString(strings.Join(b.columns, ", "))
	sqlStr.WriteString(") VALUES (")

	// 添加占位符
	placeholders := make([]string, len(b.values))
	for i := range b.values {
		placeholders[i] = "?"
	}
	sqlStr.WriteString(strings.Join(placeholders, ", "))
	sqlStr.WriteString(") RETURNING *")

	// 执行查询
	rows, err := b.tx.QueryContext(b.ctx, sqlStr.String(), b.values...)
	if err != nil {
		return uerror.Wrap(err, "事务插入失败")
	}
	defer rows.Close()

	if !rows.Next() {
		return ErrNoRows
	}

	columns, err := rows.Columns()
	if err != nil {
		return uerror.Wrap(err, "获取列名失败")
	}

	values := make([]any, len(columns))
	valuePtrs := make([]any, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	if err := rows.Scan(valuePtrs...); err != nil {
		return uerror.Wrap(err, "获取插入结果失败")
	}

	// 逐个字段调用 Scan
	for i, col := range columns {
		if err := dest.Scan(col, values[i]); err != nil {
			return uerror.Wrap(err, "扫描插入结果失败")
		}
	}

	return nil
}
