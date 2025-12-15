package postgresql

// joinClause JOIN 子句
type joinClause struct {
	joinType string // INNER, LEFT, RIGHT, FULL
	table    string
	on       string
}

// whereClause WHERE/HAVING 子句
type whereClause struct {
	condition string
	args      []any
	isOr      bool // 是否为 OR 条件
}

// orderClause ORDER BY 子句
type orderClause struct {
	field string
	desc  bool
}
