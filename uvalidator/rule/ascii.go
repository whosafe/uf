package rule

import (
	"iutime.com/utime/uf/uvalidator"
	"iutime.com/utime/uf/uvalidator/i18n"
)

// ASCII ASCII字符验证规则 (0-127)
type ASCII struct{}

// Validate 执行验证
func (a *ASCII) Validate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	if str == "" {
		return true
	}

	for _, r := range str {
		if r > 127 {
			return false
		}
	}
	return true
}

// GetMessage 获取错误消息
func (a *ASCII) GetMessage(field string, params map[string]string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("ascii", lang...)
	return replaceAll(template, "{field}", field)
}

// Name 规则名称
func (a *ASCII) Name() string {
	return "ascii"
}

// NewASCII 创建ASCII验证规则
func NewASCII() *ASCII {
	return &ASCII{}
}
