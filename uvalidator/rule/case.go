package rule

import (
	"strings"

	"github.com/whosafe/uf/uvalidator"

	"github.com/whosafe/uf/uvalidator/i18n"
)

// Lowercase 小写验证规则
type Lowercase struct{}

// Validate 执行验证
func (l *Lowercase) Validate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	if str == "" {
		return true
	}

	return str == strings.ToLower(str)
}

// GetMessage 获取错误消息
func (l *Lowercase) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("lowercase", lang...)
	return replaceAll(template, "{field}", field)
}

// Name 规则名称
func (l *Lowercase) Name() string {
	return "lowercase"
}

// NewLowercase 创建小写验证规则
func NewLowercase() *Lowercase {
	return &Lowercase{}
}

// Uppercase 大写验证规则
type Uppercase struct{}

// Validate 执行验证
func (u *Uppercase) Validate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	if str == "" {
		return true
	}

	return str == strings.ToUpper(str)
}

// GetMessage 获取错误消息
func (u *Uppercase) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("uppercase", lang...)
	return replaceAll(template, "{field}", field)
}

// Name 规则名称
func (u *Uppercase) Name() string {
	return "uppercase"
}

// NewUppercase 创建大写验证规则
func NewUppercase() *Uppercase {
	return &Uppercase{}
}
