package rule

import (
	"fmt"

	"github.com/whosafe/uf/uvalidator"
	"github.com/whosafe/uf/uvalidator/i18n"
)

// MinValue 最小值验证规则（专门用于数值）
type MinValue struct {
	MinValue int // 最小值
}

// Validate 执行验证
func (m *MinValue) Validate(value any) bool {
	switch v := value.(type) {
	case int:
		return v >= m.MinValue
	case int64:
		return int(v) >= m.MinValue
	case float64:
		return int(v) >= m.MinValue
	default:
		return false
	}
}

// GetMessage 获取错误消息
func (m *MinValue) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("min", lang...)

	// 替换占位符
	msg := template
	msg = replaceAll(msg, "{field}", field)
	msg = replaceAll(msg, "{param}", fmt.Sprintf("%d", m.MinValue))

	return msg
}

// Name 规则名称
func (m *MinValue) Name() string {
	return "min"
}

// NewMinValue 创建最小值验证规则
func NewMinValue(minValue int) *MinValue {
	return &MinValue{MinValue: minValue}
}
