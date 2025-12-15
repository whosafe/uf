package postgresql

import (
	"regexp"
	"strings"

	"github.com/whosafe/uf/uerror"
)

// 标识符验证正则表达式
// 允许字母、数字、下划线,必须以字母或下划线开头
var identifierRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

// 限定标识符正则表达式 (支持 table.column 格式)
var qualifiedIdentifierRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*(\.[a-zA-Z_][a-zA-Z0-9_]*)?$`)

// ValidateIdentifier 验证 SQL 标识符(表名、字段名等)
// 【安全修复】防止通过标识符注入恶意 SQL
func ValidateIdentifier(name string) error {
	if name == "" {
		return uerror.New("标识符不能为空")
	}

	// 长度限制(PostgreSQL 标识符最大 63 字节)
	if len(name) > 63 {
		return uerror.New("标识符长度不能超过 63 个字符")
	}

	// 检查是否符合标识符规范
	if !identifierRegex.MatchString(name) {
		return uerror.New("标识符只能包含字母、数字和下划线,且必须以字母或下划线开头")
	}

	// 检查是否是 PostgreSQL 保留字
	if isReservedKeyword(name) {
		return uerror.New("不能使用 PostgreSQL 保留字作为标识符")
	}

	return nil
}

// ValidateIdentifiers 批量验证标识符
func ValidateIdentifiers(names ...string) error {
	for _, name := range names {
		if err := ValidateIdentifier(name); err != nil {
			return err
		}
	}
	return nil
}

// ValidateFieldExpression 验证字段表达式 (用于 SELECT, ORDER BY, GROUP BY)
// 支持:
// - 简单标识符: id, name
// - 限定标识符: users.id, u.name
// - SQL 表达式: COUNT(*), SUM(amount), COUNT(o.id) as order_count
// 【安全修复】防止 SQL 注入,同时支持常见的 SQL 用法
func ValidateFieldExpression(expr string) error {
	if expr == "" {
		return uerror.New("字段表达式不能为空")
	}

	// 移除首尾空格
	expr = strings.TrimSpace(expr)

	// 检查是否包含危险的 SQL 关键字或字符
	lowerExpr := strings.ToLower(expr)
	dangerousPatterns := []string{
		";", "--", "/*", "*/", "xp_", "sp_",
		"drop ", "delete ", "insert ", "update ",
		"exec ", "execute ", "union ", "script",
	}
	for _, pattern := range dangerousPatterns {
		if strings.Contains(lowerExpr, pattern) {
			return uerror.New("字段表达式包含不安全的内容")
		}
	}

	// 如果是简单标识符或限定标识符,直接验证
	if qualifiedIdentifierRegex.MatchString(expr) {
		return nil
	}

	// 如果包含函数调用或 AS 别名,进行更宽松的验证
	// 允许的字符: 字母、数字、下划线、点、括号、星号、逗号、空格、AS
	allowedCharsRegex := regexp.MustCompile(`^[a-zA-Z0-9_\.\(\)\*\,\s]+$`)
	if !allowedCharsRegex.MatchString(expr) {
		return uerror.New("字段表达式包含不允许的字符")
	}

	return nil
}

// ValidateFieldExpressions 批量验证字段表达式
func ValidateFieldExpressions(exprs ...string) error {
	for _, expr := range exprs {
		if err := ValidateFieldExpression(expr); err != nil {
			return err
		}
	}
	return nil
}

// ValidateTableExpression 验证表名表达式 (支持表别名)
// 支持:
// - 简单表名: users, orders
// - 带别名的表名: users u, orders o
// 【安全修复】防止 SQL 注入,同时支持表别名
func ValidateTableExpression(expr string) error {
	if expr == "" {
		return uerror.New("表名表达式不能为空")
	}

	// 移除首尾空格
	expr = strings.TrimSpace(expr)

	// 检查是否包含危险的 SQL 关键字或字符
	lowerExpr := strings.ToLower(expr)
	dangerousPatterns := []string{
		";", "--", "/*", "*/", "xp_", "sp_",
		"drop ", "delete ", "insert ", "update ",
		"exec ", "execute ", "union ", "script",
	}
	for _, pattern := range dangerousPatterns {
		if strings.Contains(lowerExpr, pattern) {
			return uerror.New("表名表达式包含不安全的内容")
		}
	}

	// 分割表名和别名
	parts := strings.Fields(expr)
	if len(parts) > 2 {
		return uerror.New("表名表达式格式错误")
	}

	// 验证表名
	if !identifierRegex.MatchString(parts[0]) {
		return uerror.New("表名只能包含字母、数字和下划线,且必须以字母或下划线开头")
	}

	// 如果有别名,验证别名
	if len(parts) == 2 {
		if !identifierRegex.MatchString(parts[1]) {
			return uerror.New("表别名只能包含字母、数字和下划线,且必须以字母或下划线开头")
		}
	}

	return nil
}

// isReservedKeyword 检查是否是 PostgreSQL 保留字
func isReservedKeyword(name string) bool {
	// 常见的 PostgreSQL 保留字(小写)
	reserved := map[string]bool{
		"select": true, "insert": true, "update": true, "delete": true,
		"from": true, "where": true, "join": true, "on": true,
		"and": true, "or": true, "not": true, "null": true,
		"true": true, "false": true, "table": true, "column": true,
		"index": true, "view": true, "database": true, "schema": true,
		"user": true, "group": true, "order": true, "by": true,
		"limit": true, "offset": true, "union": true, "all": true,
		"distinct": true, "as": true, "case": true, "when": true,
		"then": true, "else": true, "end": true, "exists": true,
		"in": true, "between": true, "like": true, "is": true,
	}

	return reserved[strings.ToLower(name)]
}
