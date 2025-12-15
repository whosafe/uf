package postgresql

import (
	"fmt"
	"strings"
)

// BuildSQL 构建 SQL 语句
func (q *QueryBuilder) BuildSQL() (string, []any) {
	var sql strings.Builder
	var args []any
	paramIndex := 1

	// SELECT
	sql.WriteString("SELECT ")
	if q.distinct {
		sql.WriteString("DISTINCT ")
	}
	if len(q.selectFields) > 0 {
		sql.WriteString(strings.Join(q.selectFields, ", "))
	} else {
		sql.WriteString("*")
	}

	// FROM
	sql.WriteString(" FROM ")
	sql.WriteString(q.table)

	// JOIN
	for _, join := range q.joins {
		sql.WriteString(fmt.Sprintf(" %s JOIN %s ON %s",
			join.joinType, join.table, join.on))
	}

	// WHERE
	if q.whereBuilder.HasConditions() {
		whereSQL, whereArgs, newParamIndex := buildWhereClause(q.whereBuilder.GetConditions(), paramIndex)
		sql.WriteString(" WHERE ")
		sql.WriteString(whereSQL)
		args = append(args, whereArgs...)
		paramIndex = newParamIndex
	}

	// GROUP BY
	if len(q.groupBy) > 0 {
		sql.WriteString(" GROUP BY ")
		sql.WriteString(strings.Join(q.groupBy, ", "))
	}

	// HAVING
	if len(q.having) > 0 {
		sql.WriteString(" HAVING ")
		conditions := make([]string, 0, len(q.having))
		for _, h := range q.having {
			condition := h.condition
			for range h.args {
				condition = strings.Replace(condition, "?",
					fmt.Sprintf("$%d", paramIndex), 1)
				paramIndex++
			}
			conditions = append(conditions, condition)
			args = append(args, h.args...)
		}
		sql.WriteString(strings.Join(conditions, " AND "))
	}

	// ORDER BY
	if len(q.orderBy) > 0 {
		sql.WriteString(" ORDER BY ")
		orders := make([]string, 0, len(q.orderBy))
		for _, o := range q.orderBy {
			order := o.field
			if o.desc {
				order += " DESC"
			}
			orders = append(orders, order)
		}
		sql.WriteString(strings.Join(orders, ", "))
	}

	// LIMIT
	if q.limit > 0 {
		sql.WriteString(fmt.Sprintf(" LIMIT %d", q.limit))
	}

	// OFFSET
	if q.offset > 0 {
		sql.WriteString(fmt.Sprintf(" OFFSET %d", q.offset))
	}

	return sql.String(), args
}
