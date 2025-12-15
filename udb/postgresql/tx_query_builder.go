package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/whosafe/uf/uerror"
)

// TxQueryBuilder 事务查询构建器
type TxQueryBuilder struct {
	ctx    context.Context
	tx     pgx.Tx
	config *Config

	// 查询构建器字段(与 QueryBuilder 相同)
	selectFields []string
	distinct     bool
	table        string
	joins        []joinClause
	whereBuilder WhereBuilder
	groupBy      []string
	having       []whereClause
	orderBy      []orderClause
	limit        int
	offset       int
}

// Table 设置表名
func (q *TxQueryBuilder) Table(name string) *TxQueryBuilder {
	q.table = name
	return q
}

// Select 设置查询字段
func (q *TxQueryBuilder) Select(fields ...string) *TxQueryBuilder {
	q.selectFields = fields
	return q
}

// Where 添加 WHERE 条件
func (q *TxQueryBuilder) Where(condition string, args ...any) *TxQueryBuilder {
	q.whereBuilder.Where(condition, args...)
	return q
}

// OrWhere 添加 OR WHERE 条件
func (q *TxQueryBuilder) OrWhere(condition string, args ...any) *TxQueryBuilder {
	q.whereBuilder.OrWhere(condition, args...)
	return q
}

// WhereIn 添加 IN 条件
func (q *TxQueryBuilder) WhereIn(field string, values []any) *TxQueryBuilder {
	q.whereBuilder.WhereIn(field, values)
	return q
}

// WhereNotIn 添加 NOT IN 条件
func (q *TxQueryBuilder) WhereNotIn(field string, values []any) *TxQueryBuilder {
	q.whereBuilder.WhereNotIn(field, values)
	return q
}

// WhereBetween 添加 BETWEEN 条件
func (q *TxQueryBuilder) WhereBetween(field string, min, max any) *TxQueryBuilder {
	q.whereBuilder.WhereBetween(field, min, max)
	return q
}

// WhereNotBetween 添加 NOT BETWEEN 条件
func (q *TxQueryBuilder) WhereNotBetween(field string, min, max any) *TxQueryBuilder {
	q.whereBuilder.WhereNotBetween(field, min, max)
	return q
}

// WhereNull 添加 IS NULL 条件
func (q *TxQueryBuilder) WhereNull(field string) *TxQueryBuilder {
	q.whereBuilder.WhereNull(field)
	return q
}

// WhereNotNull 添加 IS NOT NULL 条件
func (q *TxQueryBuilder) WhereNotNull(field string) *TxQueryBuilder {
	q.whereBuilder.WhereNotNull(field)
	return q
}

// WhereLike 添加 LIKE 条件
func (q *TxQueryBuilder) WhereLike(field string, pattern string) *TxQueryBuilder {
	q.whereBuilder.WhereLike(field, pattern)
	return q
}

// Limit 限制数量
func (q *TxQueryBuilder) Limit(n int) *TxQueryBuilder {
	q.limit = n
	return q
}

// Offset 偏移量
func (q *TxQueryBuilder) Offset(n int) *TxQueryBuilder {
	q.offset = n
	return q
}

// OrderBy 升序排序
func (q *TxQueryBuilder) OrderBy(field string) *TxQueryBuilder {
	q.orderBy = append(q.orderBy, orderClause{
		field: field,
		desc:  false,
	})
	return q
}

// BuildSQL 构建 SQL(复用 QueryBuilder 的逻辑)
func (q *TxQueryBuilder) BuildSQL() (string, []any) {
	// 创建临时 QueryBuilder 来复用 BuildSQL 逻辑
	qb := &QueryBuilder{
		selectFields: q.selectFields,
		distinct:     q.distinct,
		table:        q.table,
		joins:        q.joins,
		whereBuilder: q.whereBuilder,
		groupBy:      q.groupBy,
		having:       q.having,
		orderBy:      q.orderBy,
		limit:        q.limit,
		offset:       q.offset,
	}
	return qb.BuildSQL()
}

// Scan 扫描单行
func (q *TxQueryBuilder) Scan(dest Scanner) error {
	sql, args := q.BuildSQL()

	rows, err := q.tx.Query(q.ctx, sql, args...)
	if err != nil {
		return uerror.Wrap(err, "事务查询失败")
	}
	defer rows.Close()

	if !rows.Next() {
		return ErrNoRows
	}

	fieldDescs := rows.FieldDescriptions()
	values, err := rows.Values()
	if err != nil {
		return uerror.Wrap(err, "获取行数据失败")
	}

	// 逐个字段调用 Scan
	for i, fd := range fieldDescs {
		key := string(fd.Name)
		if err := dest.Scan(key, values[i]); err != nil {
			return uerror.Wrap(err, "扫描字段失败")
		}
	}

	return nil
}
