package postgresql

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/whosafe/uf/uerror"
)

// TxInsertBuilder 事务插入构建器
type TxInsertBuilder struct {
	ctx     context.Context
	tx      pgx.Tx
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

	// 执行
	result, err := b.tx.Exec(b.ctx, sql.String(), b.values...)
	if err != nil {
		return 0, uerror.Wrap(err, "事务插入失败")
	}
	return result.RowsAffected(), nil
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
	rows, err := b.tx.Query(b.ctx, sql.String(), b.values...)
	if err != nil {
		return uerror.Wrap(err, "事务插入失败")
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
