package postgresql

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/whosafe/uf/uerror"
)

// TxDeleteBuilder 事务删除构建器
type TxDeleteBuilder struct {
	ctx          context.Context
	tx           pgx.Tx
	table        string
	whereBuilder WhereBuilder
}

// Table 设置表名
func (b *TxDeleteBuilder) Table(table string) *TxDeleteBuilder {
	b.table = table
	return b
}

// Where 添加 WHERE 条件
func (b *TxDeleteBuilder) Where(condition string, args ...any) *TxDeleteBuilder {
	b.whereBuilder.Where(condition, args...)
	return b
}

// Exec 执行删除
func (b *TxDeleteBuilder) Exec() (int64, error) {
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
	result, err := b.tx.Exec(b.ctx, sql.String(), args...)
	if err != nil {
		return 0, uerror.Wrap(err, "事务删除失败")
	}
	return result.RowsAffected(), nil
}
