package rule

import (
	"regexp"

	"iutime.com/utime/uf/uvalidator"

	"iutime.com/utime/uf/uvalidator/i18n"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// Email 邮箱验证规则
type Email struct{}

// Validate 执行验证
func (e *Email) Validate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	if str == "" {
		return true // 空值由 required 验证
	}

	return emailRegex.MatchString(str)
}

// GetMessage 获取错误消息
func (e *Email) GetMessage(field string, params map[string]string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("email", lang...)
	return replaceAll(template, "{field}", field)
}

// Name 规则名称
func (e *Email) Name() string {
	return "email"
}

// NewEmail 创建邮箱验证规则
func NewEmail() *Email {
	return &Email{}
}
