package postgresql

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/whosafe/uf/uerror"
)

// TxUpdateBuilder 事务更新构建器
type TxUpdateBuilder struct {
	ctx          context.Context
	tx           pgx.Tx
	table        string
	sets         map[string]any
	whereBuilder WhereBuilder
}

// Table 设置表名
func (b *TxUpdateBuilder) Table(table string) *TxUpdateBuilder {
	b.table = table
	return b
}

// Set 设置字段值
func (b *TxUpdateBuilder) Set(column string, value any) *TxUpdateBuilder {
	if b.sets == nil {
		b.sets = make(map[string]any)
	}
	b.sets[column] = value
	return b
}

// SetMap 批量设置字段
func (b *TxUpdateBuilder) SetMap(data map[string]any) *TxUpdateBuilder {
	if b.sets == nil {
		b.sets = make(map[string]any)
	}
	for k, v := range data {
		b.sets[k] = v
	}
	return b
}

// Where 添加 WHERE 条件
func (b *TxUpdateBuilder) Where(condition string, args ...any) *TxUpdateBuilder {
	b.whereBuilder.Where(condition, args...)
	return b
}

// Exec 执行更新
func (b *TxUpdateBuilder) Exec() (int64, error) {
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
	result, err := b.tx.Exec(b.ctx, sql.String(), args...)
	if err != nil {
		return 0, uerror.Wrap(err, "事务更新失败")
	}
	return result.RowsAffected(), nil
}
