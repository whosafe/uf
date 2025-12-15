package rule

import (
	"fmt"

	"github.com/whosafe/uf/uvalidator"
	"github.com/whosafe/uf/uvalidator/i18n"
)

// MinLength 最小长度验证规则（专门用于字符串）
type MinLength struct {
	MinValue int // 最小长度
}

// Validate 执行验证
func (m *MinLength) Validate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}
	return len(str) >= m.MinValue
}

// GetMessage 获取错误消息
func (m *MinLength) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("min_length", lang...)

	// 替换占位符
	msg := template
	msg = replaceAll(msg, "{field}", field)
	msg = replaceAll(msg, "{param}", fmt.Sprintf("%d", m.MinValue))

	return msg
}

// Name 规则名称
func (m *MinLength) Name() string {
	return "min_length"
}

// NewMinLength 创建最小长度验证规则
func NewMinLength(minValue int) *MinLength {
	return &MinLength{MinValue: minValue}
}
