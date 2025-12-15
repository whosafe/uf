package rule

import (
	"fmt"

	"github.com/whosafe/uf/uvalidator"
	"github.com/whosafe/uf/uvalidator/i18n"
)

// MaxLength 最大长度验证规则（专门用于字符串）
type MaxLength struct {
	MaxValue int // 最大长度
}

// Validate 执行验证
func (m *MaxLength) Validate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}
	return len(str) <= m.MaxValue
}

// GetMessage 获取错误消息
func (m *MaxLength) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("max_length", lang...)

	// 替换占位符
	msg := template
	msg = replaceAll(msg, "{field}", field)
	msg = replaceAll(msg, "{param}", fmt.Sprintf("%d", m.MaxValue))

	return msg
}

// Name 规则名称
func (m *MaxLength) Name() string {
	return "max_length"
}

// NewMaxLength 创建最大长度验证规则
func NewMaxLength(maxValue int) *MaxLength {
	return &MaxLength{MaxValue: maxValue}
}
