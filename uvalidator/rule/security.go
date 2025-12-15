package rule

import (
	"regexp"
	"strings"

	"github.com/whosafe/uf/uvalidator"

	"github.com/whosafe/uf/uvalidator/i18n"
)

// HTML标签正则
var htmlRegex = regexp.MustCompile(`<[^>]+>`)

// SQL注入关键
var sqlKeywords = []string{
	"select", "insert", "update", "delete", "drop", "create",
	"alter", "exec", "execute", "script", "union", "where",
	"or", "and", "--", "/*", "*/", "xp_", "sp_",
}

// XSS攻击字符
var xssPatterns = []*regexp.Regexp{
	regexp.MustCompile(`<script[^>]*>.*?</script>`),
	regexp.MustCompile(`javascript:`),
	regexp.MustCompile(`on\w+\s*=`),
	regexp.MustCompile(`<iframe[^>]*>`),
}

// StrongPassword 强密码验证规则
type StrongPassword struct {
	MinLength int // 最小长默认8
}

// Validate 执行验证
func (s *StrongPassword) Validate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	if str == "" {
		return true
	}

	minLen := s.MinLength
	if minLen == 0 {
		minLen = 8
	}

	// 检查长
	if len(str) < minLen {
		return false
	}

	// 检查是否包含大写字
	hasUpper := false
	// 检查是否包含小写字
	hasLower := false
	// 检查是否包含数
	hasDigit := false
	// 检查是否包含特殊字
	hasSpecial := false

	for _, r := range str {
		switch {
		case r >= 'A' && r <= 'Z':
			hasUpper = true
		case r >= 'a' && r <= 'z':
			hasLower = true
		case r >= '0' && r <= '9':
			hasDigit = true
		default:
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecial
}

// GetMessage 获取错误消息
func (s *StrongPassword) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("strong_password", lang...)
	return replaceAll(template, "{field}", field)
}

// Name 规则名称
func (s *StrongPassword) Name() string {
	return "strong_password"
}

// NewStrongPassword 创建强密码验证规则
func NewStrongPassword(minLength ...int) *StrongPassword {
	sp := &StrongPassword{MinLength: 8}
	if len(minLength) > 0 {
		sp.MinLength = minLength[0]
	}
	return sp
}

// NoHTML 不包含HTML标签验证规则
type NoHTML struct{}

// Validate 执行验证
func (n *NoHTML) Validate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	if str == "" {
		return true
	}

	return !htmlRegex.MatchString(str)
}

// GetMessage 获取错误消息
func (n *NoHTML) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("no_html", lang...)
	return replaceAll(template, "{field}", field)
}

// Name 规则名称
func (n *NoHTML) Name() string {
	return "no_html"
}

// NewNoHTML 创建不包含HTML验证规则
func NewNoHTML() *NoHTML {
	return &NoHTML{}
}

// NoSQL 不包含SQL注入字符验证规则
type NoSQL struct{}

// Validate 执行验证
func (n *NoSQL) Validate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	if str == "" {
		return true
	}

	lowerStr := strings.ToLower(str)
	for _, keyword := range sqlKeywords {
		if strings.Contains(lowerStr, keyword) {
			return false
		}
	}

	return true
}

// GetMessage 获取错误消息
func (n *NoSQL) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("no_sql", lang...)
	return replaceAll(template, "{field}", field)
}

// Name 规则名称
func (n *NoSQL) Name() string {
	return "no_sql"
}

// NewNoSQL 创建不包含SQL注入验证规则
func NewNoSQL() *NoSQL {
	return &NoSQL{}
}

// NoXSS 不包含XSS攻击字符验证规则
type NoXSS struct{}

// Validate 执行验证
func (n *NoXSS) Validate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	if str == "" {
		return true
	}

	lowerStr := strings.ToLower(str)
	for _, pattern := range xssPatterns {
		if pattern.MatchString(lowerStr) {
			return false
		}
	}

	return true
}

// GetMessage 获取错误消息
func (n *NoXSS) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("no_xss", lang...)
	return replaceAll(template, "{field}", field)
}

// Name 规则名称
func (n *NoXSS) Name() string {
	return "no_xss"
}

// NewNoXSS 创建不包含XSS验证规则
func NewNoXSS() *NoXSS {
	return &NoXSS{}
}
