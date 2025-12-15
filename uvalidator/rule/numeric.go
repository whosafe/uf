package rule

import (
	"unicode"

	"github.com/whosafe/uf/uvalidator"

	"github.com/whosafe/uf/uvalidator/i18n"
)

// Numeric 只包含数字验证规则
type Numeric struct{}

// Validate 执行验证
func (n *Numeric) Validate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	if str == "" {
		return true
	}

	for _, r := range str {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// GetMessage 获取错误消息
func (n *Numeric) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("numeric", lang...)
	return replaceAll(template, "{field}", field)
}

// Name 规则名称
func (n *Numeric) Name() string {
	return "numeric"
}

// NewNumeric 创建数字验证规则
func NewNumeric() *Numeric {
	return &Numeric{}
}
