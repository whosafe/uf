package rule

import (
	"fmt"

	"github.com/whosafe/uf/uvalidator"
	"github.com/whosafe/uf/uvalidator/i18n"
)

// MaxValue 最大值验证规则（专门用于数值）
type MaxValue struct {
	MaxValue int // 最大值
}

// Validate 执行验证
func (m *MaxValue) Validate(value any) bool {
	switch v := value.(type) {
	case int:
		return v <= m.MaxValue
	case int64:
		return int(v) <= m.MaxValue
	case float64:
		return int(v) <= m.MaxValue
	default:
		return false
	}
}

// GetMessage 获取错误消息
func (m *MaxValue) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("max", lang...)

	// 替换占位符
	msg := template
	msg = replaceAll(msg, "{field}", field)
	msg = replaceAll(msg, "{param}", fmt.Sprintf("%d", m.MaxValue))

	return msg
}

// Name 规则名称
func (m *MaxValue) Name() string {
	return "max"
}

// NewMaxValue 创建最大值验证规则
func NewMaxValue(maxValue int) *MaxValue {
	return &MaxValue{MaxValue: maxValue}
}
