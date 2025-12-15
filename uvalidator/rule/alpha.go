package rule

import (
	"unicode"

	"iutime.com/utime/uf/uvalidator"
	"iutime.com/utime/uf/uvalidator/i18n"
)

// Alpha 只包含字母验证规�?
type Alpha struct{}

// Validate 执行验证
func (a *Alpha) Validate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	if str == "" {
		return true
	}

	for _, r := range str {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

// GetMessage 获取错误消息
func (a *Alpha) GetMessage(field string, params map[string]string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("alpha", lang...)
	return replaceAll(template, "{field}", field)
}

// Name 规则名称
func (a *Alpha) Name() string {
	return "alpha"
}

// NewAlpha 创建字母验证规则
func NewAlpha() *Alpha {
	return &Alpha{}
}
