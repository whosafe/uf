package postgresql

import (
	"fmt"
	"strings"
)

// buildWhereClause 构建 WHERE 子句
// 返回 WHERE SQL 字符串和参数列表
func buildWhereClause(whereClauses []whereClause, startParamIndex int) (string, []any, int) {
	if len(whereClauses) == 0 {
		return "", nil, startParamIndex
	}

	var sql strings.Builder
	var args []any
	paramIndex := startParamIndex

	for i, w := range whereClauses {
		// 替换 ? 为 $1, $2, ...
		condition := w.condition
		for range w.args {
			condition = strings.Replace(condition, "?",
				fmt.Sprintf("$%d", paramIndex), 1)
			paramIndex++
		}

		// 添加连接符
		if i > 0 {
			if w.isOr {
				sql.WriteString(" OR ")
			} else {
				sql.WriteString(" AND ")
			}
		}
		sql.WriteString(condition)
		args = append(args, w.args...)
	}

	return sql.String(), args, paramIndex
}

// buildInCondition 构建 IN 或 NOT IN 条件
func buildInCondition(field string, values []any, notIn bool) whereClause {
	if len(values) == 0 {
		return whereClause{}
	}

	placeholders := make([]string, len(values))
	for i := range values {
		placeholders[i] = "?"
	}

	operator := "IN"
	if notIn {
		operator = "NOT IN"
	}

	condition := field + " " + operator + " (" + strings.Join(placeholders, ", ") + ")"
	return whereClause{
		condition: condition,
		args:      values,
		isOr:      false,
	}
}

// buildBetweenCondition 构建 BETWEEN 或 NOT BETWEEN 条件
func buildBetweenCondition(field string, min, max any, notBetween bool) whereClause {
	operator := "BETWEEN"
	if notBetween {
		operator = "NOT BETWEEN"
	}

	condition := field + " " + operator + " ? AND ?"
	return whereClause{
		condition: condition,
		args:      []any{min, max},
		isOr:      false,
	}
}

// buildNullCondition 构建 IS NULL 或 IS NOT NULL 条件
func buildNullCondition(field string, notNull bool) whereClause {
	operator := "IS NULL"
	if notNull {
		operator = "IS NOT NULL"
	}

	condition := field + " " + operator
	return whereClause{
		condition: condition,
		args:      nil,
		isOr:      false,
	}
}

// buildLikeCondition 构建 LIKE 条件
func buildLikeCondition(field string, pattern string) whereClause {
	condition := field + " LIKE ?"
	return whereClause{
		condition: condition,
		args:      []any{pattern},
		isOr:      false,
	}
}
