package sqlite

import (
	"context"
	"database/sql"
	"strings"

	"github.com/whosafe/uf/uerror"
)

// TxDeleteBuilder 事务删除构建器
type TxDeleteBuilder struct {
	ctx          context.Context
	tx           *sql.Tx
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
	var sqlStr strings.Builder
	var args []any

	sqlStr.WriteString("DELETE FROM ")
	sqlStr.WriteString(b.table)

	// WHERE 子句
	if b.whereBuilder.HasConditions() {
		whereSQL, whereArgs := buildWhereClause(b.whereBuilder.GetConditions())
		sqlStr.WriteString(" WHERE ")
		sqlStr.WriteString(whereSQL)
		args = append(args, whereArgs...)
	}

	// 执行
	result, err := b.tx.ExecContext(b.ctx, sqlStr.String(), args...)
	if err != nil {
		return 0, uerror.Wrap(err, "事务删除失败")
	}
	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, nil
}
