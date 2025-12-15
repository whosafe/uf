package rule

import (
	"strings"

	"iutime.com/utime/uf/uvalidator"

	"iutime.com/utime/uf/uvalidator/i18n"
)

// Contains 包含子串验证规则
type Contains struct {
	Substring string
}

// Validate 执行验证
func (c *Contains) Validate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	if str == "" {
		return true
	}

	return strings.Contains(str, c.Substring)
}

// GetMessage 获取错误消息
func (c *Contains) GetMessage(field string, params map[string]string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("contains", lang...)
	msg := replaceAll(template, "{field}", field)
	msg = replaceAll(msg, "{param}", c.Substring)
	return msg
}

// Name 规则名称
func (c *Contains) Name() string {
	return "contains"
}

// NewContains 创建包含验证规则
func NewContains(substring string) *Contains {
	return &Contains{Substring: substring}
}

// StartsWith �?..开头验证规�?
type StartsWith struct {
	Prefix string
}

// Validate 执行验证
func (s *StartsWith) Validate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	if str == "" {
		return true
	}

	return strings.HasPrefix(str, s.Prefix)
}

// GetMessage 获取错误消息
func (s *StartsWith) GetMessage(field string, params map[string]string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("starts_with", lang...)
	msg := replaceAll(template, "{field}", field)
	msg = replaceAll(msg, "{param}", s.Prefix)
	return msg
}

// Name 规则名称
func (s *StartsWith) Name() string {
	return "starts_with"
}

// NewStartsWith 创建�?..开头验证规�?
func NewStartsWith(prefix string) *StartsWith {
	return &StartsWith{Prefix: prefix}
}

// EndsWith �?..结尾验证规则
type EndsWith struct {
	Suffix string
}

// Validate 执行验证
func (e *EndsWith) Validate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	if str == "" {
		return true
	}

	return strings.HasSuffix(str, e.Suffix)
}

// GetMessage 获取错误消息
func (e *EndsWith) GetMessage(field string, params map[string]string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("ends_with", lang...)
	msg := replaceAll(template, "{field}", field)
	msg = replaceAll(msg, "{param}", e.Suffix)
	return msg
}

// Name 规则名称
func (e *EndsWith) Name() string {
	return "ends_with"
}

// NewEndsWith 创建�?..结尾验证规则
func NewEndsWith(suffix string) *EndsWith {
	return &EndsWith{Suffix: suffix}
}
